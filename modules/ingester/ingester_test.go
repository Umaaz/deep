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

package ingester

import (
	"context"
	"github.com/go-kit/log"
	"github.com/grafana/dskit/flagext"
	"github.com/grafana/dskit/kv/consul"
	"github.com/grafana/dskit/ring"
	"github.com/intergral/deep/modules/overrides"
	"github.com/intergral/deep/modules/storage"
	"github.com/intergral/deep/pkg/deepdb"
	"github.com/intergral/deep/pkg/deepdb/backend"
	"github.com/intergral/deep/pkg/deepdb/backend/local"
	"github.com/intergral/deep/pkg/deepdb/encoding"
	"github.com/intergral/deep/pkg/deepdb/encoding/common"
	"github.com/intergral/deep/pkg/deepdb/wal"
	"github.com/intergral/deep/pkg/deeppb"
	tp "github.com/intergral/deep/pkg/deeppb/tracepoint/v1"
	"github.com/intergral/deep/pkg/model"
	v1 "github.com/intergral/deep/pkg/model/v1"
	"github.com/intergral/deep/pkg/util"
	"github.com/intergral/deep/pkg/util/test"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

//
//import (
//	"context"
//	"crypto/rand"
//	"os"
//	"testing"
//	"time"
//
//	"github.com/go-kit/log"
//	"github.com/golang/protobuf/proto"
//	"github.com/grafana/dskit/flagext"
//	"github.com/grafana/dskit/kv/consul"
//	"github.com/grafana/dskit/ring"
//	"github.com/prometheus/client_golang/prometheus"
//	"github.com/stretchr/testify/require"
//	"github.com/weaveworks/common/user"
//
//	"github.com/intergral/deep/modules/overrides"
//	"github.com/intergral/deep/modules/storage"
//	"github.com/intergral/deep/pkg/deepdb"
//	"github.com/intergral/deep/pkg/deepdb/backend"
//	"github.com/intergral/deep/pkg/deepdb/backend/local"
//	"github.com/intergral/deep/pkg/deepdb/encoding"
//	"github.com/intergral/deep/pkg/deepdb/encoding/common"
//	"github.com/intergral/deep/pkg/deepdb/wal"
//	"github.com/intergral/deep/pkg/model"
//	"github.com/intergral/deep/pkg/model/trace"
//	model_v1 "github.com/intergral/deep/pkg/model/v1"
//	model_v2 "github.com/intergral/deep/pkg/model/v2"
//	"github.com/intergral/deep/pkg/tempopb"
//	v1 "github.com/intergral/deep/pkg/tempopb/trace/v1"
//	"github.com/intergral/deep/pkg/util/test"
//)

func TestPushQueryAllEncodings(t *testing.T) {
	for _, e := range model.AllEncodings {
		t.Run(e, func(t *testing.T) {
			tempDir := t.TempDir()
			ctx := util.InjectTenantID(context.Background(), "test")
			ingester, snapshots, snapshotsIDs := defaultIngestWithPush(t, tempDir, pushBatchV1)

			for pos, snapshotId := range snapshotsIDs {
				byId, err := ingester.FindSnapshotByID(ctx, &deeppb.SnapshotByIDRequest{
					ID: snapshotId,
				})
				require.NoError(t, err, "unexpect error querying")
				require.Equal(t, byId.Snapshot.ID, snapshots[pos].ID)
			}

			for _, tenantBlockManager := range ingester.instances {
				err := tenantBlockManager.CutSnapshots(0, true)
				require.NoError(t, err, "unexpected error cutting snapshots")
			}

			for pos, snapshotId := range snapshotsIDs {
				byId, err := ingester.FindSnapshotByID(ctx, &deeppb.SnapshotByIDRequest{
					ID: snapshotId,
				})
				require.NoError(t, err, "unexpect error querying")
				require.Equal(t, byId.Snapshot.ID, snapshots[pos].ID)
			}

		})
	}
}

//func TestPushQueryAllEncodings(t *testing.T) {
//	for _, e := range model.AllEncodings {
//		t.Run(e, func(t *testing.T) {
//			var push func(testing.TB, *Ingester, *v1.ResourceSpans, []byte)
//
//			switch e {
//			case model_v1.Encoding:
//				push = pushBatchV1
//			case model_v2.Encoding:
//				push = pushBatchV2
//			default:
//				t.Fatal("unsupported encoding", e)
//			}
//
//			tmpDir := t.TempDir()
//			ctx := util.InjectTenantID(context.Background(), "test")
//			ingester, traces, traceIDs := defaultIngesterWithPush(t, tmpDir, push)
//
//			// live trace search
//			for pos, snapshotId := range traceIDs {
//				foundTrace, err := ingester.FindTraceByID(ctx, &tempopb.TraceByIDRequest{
//					TraceID: snapshotId,
//				})
//				require.NoError(t, err, "unexpected error querying")
//				require.Equal(t, foundTrace.Trace, traces[pos])
//			}
//
//			// force cut all traces
//			for _, tenantBlockManager := range ingester.instances {
//				err := tenantBlockManager.CutSnapshots(0, true)
//				require.NoError(t, err, "unexpected error cutting traces")
//			}
//
//			// head block search
//			for i, snapshotId := range traceIDs {
//				foundTrace, err := ingester.FindTraceByID(ctx, &tempopb.TraceByIDRequest{
//					TraceID: snapshotId,
//				})
//				require.NoError(t, err, "unexpected error querying")
//				trace.SortTrace(foundTrace.Trace)
//				equal := proto.Equal(traces[i], foundTrace.Trace)
//				require.True(t, equal)
//			}
//		})
//	}
//}
//
//func TestFullTraceReturned(t *testing.T) {
//	ctx := util.InjectTenantID(context.Background(), "test")
//	ingester, _, _ := defaultIngester(t, t.TempDir())
//
//	snapshotId := make([]byte, 16)
//	_, err := rand.Read(snapshotId)
//	require.NoError(t, err)
//	testTrace := test.MakeTrace(2, snapshotId) // 2 batches
//	trace.SortTrace(testTrace)
//
//	// push the first batch
//	pushBatchV2(t, ingester, testTrace.Batches[0], snapshotId)
//
//	// force cut all traces
//	for _, tenantBlockManager := range ingester.instances {
//		err = tenantBlockManager.CutSnapshots(0, true)
//		require.NoError(t, err, "unexpected error cutting traces")
//	}
//
//	// push the 2nd batch
//	pushBatchV2(t, ingester, testTrace.Batches[1], snapshotId)
//
//	// make sure the trace comes back whole
//	foundTrace, err := ingester.FindTraceByID(ctx, &tempopb.TraceByIDRequest{
//		TraceID: snapshotId,
//	})
//	require.NoError(t, err, "unexpected error querying")
//	require.True(t, proto.Equal(testTrace, foundTrace.Trace))
//
//	// force cut all traces
//	for _, tenantBlockManager := range ingester.instances {
//		err = tenantBlockManager.CutSnapshots(0, true)
//		require.NoError(t, err, "unexpected error cutting traces")
//	}
//
//	// make sure the trace comes back whole
//	foundTrace, err = ingester.FindTraceByID(ctx, &tempopb.TraceByIDRequest{
//		TraceID: snapshotId,
//	})
//	require.NoError(t, err, "unexpected error querying")
//	require.True(t, proto.Equal(testTrace, foundTrace.Trace))
//}
//
//func TestWal(t *testing.T) {
//	tmpDir := t.TempDir()
//
//	ctx := util.InjectTenantID(context.Background(), "test")
//	ingester, traces, traceIDs := defaultIngester(t, tmpDir)
//
//	for pos, snapshotId := range traceIDs {
//		foundTrace, err := ingester.FindTraceByID(ctx, &tempopb.TraceByIDRequest{
//			TraceID: snapshotId,
//		})
//		require.NoError(t, err, "unexpected error querying")
//		require.Equal(t, foundTrace.Trace, traces[pos])
//	}
//
//	// force cut all traces
//	for _, tenantBlockManager := range ingester.instances {
//		err := tenantBlockManager.CutSnapshots(0, true)
//		require.NoError(t, err, "unexpected error cutting traces")
//	}
//
//	// create new ingester.  this should replay wal!
//	ingester, _, _ = defaultIngester(t, tmpDir)
//
//	// should be able to find old traces that were replayed
//	for i, snapshotId := range traceIDs {
//		foundTrace, err := ingester.FindTraceByID(ctx, &tempopb.TraceByIDRequest{
//			TraceID: snapshotId,
//		})
//		require.NoError(t, err, "unexpected error querying")
//		require.NotNil(t, foundTrace.Trace)
//		trace.SortTrace(foundTrace.Trace)
//		equal := proto.Equal(traces[i], foundTrace.Trace)
//		require.True(t, equal)
//	}
//
//	// a block that has been replayed should have a flush queue entry to complete it
//	// wait for the flush queues to be empty and then confirm there is a complete block
//	for !ingester.flushQueues.IsEmpty() {
//		time.Sleep(100 * time.Millisecond)
//	}
//
//	require.Len(t, ingester.instances["test"].completingBlocks, 0)
//	require.Len(t, ingester.instances["test"].completeBlocks, 1)
//
//	// should be able to find old traces that were replayed
//	for i, snapshotId := range traceIDs {
//		foundTrace, err := ingester.FindTraceByID(ctx, &tempopb.TraceByIDRequest{
//			TraceID: snapshotId,
//		})
//		require.NoError(t, err, "unexpected error querying")
//
//		trace.SortTrace(foundTrace.Trace)
//		test.TracesEqual(t, traces[i], foundTrace.Trace)
//	}
//}
//
//func TestWalDropsZeroLength(t *testing.T) {
//	tmpDir := t.TempDir()
//	ingester, _, _ := defaultIngester(t, tmpDir)
//
//	// force cut all traces and wipe wal
//	for _, tenantBlockManager := range ingester.instances {
//		err := tenantBlockManager.CutSnapshots(0, true)
//		require.NoError(t, err, "unexpected error cutting traces")
//
//		blockID, err := tenantBlockManager.CutBlockIfReady(0, 0, true)
//		require.NoError(t, err)
//
//		err = tenantBlockManager.CompleteBlock(blockID)
//		require.NoError(t, err)
//
//		err = tenantBlockManager.ClearCompletingBlock(blockID)
//		require.NoError(t, err)
//
//		err = ingester.local.ClearBlock(blockID, tenantBlockManager.instanceID)
//		require.NoError(t, err)
//	}
//
//	// create new ingester. we should have no tenants b/c we all our wals should have been 0 length
//	ingester, _, _ = defaultIngesterWithPush(t, tmpDir, func(t testing.TB, i *Ingester, rs *v1.ResourceSpans, b []byte) {})
//	require.Equal(t, 0, len(ingester.instances))
//}
//
//func TestSearchWAL(t *testing.T) {
//	tmpDir, err := os.MkdirTemp("/tmp", "")
//	require.NoError(t, err, "unexpected error getting tempdir")
//	defer os.RemoveAll(tmpDir)
//
//	i := defaultIngesterModule(t, tmpDir)
//	inst, _ := i.getOrCreateInstance("test")
//	require.NotNil(t, inst)
//
//	dec := model.MustNewSegmentDecoder(model.CurrentEncoding)
//
//	// create some search data
//	id := make([]byte, 16)
//	_, err = rand.Read(id)
//	require.NoError(t, err)
//	trace := test.MakeTrace(10, id)
//	b1, err := dec.PrepareForWrite(trace, 0, 0)
//	require.NoError(t, err)
//
//	// push to tenantBlockManager
//	require.NoError(t, inst.PushBytes(context.Background(), id, b1))
//
//	// Write wal
//	require.NoError(t, inst.CutSnapshots(0, true))
//
//	// search WAL
//	ctx := util.InjectTenantID(context.Background(), "test")
//	searchReq := &tempopb.SearchRequest{Tags: map[string]string{
//		"foo": "bar",
//	}}
//	results, err := inst.Search(ctx, searchReq)
//	require.NoError(t, err)
//	require.Equal(t, uint32(1), results.Metrics.InspectedTraces)
//
//	// Shutdown
//	require.NoError(t, i.stopping(nil))
//
//	// the below tests sometimes fail in CI. Awful bandaid sleep slammed in here
//	// to reduce occurrence of this issue. todo: fix it properly
//	time.Sleep(500 * time.Millisecond)
//
//	// replay wal
//	i = defaultIngesterModule(t, tmpDir)
//	inst, ok := i.getInstanceByID("test")
//	require.True(t, ok)
//
//	results, err = inst.Search(ctx, searchReq)
//	require.NoError(t, err)
//	require.Equal(t, uint32(1), results.Metrics.InspectedTraces)
//}
//
//// TODO - This test is flaky and commented out until it's fixed
//// TestWalReplayDeletesLocalBlocks simulates the condition where an ingester restarts after a wal is completed
//// to the local disk, but before the wal is deleted. On startup both blocks exist, and the ingester now errs
//// on the side of caution and chooses to replay the wal instead of rediscovering the local block.
///*func TestWalReplayDeletesLocalBlocks(t *testing.T) {
//	tmpDir := t.TempDir()
//
//	i, _, _ := defaultIngester(t, tmpDir)
//	inst, ok := i.getInstanceByID("test")
//	require.True(t, ok)
//
//	// Write wal
//	err := inst.CutSnapshots(0, true)
//	require.NoError(t, err)
//	blockID, err := inst.CutBlockIfReady(0, 0, true)
//	require.NoError(t, err)
//
//	// Complete block
//	err = inst.CompleteBlock(blockID)
//	require.NoError(t, err)
//
//	// Shutdown
//	err = i.stopping(nil)
//	require.NoError(t, err)
//
//	// At this point both wal and complete block exist.
//	// The wal still exists because we manually completed it and
//	// didn't delete it with ClearCompletingBlock()
//	require.Len(t, inst.completingBlocks, 1)
//	require.Len(t, inst.completeBlocks, 1)
//	require.Equal(t, blockID, inst.completingBlocks[0].BlockID())
//	require.Equal(t, blockID, inst.completeBlocks[0].BlockMeta().BlockID)
//
//	// Simulate a restart by creating a new ingester.
//	i, _, _ = defaultIngester(t, tmpDir)
//	inst, ok = i.getInstanceByID("test")
//	require.True(t, ok)
//
//	// After restart we only have the 1 wal block
//	// TODO - fix race conditions here around access inst fields outside of mutex
//	require.Len(t, inst.completingBlocks, 1)
//	require.Len(t, inst.completeBlocks, 0)
//	require.Equal(t, blockID, inst.completingBlocks[0].BlockID())
//
//	// Shutdown, cleanup
//	err = i.stopping(nil)
//	require.NoError(t, err)
//}
//*/

//func TestFlush(t *testing.T) {
//	tmpDir := t.TempDir()
//
//	ctx := util.InjectTenantID(context.Background(), "test")
//	ingester, traces, traceIDs := defaultIngester(t, tmpDir)
//
//	for pos, snapshotId := range traceIDs {
//		foundTrace, err := ingester.FindTraceByID(ctx, &tempopb.TraceByIDRequest{
//			TraceID: snapshotId,
//		})
//		require.NoError(t, err, "unexpected error querying")
//		require.Equal(t, traces[pos], foundTrace.Trace)
//	}
//
//	// stopping the ingester should force cut all live traces to disk
//	require.NoError(t, ingester.stopping(nil))
//
//	// create new ingester.  this should replay wal!
//	ingester, _, _ = defaultIngester(t, tmpDir)
//
//	// should be able to find old traces that were replayed
//	for i, snapshotId := range traceIDs {
//		foundTrace, err := ingester.FindTraceByID(ctx, &tempopb.TraceByIDRequest{
//			TraceID: snapshotId,
//		})
//		require.NoError(t, err, "unexpected error querying")
//		trace.SortTrace(foundTrace.Trace)
//		equal := proto.Equal(traces[i], foundTrace.Trace)
//		require.True(t, equal)
//	}
//}

func defaultIngesterModule(t testing.TB, tmpDir string) *Ingester {
	ingesterConfig := defaultIngesterTestConfig()
	limits, err := overrides.NewOverrides(defaultLimitsTestConfig())
	require.NoError(t, err, "unexpected error creating overrides")

	s, err := storage.NewStore(storage.Config{
		TracePoint: deepdb.Config{
			Backend: "local",
			Local: &local.Config{
				Path: tmpDir,
			},
			Block: &common.BlockConfig{
				BloomFP:             0.01,
				BloomShardSizeBytes: 100_000,
				Version:             encoding.DefaultEncoding().Version(),
				Encoding:            backend.EncLZ4_1M,
			},
			WAL: &wal.Config{
				Filepath: tmpDir,
			},
		},
	}, log.NewNopLogger())
	require.NoError(t, err, "unexpected error store")

	ingester, err := New(ingesterConfig, s, limits, prometheus.NewPedanticRegistry())
	require.NoError(t, err, "unexpected error creating ingester")
	ingester.replayJitter = false

	err = ingester.starting(context.Background())
	require.NoError(t, err, "unexpected error starting ingester")

	return ingester
}

//
//func defaultIngester(t testing.TB, tmpDir string) (*Ingester, []*tempopb.Trace, [][]byte) {
//	return defaultIngesterWithPush(t, tmpDir, pushBatchV2)
//}

func defaultIngestWithPush(t testing.TB, tmpDir string, push func(testing.TB, *Ingester, *tp.Snapshot, []byte)) (*Ingester, []*tp.Snapshot, [][]byte) {
	ingester := defaultIngesterModule(t, tmpDir)

	// make some fake traceIDs/requests
	snapshots := make([]*tp.Snapshot, 0)

	ids := make([][]byte, 0)
	for i := 0; i < 10; i++ {
		testSnapshot := test.GenerateSnapshot(i, nil)

		snapshots = append(snapshots, testSnapshot)
		ids = append(ids, testSnapshot.ID)
	}

	for i, snapshot := range snapshots {
		push(t, ingester, snapshot, ids[i])
	}

	return ingester, snapshots, ids
}

func defaultIngesterTestConfig() Config {
	cfg := Config{}

	flagext.DefaultValues(&cfg.LifecyclerConfig)
	mockStore, _ := consul.NewInMemoryClient(
		ring.GetCodec(),
		log.NewNopLogger(),
		nil,
	)

	cfg.FlushCheckPeriod = 99999 * time.Hour
	cfg.MaxSnapshotIdle = 99999 * time.Hour
	cfg.ConcurrentFlushes = 1
	cfg.LifecyclerConfig.RingConfig.KVStore.Mock = mockStore
	cfg.LifecyclerConfig.NumTokens = 1
	cfg.LifecyclerConfig.ListenPort = 0
	cfg.LifecyclerConfig.Addr = "localhost"
	cfg.LifecyclerConfig.ID = "localhost"
	cfg.LifecyclerConfig.FinalSleep = 0

	return cfg
}

func defaultLimitsTestConfig() overrides.Limits {
	limits := overrides.Limits{}
	flagext.DefaultValues(&limits)
	return limits
}

func pushBatchV1(t testing.TB, i *Ingester, snapshot *tp.Snapshot, _ []byte) {
	ctx := util.InjectTenantID(context.Background(), "test")

	batchDecoder := model.MustNewSegmentDecoder(v1.Encoding)

	buffer, err := batchDecoder.PrepareForWrite(snapshot, uint32(snapshot.TsNanos))
	require.NoError(t, err)

	_, err = i.PushBytes(ctx, &deeppb.PushBytesRequest{
		Snapshot: buffer,
		ID:       snapshot.ID,
	})
	require.NoError(t, err)
}
