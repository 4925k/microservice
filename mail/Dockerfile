FROM golang:alpine AS builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o mail ./cmd/api

RUN chmod +x /app/mail


# tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/mail /app

CMD [ "/app/mail" ]