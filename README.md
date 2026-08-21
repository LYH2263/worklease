# go-worklease

带租约的工作队列库 + `workd` 管理服务：Enqueue、Claim、Heartbeat、Ack/Nack、租约过期回收。

## Build / Test

```bash
go build ./...
go test ./... -count=1
```

## Daemon

```bash
go run ./cmd/workd -addr :8107 -web web
```
