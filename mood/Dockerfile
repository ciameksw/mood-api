FROM golang:1.25-alpine

RUN mkdir mood
COPY . /mood

WORKDIR /mood

RUN go build ./cmd/mood

CMD ["/mood/mood"]