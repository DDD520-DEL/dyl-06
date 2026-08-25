package metrics

// Snapshot 返回全部计数快照。
func (m *Metrics) Snapshot() map[string]int64 {
	return map[string]int64{
		"ingested": m.Ingested.Load(),
		"dropped":  m.Dropped.Load(),
		"rejected": m.Rejected.Load(),
		"alerts":   m.Alerts.Load(),
		"events":   m.Events.Load(),
	}
}
