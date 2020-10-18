# Linux常用工具


## 测试局域网内的网络

**1、环境准备**

- 局域网内，2台设备
- 安装测试软件：`iperf3`

我使用的设备为：
1. Macbook Pro：安装`iperf3`。
2. iPad Pro：安装`he.net - Network Tools`。

其他平台的支持见：

- [iperf download page](https://iperf.fr/iperf-download.php)


**2、测试**

2.1、在MBP上，建立Server

`iperf3`解压后，运行Server。

```bash
./iperf3 -s
```

2.2、从iPad Pro上，请求MBP Server的地址，进行测试


- 打开`HE.NET`工具，菜单点击：`Iperf`
- 选择`Iperf`，地址栏中填写：`{server_ip}`
- 配置参数：
  - Interval：表示间隔多久输出日志信息
  - Bytes：总共请求多少数据。填入：`10240M`，即10G数据
- 回车


2.3、从其他PC上面使用命令测试：

运行命令：

```bash
iperf3 -c {server_ip} -i1 -t10
```

- c：后面接iperf server地址
- i：打印间隔
- t：运行时长



**3、结果**


```
[  5] 387.00-388.00 sec  33.4 MBytes   280 Mbits/sec
[  5] 388.00-389.00 sec  33.0 MBytes   277 Mbits/sec
[  5] 389.00-390.00 sec  33.3 MBytes   280 Mbits/sec
[  5] 390.00-391.00 sec  34.8 MBytes   292 Mbits/sec
[  5] 391.00-392.00 sec  36.1 MBytes   303 Mbits/sec
[  5] 392.00-392.65 sec  22.7 MBytes   293 Mbits/sec
- - - - - - - - - - - - - - - - - - - - - - - - -
[ ID] Interval           Transfer     Bandwidth
[  5]   0.00-392.65 sec  0.00 Bytes  0.00 bits/sec                  sender
[  5]   0.00-392.65 sec  9.76 GBytes   214 Mbits/sec                  receiver
-----------------------------------------------------------
Server listening on 5201
-----------------------------------------------------------
```
