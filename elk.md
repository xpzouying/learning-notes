# ELK

## 日志系统

### 参考资料

- [Filebeat vs. Logstash: The Evolution of a Log Shipper](https://logz.io/blog/filebeat-vs-logstash/)
  - Filebeat - 作为日志收集器
  - Logstash - 作为聚合起，从不同的数据源获取数据，然后导入到Elasticsearch中。

### 说明

1. Filebeat典型的配置文件：

   ```
   filebeat.inputs:
   - input_type: log
     paths:
       - /var/log/httpd/access.log
   
   document_type: apache-access
   fields_under_root: true
   
   output.logstash:
     hosts: ["127.0.0.1:5044"]
   ```

   

2. Logstash典型的配置文件：

   ```
   input {
     beats {
       port => 5044
     }
    }
   
   filter {
     grok {
       match => { "message" => "%{IPORHOST:clientip} %{USER:ident} %{USER:auth} \[%{HTTPDATE:time}\] "(?:%{WORD:verb} %{NOTSPACE:request}(?: HTTP/%{NUMBER:httpversion})?|%{DATA:rawrequest})" %{NUMBER:response} (?:%{NUMBER:bytes}|-)" }
     }
     date {
       match => [ "timestamp" , "dd/MMM/yyyy:HH:mm:ss Z" ]
     }
    }
   
   output {
     elasticsearch { hosts => ["localhost:9200"] }
    }
   ```

   
