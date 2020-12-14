# Clickhouse笔记



## 准备工作



1. 下载数据

```bash
curl https://clickhouse-datasets.s3.yandex.net/hits/tsv/hits_v1.tsv.xz | unxz --threads=`nproc` > hits_v1.tsv
curl https://clickhouse-datasets.s3.yandex.net/visits/tsv/visits_v1.tsv.xz | unxz --threads=`nproc` > visits_v1.tsv
```



2. 启动clickhouse服务

   使用docker-compose启动。

   clone整个工程：[ClickHouse](https://github.com/ClickHouse/ClickHouse/tree/master/docker/server)

   

3. 导入数据

   ```bash
   docker run -i yandex/clickhouse-client 
   ```

   

4. 

   



## 参考资料

- [clickhouse docker offical](https://hub.docker.com/r/yandex/clickhouse-server/)
- [clickhouse-docker-compose](https://github.com/rongfengliang/clickhouse-docker-compose)：docker-compose的示例、导入数据的命令
- [clickhouse教程](https://clickhouse.tech/docs/zh/getting-started/tutorial/)
- [使用Docker安装ClickHouse](https://my.oschina.net/u/4330033/blog/3264678)

