FROM alpine:latest
COPY papagaio-api /app/

RUN apk add tzdata
ENTRYPOINT ["/app/papagaio-api", "serve"]