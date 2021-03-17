FROM alpine:latest
COPY papagaio-api /app/

RUN apk update && apk add bash
RUN apk add tzdata
ENTRYPOINT ["/app/papagaio-api", "serve"]