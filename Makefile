GOFMT ?= gofmt -s
GOFILES := find . -name "*.go" -type f

.PHONY: fmt
fmt:
	$(GOFILES) | xargs $(GOFMT) -w

.PHONY: fmt-check
fmt-check:
	# get all go files and run go fmt on them
	@files=$$($(GOFILES) | xargs $(GOFMT) -l); if [ -n "$$files" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${files}"; \
		exit 1; \
		fi;

build:
	go build .

install:
	go install .

.PHONY: test
test:
	./test.sh

clean:
	go clean .

