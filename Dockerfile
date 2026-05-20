FROM golang:1.26.2 AS builder
WORKDIR /app

ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.org

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=mod -o /bin/gin_blog_upload ./cmd/server

FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates

COPY --from=builder /bin/gin_blog_upload /usr/local/bin/gin_blog_upload

ENV UPLOAD_SERVICE_PORT=9091
ENV UPLOAD_ROOT=/app/uploads
ENV MINIO_ENDPOINT=minio-svc:9000
ENV MINIO_BUCKET=gin-blog-upload
ENV MINIO_USE_SSL=false
ENV MINIO_PUBLIC_BASE=http://192.168.32.141:30080/minio

EXPOSE 9091
ENTRYPOINT ["/usr/local/bin/gin_blog_upload"]

