FROM golang:1.25-alpine AS builder

WORKDIR /workspace

COPY go.work go.work.sum ./
COPY pkg/go.mod pkg/go.sum ./pkg/
COPY mood/go.mod mood/go.sum ./mood/
COPY auth/go.mod auth/go.sum ./auth/
COPY advice/go.mod advice/go.sum ./advice/
COPY gateway/go.mod ./gateway/
COPY quote/go.mod ./quote/

WORKDIR /workspace/quote
RUN go mod download

WORKDIR /workspace
COPY pkg/ ./pkg/
COPY quote/ ./quote/

WORKDIR /workspace/quote
RUN go build -o /quote ./cmd/quote

FROM alpine:latest

COPY --from=builder /quote /quote

CMD ["/quote"]