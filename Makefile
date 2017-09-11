GOFMT ?= gofmt -s
GOFILES := $(shell find . -name "*.go" -type f -not -path "./vendor/*")

.PHONY: fmt
fmt:
	$(GOFMT) -w $(GOFILES)

.PHONY: build
build:
	go build .

.PHONY: install
install:
	go install .

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	go clean .

