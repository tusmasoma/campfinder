# ベースイメージ
FROM golang:1.21.3

# MySQLクライアントをインストール
RUN apt-get update && apt-get install -y default-mysql-client

# 作業ディレクトリの設定
WORKDIR /app

# 依存関係をコピー
COPY docker/back/go.mod ./docker/back/
COPY docker/back/go.sum ./docker/back/

# 依存関係のインストール
WORKDIR /app/docker/back
RUN go mod download

WORKDIR /app

# ソースコードをコピー
COPY . .

# エントリポイントスクリプトをコピー
COPY docker/back/entrypoint.sh /usr/local/bin/

# エントリポイントスクリプトを実行可能にする
RUN chmod +x /usr/local/bin/entrypoint.sh

# アプリケーションのビルド
# GOARCH=amd64 は64ビットのx86アーキテクチャ用にバイナリをビルドする
WORKDIR /app/docker/back
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/main.go

# 作業ディレクトリの設定
WORKDIR /root/

# ビルドされたバイナリを最終的な作業ディレクトリに移動
RUN cp /app/docker/back/server .
RUN cp -r /app/docker/back/.certificate/ .certificate/

# コンテナが起動するときに実行されるコマンド (バイナリにしたgolangのファイルを実行)
ENTRYPOINT ["entrypoint.sh"]
CMD ["./server"]