# Prometheus / Grafana





## 错误统计

```
sum by (url) (rate(mmrpc_error_count{svtype="go_gossip_darwin"}[5m]) / rate(go_gossip_darwin_mmrpc_count{svtype="go_gossip_darwin"}[5m]))
```

解释：

1. `sum by (url)` - 根据url累加。
2. `5m` - 5分钟一个维度统计。
3. `mmrpc_error_count`、`go_gossip_darwin_mmrpc_count` - label名字。