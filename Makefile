GOFMT ?= gofmt -s
GOFILES := find . -name "*.go" -type f

.PHONY: fmt
fmt:
	$(GOFILES) | xargs $(GOFMT) -w

.PHONY: build
build:
	go build .

.PHONY: install
install:
	go install .

.PHONY: test
test:
	./test.sh

.PHONY: clean
clean:
	go clean .

