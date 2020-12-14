# Linux常用工具



## Linux Cmd



### awk

**参考资料：**

- [How To Use the AWK language to Manipulate Text in Linux](https://www.digitalocean.com/community/tutorials/how-to-use-the-awk-language-to-manipulate-text-in-linux)

**笔记：**

- 基本命令

  ```bash
  awk '/search_pattern/ { action_to_take_on_matches; another_action; }' file_to_parse
  ```

  

- 



## 结合UPS实现自动关机

**1、大致原理：**

1. 主机插入UPS电源；开机启动该脚本；
1. 脚本每隔一段时间间隔后，ping局域网中的一台设备，且该设备未接入UPS电源，比如路由器（192.168.1.1）。如果ping不通，则表示可能由于市电导致关机，则准备使用脚本进行关机。



```bash
#!/bin/bash

target_ip=192.168.1.1
failure_count=0
shutdown_failure_count_threshold=15

while :
do
  ping -c 1 $target_ip &> /dev/null
  if [ $? -eq 0 ]; then
    ((failure_count=0))
  else
    ((failure_count++))
  fi
  sleep 10s
  if [ $failure_count -eq $shutdown_failure_count_threshold ]; then
    /sbin/shutdown -hP now
    break
  fi
done

exit 0
```

**添加为开机启动**

不使用`sudo`来运行该命令，直接增加`root`用户的crontab。例如：直接运行：

```bash
# 增加root的crontab
sudo crontab -e
```

crontab的命令示例如下，

```bash
@reboot  /path/to/job
@reboot  /path/to/shell.script
@reboot  /path/to/command arg1 arg2
```

增加我们的命令：

```bash
# 注意脚本中的 /sbin/shutdown需要使用绝对路径
@reboot /usr/local/bin/ups-safe-shutdown.sh &
```

**参考来源：**

- [Proxmox VE下配合UPS使用的断电关机脚本](https://juejin.im/post/6874098313839575047)
- [running command  at startup on crontab](https://askubuntu.com/questions/735935/running-command-at-startup-on-crontab)
- [Linux execute cron job after system reboot](https://www.cyberciti.biz/faq/linux-execute-cron-job-after-system-reboot/)


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


## MySQL主从启动

```bash
# git clone
git clone https://github.com/vbabak/docker-mysql-master-slave

# run
./build.sh
```


mydb数据库会进行同步。


## Mac OS 自动根据WI-FI名字改变网络位置

- [Mac OS 自动根据WI-FI名字改变网络位置](https://razeencheng.com/post/auto-change-network-location-base-on-name-of-wifi.html)



