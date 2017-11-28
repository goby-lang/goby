GOFMT ?= gofmt -s
GOFILES := $(shell find . -name "*.go" -type f -not -path "./vendor/*")
OPTIONS := -ldflags "-s -w"

.PHONY: fmt
fmt:
	$(GOFMT) -w $(GOFILES)

.PHONY: build
build:
	go build $(OPTIONS) .

.PHONY: install
install:
	go install $(OPTIONS) .

.PHONY: test
test:
	go test $(OPTIONS) ./...

.PHONY: clean
clean:
	go clean .

