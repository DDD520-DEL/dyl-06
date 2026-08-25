package verifycase

import (
	"testing"

	"github.com/dyl-06/telemetry/internal/model"
	"github.com/dyl-06/telemetry/internal/seq"
	"github.com/dyl-06/telemetry/internal/shard"
	"github.com/dyl-06/telemetry/internal/store"
)

func TestLatePointDoesNotOverwriteNewer(t *testing.T) {
	router := shard.New([]string{"shard-a"})
	st := store.New(router)
	guard := seq.New()
	st.WritePointWithGuard(guard, model.Point{Series: "net.bytes", Seq: 5, Value: 500, Ts: 60})
	st.WritePointWithGuard(guard, model.Point{Series: "net.bytes", Seq: 3, Value: 300, Ts: 60})
	point, ok := st.LastPoint("net.bytes")
	if !ok {
		t.Fatal("series missing")
	}
	if point.Value != 500 {
		t.Fatalf("value = %v, want 500 (late point overwrote newer)", point.Value)
	}
}
