# Goのビルド環境
FROM golang:1.21-alpine as builder

# 作業ディレクトリを設定
WORKDIR /app

# 依存関係ファイルをコピー
COPY go.mod ./
COPY go.sum ./

# 依存関係をインストール
RUN go mod download

# プロジェクトの全ファイルをコピー
COPY . .

# プロジェクトルートのmain.goをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -o /discord-bot .

# 実行環境
FROM alpine:latest

# 必要なファイルやディレクトリを新しいイメージにコピー
COPY --from=builder /discord-bot /discord-bot
COPY --from=builder /app/public /public
COPY --from=builder /app/.env /.env

# 実行
CMD ["/discord-bot"]
