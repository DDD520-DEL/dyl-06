package verifycase

import (
	"sync"
	"testing"

	"github.com/dyl-06/telemetry/internal/agg"
	"github.com/dyl-06/telemetry/internal/model"
)

func TestConcurrentBucketMergeNoLostUpdate(t *testing.T) {
	bucket := agg.NewBucket()
	const workers = 16
	const perWorker = 100
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < perWorker; j++ {
				bucket.Merge(model.Point{Series: "latency", Seq: int64(j), Value: 1, Ts: 60})
			}
		}()
	}
	wg.Wait()
	got := bucket.Read()
	if got.Count != workers*perWorker {
		t.Fatalf("count = %d, want %d", got.Count, workers*perWorker)
	}
}
