package model

// Point 表示一个时序点位。
type Point struct {
	Series string
	Seq    int64
	Value  float64
	Ts     int64
}

// BucketKey 定位一个时间桶。
type BucketKey struct {
	Series string
	Slot   int64
}

// AggResult 是时间桶聚合结果。
type AggResult struct {
	Count int64
	Sum   float64
	Max   float64
	Min   float64
}

// RollupResult 是降采样后的粗粒度结果。
type RollupResult struct {
	Slot  int64
	Count int64
	Sum   float64
	Max   float64
	Min   float64
}

// BucketResult 带槽位的聚合结果。
type BucketResult struct {
	Slot int64
	Agg  AggResult
}

// Rule 是一条告警规则。
type Rule struct {
	Name     string
	Metric   string
	Op       string
	Threshold float64
	Window   int64
	Webhooks []string
}

// RuleState 记录规则当前状态。
type RuleState struct {
	Status string // ok / alert
	Since  int64
}

// Subscription 是一个通知订阅者。
type Subscription struct {
	ID      string
	Webhook string
	Rules   []string
}

// Worker 是摄取 worker 的注册信息。
type Worker struct {
	ID       string
	LastSeen int64
}

// Event 是一次告警通知事件。
type Event struct {
	Rule   string
	Series string
	Status string
	Value  float64
}

// StatusReport 是状态接口返回的汇总。
type StatusReport struct {
	Workers      int
	Rules        int
	Subscribers  int
	Pending      int
	Ingested     int64
	Dropped      int64
	Rejected     int64
	Alerts       int64
	Events       int64
	KnownSeries  int
	RouteVersion int
	Buckets      int
	ClosedBuckets int
}
