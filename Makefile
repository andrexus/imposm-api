PROJECT_NAME?=imposm-api
PROJECT_NAMESPACE?=andrexus
PROJECT?=github.com/${PROJECT_NAMESPACE}/${PROJECT_NAME}
BUILD_DIR=build
BUILDTAGS=
OSARCH=darwin/amd64 linux/amd64 linux/arm windows/amd64
DIST_USER=andrexus
GIT_VERSION=$(shell git describe --abbrev=0 --tags 2> /dev/null || git rev-parse HEAD)
GOLDFLAGS=-s -w -X $(PROJECT)/cmd.Version=$(GIT_VERSION)
GOTOOLS = \
	github.com/golang/dep/cmd/dep \
    github.com/golang/lint/golint \
	github.com/alvaroloes/enumer \
    github.com/mitchellh/gox \
    github.com/tcnksm/ghr

all: clean vendor generate test build

build: generate
	@echo "==> Building ${PROJECT_NAME}..."
	CGO_ENABLED=0 gox -ldflags "${GOLDFLAGS}" -osarch "${OSARCH}" -output "${BUILD_DIR}/{{.OS}}_{{.Arch}}_{{.Dir}}"

dist:
	ghr -u ${DIST_USER} --token ${GITHUB_TOKEN} --replace ${GIT_VERSION} build/

vendor: tools
	@echo "==> Vendoring dependencies..."
	dep ensure

fmt:
	@echo "+ $@"
	@go list -f '{{if len .TestGoFiles}}"gofmt -s -l {{.Dir}}"{{end}}' $(shell go list ${PROJECT}/... | grep -v vendor) | xargs -L 1 sh -c

lint:
	@echo "+ $@"
	@go list -f '{{if len .TestGoFiles}}"golint {{.Dir}}/..."{{end}}' $(shell go list ${PROJECT}/... | grep -v vendor) | xargs -L 1 sh -c

vet:
	@echo "+ $@"
	@go vet $(shell go list ${PROJECT}/... | grep -v vendor)

test: fmt lint vet
	@echo "+ $@"
	@go test -v -race -tags "$(BUILDTAGS) cgo" $(shell go list ${PROJECT}/... | grep -v vendor)

cover:
	@echo "+ $@"
	@go list -f '{{if len .TestGoFiles}}"go test -coverprofile={{.Dir}}/.coverprofile {{.ImportPath}}"{{end}}' $(shell go list ${PROJECT}/... | grep -v vendor) | xargs -L 1 sh -c

tools:
	@echo "==> Installing additional tools..."
	go get -u $(GOTOOLS)

generate:
	go generate

clean:
	rm -rf ${BUILD_DIR}/*

.PHONY: clean dist