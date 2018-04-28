GOFMT ?= gofmt -s
GOFILES := $(shell find . -name "*.go" -type f -not -path "./vendor/*")
RELEASE_OPTIONS := -ldflags "-s -w -X github.com/goby-lang/goby/vm.DefaultLibPath=${GOBY_LIBPATH}" -tags release
TEST_OPTIONS := -ldflags "-s -w"
BENCHMARK_OPTIONS := -run '^$$' -bench '.' -benchmem -benchtime 2s

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

.PHONY: update_benchmarks
update_benchmarks:
	go test $(BENCHMARK_OPTIONS) ./... > current_benchmarks

.PHONY: compare_benchmarks
compare_benchmarks: 
	go test $(BENCHMARK_OPTIONS) ./... > .tmp_benchmarks 
	benchcmp current_benchmarks .tmp_benchmarks



