PROJDIR=$(dir $(realpath $(firstword $(MAKEFILE_LIST))))

# change to project dir so we can express all as relative paths
$(shell cd $(PROJDIR))

REPO_PATH=wecode.sorint.it/opensource/papagaio-api

VERSION ?= $(shell scripts/git-version.sh)

$(shell mkdir -p bin )
$(shell mkdir -p tools/bin )

PAPAGAIO_WEBBUNDLE_DEPS = webbundle/bindata.go
PAPAGAIO_WEBBUNDLE_TAGS = webbundle

ifdef WEBBUNDLE

ifndef WEBDISTPATH
$(error WEBDISTPATH must be provided when building the webbundle)
endif

PAPAGAIO_DEPS = $(PAPAGAIO_WEBBUNDLE_DEPS)
PAPAGAIO_TAGS = $(PAPAGAIO_WEBBUNDLE_TAGS)
endif

.PHONY: all
all: build

.PHONY: build
build: papagaio

# don't use existing file names and track go sources, let's do this to the go tool
.PHONY: papagaio
papagaio: $(PAPAGAIO_DEPS)
	GO111MODULE=on go build $(if $(PAPAGAIO_TAGS),-tags "$(PAPAGAIO_TAGS)") -o $(PROJDIR)/bin/papagaio $(REPO_PATH)

.PHONY: go-bindata
go-bindata:
	GOBIN=$(PROJDIR)/tools/bin go install github.com/go-bindata/go-bindata/go-bindata

webbundle/bindata.go: go-bindata $(WEBDISTPATH)
	./tools/bin/go-bindata -o webbundle/bindata.go -tags webbundle -pkg webbundle -prefix "$(WEBDISTPATH)" -nocompress=true "$(WEBDISTPATH)/..."

.PHONY: docker-papagaio
docker-papagaio:
	docker build --target papagaio . -t papagaio