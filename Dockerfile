FROM registry.sorintdev.it/fedora:33
COPY papagaio-api /app/

ENTRYPOINT ["/app/papagaio-api", "serve"]