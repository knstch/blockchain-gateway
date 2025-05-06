FROM golang:1.24 AS base

FROM base AS builder

WORKDIR /build
COPY . ./
RUN go build ./cmd/blockchain-gateway

FROM base AS final

ARG PORT

WORKDIR /app
COPY --from=builder /build/blockchain-gateway /build/.env ./
COPY --from=builder /build/blockchain-gateway ./

EXPOSE ${PUBLIC_HTTP_ADDR}
EXPOSE ${PRIVATE_GRPC_ADDR}
CMD ["/app/blockchain-gateway"]