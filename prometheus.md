# Prometheus / Grafana





## 错误统计

```
sum by (url) (rate(mmrpc_error_count{svtype="go_gossip_darwin"}[5m]) / rate(go_gossip_darwin_mmrpc_count{svtype="go_gossip_darwin"}[5m]))
```

解释：

1. `sum by (url)` - 根据url累加。
2. `5m` - 5分钟一个维度统计。
3. `mmrpc_error_count`、`go_gossip_darwin_mmrpc_count` - label名字。



## 错误数量的变化值

```
sum by (url) (delta(go_gossip_darwin_mmrpc_error_count{svtype="go_gossip_darwin"}[5m]))
```

解释：

1. `5m` - 5分钟的差值。如果是前一天的话，那么使用`1d`。

2. `delta` - 差值函数。
3. `sum by (url)` - 按照url进行counter聚合。