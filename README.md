# 反向代理服务器，使用Caddy作为服务器实现

## 支持配置的环境变量

- PORT: 对外服务的端口号（默认：80）
- DEFAULT_URL：默认后端服务URI。如：`https://www.google.com:443`
- TRUSTED_PROXIES：信任的下游代理服务`CIDR IP`，多个IP使用逗号分隔，对于信任的IP服务器会采纳其传入的`X-Forward-*`等请求头，
否则就抛弃。如：`192.168.0.1/32`
- METRICS: 是否开启服务指标（默认：true）
- DEBUG：是否开启Debug模式（默认：false）
- LOG_LEVEL：日志等级（默认：INFO）可选DEBUG、WARN、ERROR等
- PATH_$i：代理路径匹配器。如：`PATH_0=/api/**`
- URL_$i：后端服务URI。如：`URL_0=https://www.google.com:443`
