FROM golang:alpine AS builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o fe ./cmd/web

RUN chmod +x /app/fe


# tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/fe /app

CMD [ "/app/fe" ]