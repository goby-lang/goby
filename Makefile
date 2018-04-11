GOFMT ?= gofmt -s
GOFILES := $(shell find . -name "*.go" -type f -not -path "./vendor/*")
RELEASE_OPTIONS := -ldflags "-s -w -X github.com/goby-lang/goby/vm.DefaultLibPath=${GOBY_LIBPATH}"
TEST_OPTIONS := -ldflags "-s -w"

.PHONY: fmt
fmt:
	$(GOFMT) -w $(GOFILES)

.PHONY: build
build:
	go build $(RELEASE_OPTIONS) .

.PHONY: install
install:
	go install $(RELEASE_OPTIONS) .

.PHONY: test
test:
	go test $(TEST_OPTIONS) ./...

.PHONY: clean
clean:
	go clean .

