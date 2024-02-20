#!/bin/sh

# 環境変数を使用して設定ファイルを生成
envsubst '$UPSTREAM_SERVER' < /etc/nginx/conf.d/nginx.conf.template > /etc/nginx/conf.d/nginx.conf

# nginxを起動
exec nginx -g 'daemon off;'
