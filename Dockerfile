FROM golang:latest as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o main .

FROM alpine:latest  
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main -v service"]
