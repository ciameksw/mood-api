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
COPY advice/ ./advice/

WORKDIR /workspace/advice
RUN go build -o /advice ./cmd/advice

FROM alpine:latest

COPY --from=builder /advice /advice

CMD ["/advice"]