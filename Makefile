DOCKER_TAG?=latest
# GIT_VERSION := $(shell git describe --tags)
VERSION?=$(GIT_VERSION)
PWD=$(shell pwd)
BUILD_DIR?=$(shell pwd)/build
GOX_OS?=linux darwin windows solaris freebsd netbsd openbsd
GOX_FLAGS="-arch=amd64"

GOPKGS_NOVENDOR=$(shell go list ./... | grep -v "/vendor/")

all: fmt check build

.PHONY: deps
deps:
	glide install

.PHONY: fmt
fmt:
	@ if [ ! $$(gofmt -e -d service |wc -l) -eq 0 ]; then \
		echo "gofmt failed" ; \
		gofmt -e -d service ; exit 1 ; \
	fi

.PHONY: check
check:
	go vet
	golint

.PHONY: build
build:
	@echo "go build"
	@go build $(GOPKGS_NOVENDOR)

.PHONY: test
test:
	go test -v

.PHONY: docker-image
docker-image:
	docker build -t webservice:$(DOCKER_TAG) -f Dockerfile .

.PHONY: crosscompile
crosscompile:
	go get github.com/mitchellh/gox
	mkdir -p ${BUILD_DIR}/bin
	gox -output="${BUILD_DIR}/bin/{{.Dir}}-{{.OS}}-{{.Arch}}" -os="${GOX_OS}" ${GOX_FLAGS}

.PHONY: clean
clean:
	@rm -rf build
