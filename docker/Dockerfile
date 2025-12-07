FROM rust:alpine AS rust-builder

RUN apk add --no-cache musl-dev openssl-dev

WORKDIR /app

COPY Cargo.toml ./
COPY src/ ./src/

RUN cargo build --release

FROM alpine:latest

RUN apk --no-cache add ca-certificates curl libgcc

WORKDIR /root/

COPY --from=rust-builder /app/target/release/nusa-node .

RUN chmod +x nusa-node

HEALTHCHECK --interval=10s --timeout=3s CMD curl -f http://localhost:8545/health || exit 1

EXPOSE 8545 8546

CMD ["./nusa-node"]
