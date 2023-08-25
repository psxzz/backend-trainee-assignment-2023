PROJECT_PATH=$(CURDIR)
BINDIR=$(CURDIR)/bin
PACKAGE=cmd/app
OS=linux

GOLANGCI_BIN:=$(BINDIR)/golangci-lint
GOLANGCI_REPO=https://github.com/golangci/golangci-lint
GOLANGCI_LATEST_VERSION:= $(shell git ls-remote --tags --refs --sort='v:refname' $(GOLANGCI_REPO)|tail -1|egrep -o "v[0-9]+.*")

all: deps build test lint

bindir:
	mkdir -p ${BINDIR}

build: bindir
	GOOS=${OS} go build -o ${BINDIR}/app ${PACKAGE}/*.go

test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

install-lint: bindir
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s $(GOLANGCI_LATEST_VERSION)

lint: install-lint
	${GOLANGCI_BIN} run --config=${PROJECT_PATH}/.golangci.yaml -v ${PROJECT_PATH}/...

deps:
	go mod tidy
	go mod vendor
	go mod verify

clean:
	go clean
	rm -rf ${BINDIR}
