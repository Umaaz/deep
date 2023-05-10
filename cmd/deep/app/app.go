/*
 *     Copyright (C) 2023  Intergral GmbH
 *
 *     This program is free software: you can redistribute it and/or modify
 *     it under the terms of the GNU Affero General Public License as published by
 *     the Free Software Foundation, either version 3 of the License, or
 *     (at your option) any later version.
 *
 *     This program is distributed in the hope that it will be useful,
 *     but WITHOUT ANY WARRANTY; without even the implied warranty of
 *     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *     GNU Affero General Public License for more details.
 *
 *     You should have received a copy of the GNU Affero General Public License
 *     along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package app

import (
	"bytes"
	"context"
	"fmt"
	"github.com/intergral/deep/modules/compactor"
	"github.com/intergral/deep/modules/distributor"
	"github.com/intergral/deep/modules/ingester"
	"github.com/intergral/deep/modules/overrides"
	"github.com/intergral/deep/modules/querier"
	"github.com/intergral/deep/modules/storage"
	"github.com/intergral/deep/modules/tracepoint"
	tpapi "github.com/intergral/deep/modules/tracepoint/api"
	"github.com/intergral/deep/modules/tracepoint/client"
	"github.com/intergral/deep/pkg/usagestats"
	"github.com/intergral/deep/pkg/util"
	"github.com/intergral/deep/pkg/util/log"
	"io"
	"net/http"
	"sort"

	"github.com/go-kit/log/level"
	"github.com/gorilla/mux"
	"github.com/grafana/dskit/grpcutil"
	"github.com/grafana/dskit/kv/memberlist"
	"github.com/grafana/dskit/modules"
	"github.com/grafana/dskit/ring"
	"github.com/grafana/dskit/services"
	frontend_v1 "github.com/intergral/deep/modules/frontend/v1"
	"github.com/intergral/deep/modules/generator"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/version"
	"github.com/weaveworks/common/middleware"
	"github.com/weaveworks/common/server"
	"github.com/weaveworks/common/signals"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"gopkg.in/yaml.v3"
)

const metricsNamespace = "deep"
const apiDocs = "https://intergral.com/docs/deep/latest/api_docs/"

var (
	metricConfigFeatDesc = prometheus.NewDesc(
		"deep_feature_enabled",
		"Boolean for configuration variables",
		[]string{"feature"},
		nil,
	)

	statFeatureEnabledAuth         = usagestats.NewInt("feature_enabled_auth_stats")
	statFeatureEnabledMultitenancy = usagestats.NewInt("feature_enabled_multitenancy")
)

// App is the root datastructure.
type App struct {
	cfg Config

	Server         *server.Server
	InternalServer *server.Server
	ring           *ring.Ring
	generatorRing  *ring.Ring
	overrides      *overrides.Overrides
	distributor    *distributor.Distributor
	querier        *querier.Querier
	frontend       *frontend_v1.Frontend
	compactor      *compactor.Compactor
	ingester       *ingester.Ingester
	generator      *generator.Generator
	store          storage.Store
	usageReport    *usagestats.Reporter
	MemberlistKV   *memberlist.KVInitService

	HTTPAuthMiddleware middleware.Interface

	ModuleManager *modules.Manager
	serviceMap    map[string]services.Service
	deps          map[string][]string

	// tracepoint config services
	tpClient          *client.TPClient
	tpRing            *ring.Ring
	tracepointService *tracepoint.TPService
	tracepointAPI     *tpapi.TracepointAPI
}

// New makes a new app.
func New(cfg Config) (*App, error) {
	app := &App{
		cfg: cfg,
	}

	usagestats.Edition("oss")

	statFeatureEnabledAuth.Set(0)
	if cfg.AuthEnabled {
		statFeatureEnabledAuth.Set(1)
	}

	statFeatureEnabledMultitenancy.Set(0)
	if cfg.MultitenancyEnabled {
		statFeatureEnabledMultitenancy.Set(1)
	}

	app.setupAuthMiddleware()

	if err := app.setupModuleManager(); err != nil {
		return nil, fmt.Errorf("failed to setup module manager %w", err)
	}

	return app, nil
}

func (t *App) setupAuthMiddleware() {
	if t.cfg.MultitenancyIsEnabled() {

		// don't check auth for these gRPC methods, since single call is used for multiple users
		noGRPCAuthOn := []string{
			"/frontend.Frontend/Process",
			"/frontend.Frontend/NotifyClientShutdown",
		}
		ignoredMethods := map[string]bool{}
		for _, m := range noGRPCAuthOn {
			ignoredMethods[m] = true
		}

		t.cfg.Server.GRPCMiddleware = []grpc.UnaryServerInterceptor{
			func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
				if ignoredMethods[info.FullMethod] {
					return handler(ctx, req)
				}
				return middleware.ServerUserHeaderInterceptor(ctx, req, info, handler)
			},
		}
		t.cfg.Server.GRPCStreamMiddleware = []grpc.StreamServerInterceptor{
			func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
				if ignoredMethods[info.FullMethod] {
					return handler(srv, ss)
				}
				return middleware.StreamServerUserHeaderInterceptor(srv, ss, info, handler)
			},
		}
		t.HTTPAuthMiddleware = middleware.AuthenticateUser
	} else {
		t.cfg.Server.GRPCMiddleware = []grpc.UnaryServerInterceptor{
			fakeGRPCAuthUniaryMiddleware,
		}
		t.cfg.Server.GRPCStreamMiddleware = []grpc.StreamServerInterceptor{
			fakeGRPCAuthStreamMiddleware,
		}
		t.HTTPAuthMiddleware = fakeHTTPAuthMiddleware
	}
}

// Run starts, and blocks until a signal is received.
func (t *App) Run() error {
	if !t.ModuleManager.IsUserVisibleModule(t.cfg.Target) {
		level.Warn(log.Logger).Log("msg", "selected target is an internal module, is this intended?", "target", t.cfg.Target)
	}

	serviceMap, err := t.ModuleManager.InitModuleServices(t.cfg.Target)
	if err != nil {
		return fmt.Errorf("failed to init module services %w", err)
	}
	t.serviceMap = serviceMap

	servs := []services.Service(nil)
	for _, s := range serviceMap {
		servs = append(servs, s)
	}

	sm, err := services.NewManager(servs...)
	if err != nil {
		return fmt.Errorf("failed to start service manager %w", err)
	}

	// before starting servers, register /ready handler and gRPC health check service.
	if t.cfg.InternalServer.Enable {
		t.InternalServer.HTTP.Path("/ready").Methods("GET").Handler(t.readyHandler(sm))
	}

	t.Server.HTTP.Path("/ready").Handler(t.readyHandler(sm))
	t.Server.HTTP.Path("/status").Handler(t.statusHandler()).Methods("GET")
	t.Server.HTTP.Path("/status/{endpoint}").Handler(t.statusHandler()).Methods("GET")
	grpc_health_v1.RegisterHealthServer(t.Server.GRPC, grpcutil.NewHealthCheck(sm))

	// Let's listen for events from this manager, and log them.
	healthy := func() { level.Info(log.Logger).Log("msg", "Deep started") }
	stopped := func() { level.Info(log.Logger).Log("msg", "Deep stopped") }
	serviceFailed := func(service services.Service) {
		// if any service fails, stop everything
		sm.StopAsync()

		// let's find out which module failed
		for m, s := range serviceMap {
			if s == service {
				switch service.FailureCase() {
				case modules.ErrStopProcess:
					level.Info(log.Logger).Log("msg", "received stop signal via return error", "module", m, "err", service.FailureCase())
				case context.Canceled:
				default:
					level.Error(log.Logger).Log("msg", "module failed", "module", m, "err", service.FailureCase())
				}

				return
			}
		}

		level.Error(log.Logger).Log("msg", "module failed", "module", "unknown", "err", service.FailureCase())
	}
	sm.AddListener(services.NewManagerListener(healthy, stopped, serviceFailed))

	// Setup signal handler. If signal arrives, we stop the manager, which stops all the services.
	handler := signals.NewHandler(t.Server.Log)
	go func() {
		handler.Loop()
		sm.StopAsync()
	}()

	// Start all services. This can really only fail if some service is already
	// in other state than New, which should not be the case.
	err = sm.StartAsync(context.Background())
	if err != nil {
		return fmt.Errorf("failed to start service manager %w", err)
	}

	return sm.AwaitStopped(context.Background())
}

func (t *App) writeStatusVersion(w io.Writer) error {
	_, err := w.Write([]byte(version.Print("deep") + "\n"))
	if err != nil {
		return err
	}

	return nil
}

func (t *App) writeStatusConfig(w io.Writer, r *http.Request) error {
	var output interface{}

	mode := r.URL.Query().Get("mode")
	switch mode {
	case "diff":
		defaultCfg := newDefaultConfig()

		defaultCfgYaml, err := util.YAMLMarshalUnmarshal(defaultCfg)
		if err != nil {
			return err
		}

		cfgYaml, err := util.YAMLMarshalUnmarshal(t.cfg)
		if err != nil {
			return err
		}

		output, err = util.DiffConfig(defaultCfgYaml, cfgYaml)
		if err != nil {
			return err
		}
	case "defaults":
		output = newDefaultConfig()
	case "":
		output = t.cfg
	default:
		return errors.Errorf("unknown value for mode query parameter: %v", mode)
	}

	out, err := yaml.Marshal(output)
	if err != nil {
		return err
	}

	_, err = w.Write([]byte("---\n"))
	if err != nil {
		return err
	}

	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func (t *App) readyHandler(sm *services.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !sm.IsHealthy() {
			msg := bytes.Buffer{}
			msg.WriteString("Some services are not Running:\n")

			byState := sm.ServicesByState()
			for st, ls := range byState {
				msg.WriteString(fmt.Sprintf("%v: %d\n", st, len(ls)))
			}

			http.Error(w, msg.String(), http.StatusServiceUnavailable)
			return
		}

		// Ingester has a special check that makes sure that it was able to register into the ring,
		// and that all other ring entries are OK too.
		if t.ingester != nil {
			if err := t.ingester.CheckReady(r.Context()); err != nil {
				http.Error(w, "Ingester not ready: "+err.Error(), http.StatusServiceUnavailable)
				return
			}
		}

		// Generator has a special check that makes sure that it was able to register into the ring,
		// and that all other ring entries are OK too.
		if t.generator != nil {
			if err := t.generator.CheckReady(r.Context()); err != nil {
				http.Error(w, "Generator not ready: "+err.Error(), http.StatusServiceUnavailable)
				return
			}
		}

		// Query Frontend has a special check that makes sure that a querier is attached before it signals
		// itself as ready
		if t.frontend != nil {
			if err := t.frontend.CheckReady(r.Context()); err != nil {
				http.Error(w, "Query Frontend not ready: "+err.Error(), http.StatusServiceUnavailable)
				return
			}
		}

		http.Error(w, "ready", http.StatusOK)
	}
}

func (t *App) writeRuntimeConfig(w io.Writer, r *http.Request) error {
	// Querier and query-frontend services do not run the overrides module
	if t.overrides == nil {
		_, err := w.Write([]byte(fmt.Sprintf("overrides module not loaded in %s\n", t.cfg.Target)))
		return err
	}
	return t.overrides.WriteStatusRuntimeConfig(w, r)
}

func (t *App) statusHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var errs []error
		msg := bytes.Buffer{}

		simpleEndpoints := map[string]func(io.Writer) error{
			"version":   t.writeStatusVersion,
			"services":  t.writeStatusServices,
			"endpoints": t.writeStatusEndpoints,
		}

		wrapStatus := func(endpoint string) {
			msg.WriteString("GET /status/" + endpoint + "\n")

			switch endpoint {
			case "runtime_config":
				err := t.writeRuntimeConfig(&msg, r)
				if err != nil {
					errs = append(errs, err)
				}
			case "config":
				err := t.writeStatusConfig(&msg, r)
				if err != nil {
					errs = append(errs, err)
				}
			default:
				err := simpleEndpoints[endpoint](&msg)
				if err != nil {
					errs = append(errs, err)
				}
			}
		}

		vars := mux.Vars(r)

		if endpoint, ok := vars["endpoint"]; ok {
			wrapStatus(endpoint)
		} else {
			wrapStatus("version")
			wrapStatus("services")
			wrapStatus("endpoints")
			wrapStatus("runtime_config")
			wrapStatus("config")
		}

		w.Header().Set("Content-Type", "text/plain")

		joinErrors := func(errs []error) error {
			if len(errs) == 0 {
				return nil
			}
			var err error

			for _, e := range errs {
				if e != nil {
					if err == nil {
						err = e
					} else {
						err = errors.Wrap(err, e.Error())
					}
				}
			}
			return err
		}

		err := joinErrors(errs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		if _, err := w.Write(msg.Bytes()); err != nil {
			level.Error(log.Logger).Log("msg", "error writing response", "err", err)
		}
	}
}

func (t *App) writeStatusServices(w io.Writer) error {
	svcNames := make([]string, 0, len(t.serviceMap))
	for name := range t.serviceMap {
		svcNames = append(svcNames, name)
	}

	sort.Strings(svcNames)

	x := table.NewWriter()
	x.SetOutputMirror(w)
	x.AppendHeader(table.Row{"service name", "status", "failure case"})

	for _, name := range svcNames {
		service := t.serviceMap[name]

		var e string

		if err := service.FailureCase(); err != nil {
			e = err.Error()
		}

		x.AppendRows([]table.Row{
			{name, service.State(), e},
		})
	}

	x.AppendSeparator()
	x.Render()

	return nil
}

func (t *App) writeStatusEndpoints(w io.Writer) error {
	type endpoint struct {
		name  string
		regex string
	}

	endpoints := []endpoint{}

	err := t.Server.HTTP.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		e := endpoint{}

		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			e.name = pathTemplate
		}

		pathRegexp, err := route.GetPathRegexp()
		if err == nil {
			e.regex = pathRegexp
		}

		endpoints = append(endpoints, e)

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "error walking routes")
	}

	sort.Slice(endpoints[:], func(i, j int) bool {
		return endpoints[i].name < endpoints[j].name
	})

	x := table.NewWriter()
	x.SetOutputMirror(w)
	x.AppendHeader(table.Row{"name", "regex"})

	for _, e := range endpoints {
		x.AppendRows([]table.Row{
			{e.name, e.regex},
		})
	}

	x.AppendSeparator()
	x.Render()

	_, err = w.Write([]byte(fmt.Sprintf("\nAPI documentation: %s\n\n", apiDocs)))
	if err != nil {
		return errors.Wrap(err, "error writing status endpoints")
	}

	return nil
}
