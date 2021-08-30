ARG PAPAGAIOWEB_IMAGE="papagaio-web"

FROM $PAPAGAIOWEB_IMAGE AS papagaio-web

#######
####### Build the backend
#######

# base build image
FROM registry.sorintdev.it/golang:1.16-buster AS build_base

WORKDIR /papagaio-api

# use go modules
ENV GO111MODULE=on

# only copy go.mod and go.sum
COPY go.mod .
COPY go.sum .

RUN go mod download

# builds the papagaio binaries
FROM build_base AS server_builder

# copy all the sources
COPY . .

# copy the papagaio-web dist
COPY --from=papagaio-web /usr/share/nginx/html/ /papagaio-web/dist/

RUN make WEBBUNDLE=1 WEBDISTPATH=/papagaio-web/dist/

#######
####### Build the final image
#######
FROM registry.sorintdev.it/fedora-minimal AS papagaio

WORKDIR /

COPY --from=server_builder /papagaio-api/bin/papagaio-api /bin/

ENTRYPOINT ["papagaio-api", "serve"]