GOFMT ?= gofmt -s
GOFILES := $(shell find . -name "*.go" -type f -not -path "./vendor/*")
RELEASE_OPTIONS := -ldflags "-s -w -X github.com/goby-lang/goby/vm.DefaultLibPath=${GOBY_LIBPATH}" -tags release
TEST_OPTIONS := -ldflags "-s -w"
ENV := CGO_ENABLED=0

.PHONY: fmt
fmt:
	$(GOFMT) -w $(GOFILES)

.PHONY: build
build:
	$(ENV) go build $(RELEASE_OPTIONS) .

.PHONY: install
install:
	$(ENV) go install $(RELEASE_OPTIONS) .

.PHONY: test
test:
	$(ENV) go test $(TEST_OPTIONS) ./...

.PHONY: clean
clean:
	go clean .
