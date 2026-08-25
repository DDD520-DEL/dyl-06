基于 Go 实现的时序遥测指标平台项目，一款后端服务，完成指标写入、分片存储与聚合查询。

## 快速开始

```bash
go build -mod=vendor ./...
go test -mod=vendor ./...
```

## 接口

- POST /ingest 批量上报点位
- GET /query 查询序列聚合结果
- POST /subscribe /unsubscribe 告警通知订阅
- POST /rules 热更新告警规则
- POST /rebalance 分片重平衡
- POST /heartbeat worker 心跳
