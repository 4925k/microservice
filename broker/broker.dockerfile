FROM golang:alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o broker ./cmd/api

RUN chmod +x /app/broker


# tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/broker /app

CMD [ "/app/broker" ]