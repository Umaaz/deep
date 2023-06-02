/*
 * Copyright (C) 2023  Intergral GmbH
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package distributor

//
//import (
//	"context"
//	"flag"
//	"sync"
//	"testing"
//	"time"
//
//	"github.com/go-kit/log"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//
//	"github.com/intergral/deep/modules/overrides"
//	v1 "github.com/intergral/deep/pkg/deeppb/trace/v1"
//	"github.com/intergral/deep/pkg/util"
//	"github.com/intergral/deep/pkg/util/test"
//)
//
//const tenantID = "tenant-id"
//
//func TestForwarder(t *testing.T) {
//	oCfg := overrides.Limits{}
//	oCfg.RegisterFlags(&flag.FlagSet{})
//
//	id, err := util.HexStringToTraceID("1234567890abcdef")
//	require.NoError(t, err)
//
//	b := test.MakeBatch(10, id)
//	keys, rebatchedTraces, err := requestsByTraceID([]*v1.ResourceSpans{b}, tenantID, 10)
//	require.NoError(t, err)
//
//	o, err := overrides.NewOverrides(oCfg)
//	require.NoError(t, err)
//
//	wg := sync.WaitGroup{}
//	f := newGeneratorForwarder(
//		log.NewNopLogger(),
//		func(ctx context.Context, userID string, k []uint32, traces []*rebatchedTrace) error {
//			assert.Equal(t, tenantID, userID)
//			assert.Equal(t, keys, k)
//			assert.Equal(t, rebatchedTraces, traces)
//			wg.Done()
//			return nil
//		},
//		o,
//	)
//	require.NoError(t, f.start(context.Background()))
//	defer func() {
//		require.NoError(t, f.stop(nil))
//	}()
//
//	wg.Add(1)
//	f.SendTraces(context.Background(), tenantID, keys, rebatchedTraces)
//	wg.Wait()
//
//	wg.Add(1)
//	f.SendTraces(context.Background(), tenantID, keys, rebatchedTraces)
//	wg.Wait()
//}
//
//func TestForwarder_shutdown(t *testing.T) {
//	oCfg := overrides.Limits{}
//	oCfg.RegisterFlags(&flag.FlagSet{})
//	oCfg.MetricsGeneratorForwarderQueueSize = 200
//
//	id, err := util.HexStringToTraceID("1234567890abcdef")
//	require.NoError(t, err)
//
//	b := test.MakeBatch(10, id)
//	keys, rebatchedTraces, err := requestsByTraceID([]*v1.ResourceSpans{b}, tenantID, 10)
//	require.NoError(t, err)
//
//	o, err := overrides.NewOverrides(oCfg)
//	require.NoError(t, err)
//
//	signalCh := make(chan struct{})
//	f := newGeneratorForwarder(
//		log.NewNopLogger(),
//		func(ctx context.Context, userID string, k []uint32, traces []*rebatchedTrace) error {
//			<-signalCh
//
//			assert.Equal(t, tenantID, userID)
//			assert.Equal(t, keys, k)
//			assert.Equal(t, rebatchedTraces, traces)
//			return nil
//		},
//		o,
//	)
//
//	require.NoError(t, f.start(context.Background()))
//	defer func() {
//		go func() {
//			// Wait to unblock processing of requests so shutdown and draining the queue is done in parallel
//			time.Sleep(time.Second)
//			close(signalCh)
//		}()
//		require.NoError(t, f.stop(nil))
//	}()
//
//	for i := 0; i < 100; i++ {
//		f.SendTraces(context.Background(), tenantID, keys, rebatchedTraces)
//	}
//}
