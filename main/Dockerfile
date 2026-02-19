FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download && \
  CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/

# Run stage
FROM alpine:3.21
WORKDIR /app

# タイムゾーンを日本時間に設定
RUN apk add --no-cache tzdata curl bash && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    echo "Asia/Tokyo" > /etc/timezone && \
    apk del tzdata
ENV TZ=Asia/Tokyo

COPY --from=builder /app/main .

# ユーザーを変更
RUN addgroup -g 1000 appuser && \
    adduser -D -s /bin/sh -u 1000 -G appuser appuser
USER appuser

EXPOSE 8080

CMD ["./main"]
