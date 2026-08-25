package verifycase

import (
	"testing"

	"github.com/dyl-06/telemetry/internal/api"
	"github.com/dyl-06/telemetry/internal/model"
	"github.com/dyl-06/telemetry/internal/shard"
	"github.com/dyl-06/telemetry/internal/store"
)

func TestRebalanceRoutesNewWrites(t *testing.T) {
	router := shard.New([]string{"shard-a"})
	st := store.New(router)
	st.WritePoint(model.Point{Series: "mem.used", Seq: 1, Value: 100, Ts: 60})
	if got := st.ShardFor("mem.used"); got != "shard-a" {
		t.Fatalf("initial route = %q, want shard-a", got)
	}
	if err := api.RebalanceRouter(router, st, []string{"shard-b"}); err != nil {
		t.Fatal(err)
	}
	st.WritePoint(model.Point{Series: "mem.used", Seq: 2, Value: 110, Ts: 120})
	if got := st.ShardFor("mem.used"); got != "shard-b" {
		t.Fatalf("route after rebalance = %q, want shard-b", got)
	}
}
