GOFMT ?= gofmt -s
GOFILES := $(shell find . -name "*.go" -type f -not -path "./vendor/*")
RELEASE_OPTIONS := -ldflags "-s -w -X github.com/goby-lang/goby/vm.DefaultLibPath=${GOBY_LIBPATH}" -tags release
TEST_OPTIONS := -ldflags "-s -w"

INSTRUCTION := compiler/bytecode/instruction.go

.PHONY: fmt
fmt:
	$(GOFMT) -w $(GOFILES)

.PHONY: build generate
build:
	go build $(RELEASE_OPTIONS) .

.PHONY: install generate
install:
	go install $(RELEASE_OPTIONS) .

.PHONY: test generate
test:
	go test $(TEST_OPTIONS) ./...

.PHONY: clean
clean:
	go clean .

.PHONY: generate
generate:
	go generate $(INSTRUCTION)
