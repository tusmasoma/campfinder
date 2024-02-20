FROM --platform=linux/amd64 nginx:1.18.0-alpine

# デフォルトのconfを消して、別の設定を追加
WORKDIR /app

RUN rm /etc/nginx/conf.d/default.conf
COPY docker/nginx/nginx.conf.template /etc/nginx/conf.d/nginx.conf.template
COPY docker/nginx/entrypoint.sh /nginx-entrypoint.sh
RUN chmod +x /nginx-entrypoint.sh

ENTRYPOINT ["/nginx-entrypoint.sh"]
