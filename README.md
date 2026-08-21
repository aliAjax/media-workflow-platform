# Media Workflow Platform

纯 Go 媒体资产与异步 DAG 处理服务，覆盖资产登记、Pipeline 校验、任务调度、内置安全 mock executor、产物 manifest、原子 JSON 持久化和本地 blob 存储。

## 快速开始

```bash
go test ./...
go run ./cmd/media-api
curl http://127.0.0.1:8084/healthz
```

API 可创建资产、Pipeline 和 Job，worker 异步执行 `probe`、`transcode`、`thumbnail`、`waveform`、`subtitle`、`package`、`quality-check` 阶段。通过 `GET /api/v1/jobs/{id}` 查询状态。真实编码器、PostgreSQL 和对象存储通过领域端口替换，本地模式不需要外部依赖。

分层为 `transport -> application -> domain ports -> repository/storage`。生产环境应启用 TLS、数据库事务、独立 worker、资源隔离与 Secret 管理。详见 `docs/operations.md`。
