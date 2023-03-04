#!/bin/bash
set -e
# 生成配置文件 /etc/caddy/Caddyfile
./gen_caddyfile
# 启动服务器
exec caddy run --config /etc/caddy/Caddyfile --adapter caddyfile
