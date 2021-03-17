#FROM alpine:latest
FROM debian:buster
COPY papagaio-api /app/

RUN apk add tzdata
ENTRYPOINT ["/app/papagaio-api", "serve"]