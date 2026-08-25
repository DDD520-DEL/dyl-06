package verifycase

import (
	"testing"

	"github.com/dyl-06/telemetry/internal/ingest"
	"github.com/dyl-06/telemetry/internal/metrics"
	"github.com/dyl-06/telemetry/internal/seq"
	"github.com/dyl-06/telemetry/internal/shard"
	"github.com/dyl-06/telemetry/internal/store"
)

func TestMalformedPointRejectedNotMerged(t *testing.T) {
	router := shard.New([]string{"shard-a"})
	st := store.New(router)
	pipeline := ingest.NewPipeline(st, seq.New(), metrics.New())
	raw := []byte(`{"series":"cpu.usage","seq":1,"value":"not-a-number","ts":60}`)
	if _, err := pipeline.HandleRaw(raw); err == nil {
		t.Fatal("expected decode error for malformed point")
	}
	if _, ok := st.LastPoint("cpu.usage"); ok {
		t.Fatal("malformed point was silently merged")
	}
}
