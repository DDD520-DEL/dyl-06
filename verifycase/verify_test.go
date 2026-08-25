package verifycase

import (
	"testing"

	"github.com/dyl-06/telemetry/internal/ingest"
	"github.com/dyl-06/telemetry/internal/model"
	"github.com/dyl-06/telemetry/internal/shard"
	"github.com/dyl-06/telemetry/internal/store"
)

func TestBatchIngestAllOrNothing(t *testing.T) {
	router := shard.New([]string{"shard-a"})
	st := store.New(router)
	points := []model.Point{
		{Series: "cpu.usage", Seq: 1, Value: 12.5, Ts: 60},
		{Series: "", Seq: 2, Value: 3, Ts: 60},
	}
	if err := ingest.Batch(st, points); err == nil {
		t.Fatal("expected batch validation error")
	}
	if _, ok := st.LastPoint("cpu.usage"); ok {
		t.Fatal("partial batch write detected: first point persisted")
	}
}
