FROM golang:1.25-alpine AS builder

WORKDIR /workspace

COPY go.work go.work.sum ./
COPY pkg/go.mod pkg/go.sum ./pkg/
COPY mood/go.mod mood/go.sum ./mood/
COPY auth/go.mod ./auth/
COPY advice/go.mod ./advice/
COPY gateway/go.mod ./gateway/
COPY quote/go.mod ./quote/

WORKDIR /workspace/mood
RUN go mod download

WORKDIR /workspace
COPY pkg/ ./pkg/
COPY mood/ ./mood/

WORKDIR /workspace/mood
RUN go build -o /mood ./cmd/mood

FROM alpine:latest

COPY --from=builder /mood /mood

CMD ["/mood"]