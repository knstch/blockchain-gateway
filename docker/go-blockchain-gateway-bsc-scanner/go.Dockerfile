FROM golang:1.24 AS base

FROM base AS builder

WORKDIR /build
COPY . ./
RUN go build ./cmd/bsc-scanner

FROM base AS final

WORKDIR /app
COPY --from=builder /build/bsc-scanner /build/.env ./
COPY --from=builder /build/bsc-scanner ./

CMD ["/app/bsc-scanner"]