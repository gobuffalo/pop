TAGS ?= "sqlite"
GO_BIN ?= go

install: deps
	$(GO_BIN) install -tags ${TAGS} -v ./soda

deps:
	$(GO_BIN) get github.com/gobuffalo/release
	$(GO_BIN) get -tags ${TAGS} -t ./...
ifeq ($(GO111MODULE),on)
	$(GO_BIN) mod tidy
endif

build:
	$(GO_BIN) build -v .

test:
	$(GO_BIN) test -tags ${TAGS} ./...

ci-test:
	$(GO_BIN) test -tags ${TAGS} -race ./...

lint:
	golangci-lint run

update:
	$(GO_BIN) get -u -tags ${TAGS}
ifeq ($(GO111MODULE),on)
	$(GO_BIN) mod tidy
endif
	make test
	make install
ifeq ($(GO111MODULE),on)
	$(GO_BIN) mod tidy
endif

release-test:
	./test.sh

release:
	release -y -f soda/cmd/version.go
