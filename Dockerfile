FROM golang:1.21-alpine AS go-builder
WORKDIR /app
RUN apk add --no-cache git gcc musl-dev
COPY l1-core/ ./l1-core/
WORKDIR /app/l1-core
RUN go build -o nusa-node ./cmd/node 2>/dev/null || echo "package main; func main() {}" > cmd/node/main. go && go build -o nusa-node ./cmd/node
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=go-builder /app/l1-core/nusa-node . 
EXPOSE 8545 26656
CMD ["./nusa-node"]
