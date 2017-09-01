VERSION        ?= $(shell git describe --tags --always )
TARGET         ?= $(shell basename `git rev-parse --show-toplevel`)
TEST           ?= $(shell go list ./... | grep -v /vendor/)
REPOSITORY     := mattdeboer/${TARGET}
DOCKER_IMAGE   ?= ${REPOSITORY}:${VERSION}
BRANCH         ?= $(shell git rev-parse --abbrev-ref HEAD)
REVISION       ?= $(shell git rev-parse HEAD)
LD_FLAGS       ?= -s -X github.com/matt-deboer/${TARGET}/pkg/version.Name=$(TARGET) \
	-X github.com/matt-deboer/${TARGET}/pkg/version.Revision=$(REVISION) \
	-X github.com/matt-deboer/${TARGET}/pkg/version.Branch=$(BRANCH) \
	-X github.com/matt-deboer/${TARGET}/pkg/version.Version=$(VERSION)

default: test build

.PHONY: dep
dep:
	@ which dep >/dev/null || go get github.com/golang/dep/cmd/dep
	dep ensure

test:
	go test -v -cover -run=$(RUN) $(TEST)

build: clean
	@go build -v -o bin/$(TARGET) -ldflags "$(LD_FLAGS)+local_changes" ./pkg/cmd

release: clean dep
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build \
		-a -tags netgo \
		-a -installsuffix cgo \
    -ldflags "$(LD_FLAGS)" \
		-o bin/$(TARGET) ./pkg/cmd

ca-certificates.crt:
	@-docker rm -f ${TARGET}_cacerts
	@docker run --name ${TARGET}_cacerts debian:latest bash -c 'apt-get update && apt-get install -y ca-certificates'
	@docker cp ${TARGET}_cacerts:/etc/ssl/certs/ca-certificates.crt .
	@docker rm -f ${TARGET}_cacerts

docker: ca-certificates.crt release
	@echo "Building ${DOCKER_IMAGE}..."
	@docker build -t ${DOCKER_IMAGE} -f Dockerfile .

clean:
	@rm -rf bin/