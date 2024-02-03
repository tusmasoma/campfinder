# ベースイメージ
FROM golang:1.21.3

# MySQLクライアントをインストール
RUN apt-get update && apt-get install -y default-mysql-client

# 作業ディレクトリの設定
WORKDIR /app

# 依存関係をコピー
COPY go.mod .
COPY go.sum .

# 依存関係のインストール
RUN go mod download

# ソースコードをコピー
COPY . .

# エントリポイントスクリプトをコピー
COPY entrypoint.sh /usr/local/bin/

# エントリポイントスクリプトを実行可能にする
RUN chmod +x /usr/local/bin/entrypoint.sh

# アプリケーションのビルド
# GOARCH=amd64 は64ビットのx86アーキテクチャ用にバイナリをビルドする
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/main.go

# 作業ディレクトリの設定
WORKDIR /root/

# ビルドされたバイナリを最終的な作業ディレクトリに移動
RUN cp /app/server .
RUN cp -r /app/.certificate/ .certificate/

# コンテナが起動するときに実行されるコマンド (バイナリにしたgolangのファイルを実行)
ENTRYPOINT ["entrypoint.sh"]
CMD ["./server"]
