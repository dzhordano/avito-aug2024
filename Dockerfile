FROM golang:1.22.5-alpine3.20 AS builder

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build -o ./.bin ./cmd/app/main.go

FROM alpine:3.20

COPY --from=builder /app .

CMD ["./.bin"]