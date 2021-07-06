FROM registry.sorintdev.it/alpine
COPY papagaio-api /app/

ENTRYPOINT ["/app/papagaio-api", "serve"]