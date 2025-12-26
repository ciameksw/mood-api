FROM golang:1.25-alpine

RUN mkdir auth
COPY . /auth

WORKDIR /auth

RUN go build ./cmd/auth

CMD ["/auth/auth"]