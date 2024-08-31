FROM golang:alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o auth ./cmd/api

RUN chmod +x /app/auth


# tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/auth /app

CMD [ "/app/auth" ]