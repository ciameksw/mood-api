FROM golang:1.25-alpine

RUN mkdir advice
COPY . /advice

WORKDIR /advice

RUN go build ./cmd/advice

CMD ["/advice/advice"]