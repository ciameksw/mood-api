FROM golang:1.25-alpine AS builder

WORKDIR /workspace

COPY go.work go.work.sum ./
COPY pkg/go.mod pkg/go.sum ./pkg/
COPY advice/go.mod advice/go.sum ./advice/
COPY mood/go.mod ./mood/
COPY auth/go.mod ./auth/
COPY gateway/go.mod ./gateway/
COPY quote/go.mod ./quote/

WORKDIR /workspace/advice
RUN go mod download

WORKDIR /workspace
COPY pkg/ ./pkg/
COPY gateway/ ./gateway/

WORKDIR /workspace/gateway
RUN go build -o /gateway ./cmd/gateway

FROM alpine:latest

COPY --from=builder /gateway /gateway

CMD ["/gateway"]