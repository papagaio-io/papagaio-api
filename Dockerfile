FROM registry.sorintdev.it/fedora:minimal
COPY papagaio-api /app/

ENTRYPOINT ["/app/papagaio-api", "serve"]