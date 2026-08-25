package verifycase

import (
	"testing"

	"github.com/dyl-06/telemetry/internal/downsample"
	"github.com/dyl-06/telemetry/internal/shard"
	"github.com/dyl-06/telemetry/internal/store"
)

func TestMissingSeriesBucketNoPanic(t *testing.T) {
	router := shard.New([]string{"shard-a"})
	st := store.New(router)
	reader := downsample.NewReader(st)
	if _, err := reader.ReadBucket("ghost.series", 1); err == nil {
		t.Fatal("expected error for missing series bucket")
	}
}
