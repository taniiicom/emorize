# usage:
# - `docker build . -t emorize`
# - `docker run emorize`

# Goのビルド環境
FROM golang:1.21-alpine as builder

# label
LABEL maintainer="taniiicom <mail@taniii.com>"

# 作業ディレクトリを設定
WORKDIR /app

# 依存関係ファイルをコピー
COPY go.mod ./
COPY go.sum ./

# 依存関係をインストール
RUN go mod download

# プロジェクトの全ファイルをコピー
COPY . .

# [prod] product 環境向け
# .env ファイルが存在しない場合は空の .env ファイルを作成
RUN [ ! -f .env ] && touch .env || exit 0

# プロジェクトルートのmain.goをビルド
# -v = verbose = ビルド中のログの冗長詳細出力
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /discord-bot .

# 実行用の新しいステージ
# FROM alpine:latest
FROM gcr.io/distroless/base-debian10

# 必要なファイルやディレクトリを新しいイメージにコピー
COPY --from=builder /discord-bot /discord-bot
COPY --from=builder /app/public /public
COPY --from=builder /app/.env /.env

# アプリケーションがリッスンするポート番号を指定
EXPOSE 8080

# アプリケーションの起動
CMD ["/discord-bot"]