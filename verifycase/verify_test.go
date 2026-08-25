package verifycase

import (
	"testing"

	"github.com/dyl-06/telemetry/internal/downsample"
	"github.com/dyl-06/telemetry/internal/model"
	"github.com/dyl-06/telemetry/internal/shard"
	"github.com/dyl-06/telemetry/internal/store"
)

func TestClosedBucketIgnoresLatePoints(t *testing.T) {
	router := shard.New([]string{"shard-a"})
	st := store.New(router)
	manager := downsample.New(st, 60)
	st.WritePoint(model.Point{Series: "disk.io", Seq: 1, Value: 10, Ts: 60})
	manager.CloseBucket("disk.io", 1)
	st.WritePoint(model.Point{Series: "disk.io", Seq: 2, Value: 20, Ts: 60})
	bucket := st.GetBucket("disk.io", 1)
	if bucket == nil {
		t.Fatal("bucket missing")
	}
	if got := bucket.Read().Count; got != 1 {
		t.Fatalf("closed bucket count = %d, want 1 (late point re-merged)", got)
	}
}
