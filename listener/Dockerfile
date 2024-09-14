FROM golang:alpine AS builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o listener .

RUN chmod +x /app/listener


# tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/listener /app

CMD [ "/app/listener" ]