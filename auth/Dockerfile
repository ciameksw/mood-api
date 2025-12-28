FROM golang:1.25-alpine AS builder

WORKDIR /workspace

COPY go.work go.work.sum ./
COPY pkg/go.mod pkg/go.sum ./pkg/
COPY auth/go.mod auth/go.sum ./auth/
COPY mood/go.mod ./mood/
COPY advice/go.mod ./advice/
COPY gateway/go.mod ./gateway/
COPY quote/go.mod ./quote/

WORKDIR /workspace/auth
RUN go mod download

WORKDIR /workspace
COPY pkg/ ./pkg/
COPY auth/ ./auth/

WORKDIR /workspace/auth
RUN go build -o /auth ./cmd/auth

FROM alpine:latest

COPY --from=builder /auth /auth

CMD ["/auth"]