FROM golang:latest as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X github.com/ljcbaby/hdu-wiki-qa/cmd.Version={{.Version}}"

FROM alpine:latest  
WORKDIR /app
COPY --from=builder /app/hdu-wiki-qa .
RUN ls -la /app  # 添加这一行检查文件是否存在
EXPOSE 8080
CMD ["./hdu-wiki-qa", "-v", "service"]
