# Load Balancer

纯 Go 标准库实现的负载均衡管理后端。除节点、策略和路由管理 API 外，服务还包含并发健康探测、不可变路由快照、请求取消重试、版本化状态回写、请求作用域审计和运行配置构建能力。

## 运行

```bash
cd origin/
go run ./cmd/server
```

## API

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/backends | 创建后端节点 |
| GET | /api/backends | 后端节点列表（支持 status、keyword 过滤） |
| GET | /api/backends/{id} | 获取后端节点详情 |
| PUT | /api/backends/{id} | 更新后端节点 |
| DELETE | /api/backends/{id} | 删除后端节点 |
| POST | /api/backends/{id}/status | 状态流转（up/down/draining） |
| POST | /api/backends/{id}/health-check | 健康检查上报 |
| POST | /api/strategies | 创建负载均衡策略 |
| GET | /api/strategies | 策略列表（支持 algorithm、keyword 过滤） |
| GET | /api/strategies/{id} | 获取策略详情 |
| PUT | /api/strategies/{id} | 更新策略 |
| DELETE | /api/strategies/{id} | 删除策略 |
| POST | /api/route-rules | 创建路由规则 |
| GET | /api/route-rules | 路由规则列表（支持 path_prefix、keyword 过滤） |
| GET | /api/route-rules/{id} | 获取路由规则详情 |
| PUT | /api/route-rules/{id} | 更新路由规则 |
| DELETE | /api/route-rules/{id} | 删除路由规则 |
| POST | /api/traffic-stats | 创建流量统计 |
| GET | /api/traffic-stats | 流量统计列表（支持 backend_id 过滤） |
| GET | /api/traffic-stats/{id} | 获取流量统计详情 |
| PUT | /api/traffic-stats/{id} | 更新流量统计 |
| DELETE | /api/traffic-stats/{id} | 删除流量统计 |
| GET | /api/traffic-stats/summary/global | 全局流量汇总 |
| GET | /api/traffic-stats/summary/by-strategy | 按策略分组统计 |
| GET | /api/traffic-stats/summary/backend/{backend_id} | 单节点统计摘要 |
| POST | /api/route/select | 路由选择（给定 path 返回命中规则+选中节点） |
| GET | /api/route/available | 获取路径下可用后端节点列表 |

## 运行组件

- `internal/health`：并发健康探测与连接租约生命周期。
- `internal/routing`：版本化路由快照和跨请求数据隔离。
- `internal/resilience`：保留错误分类的取消感知重试和状态回写。
- `internal/requestscope`：请求身份快照与异步审计隔离。
- `internal/bootstrap`：运行配置构建、恢复和只读缓存发布。
