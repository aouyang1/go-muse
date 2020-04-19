usage:
	@echo "make all       : Runs all tests, examples, and benchmarks"
	@echo "make test      : Runs test suite"
	@echo "make bench     : Runs benchmarks"
	@echo "make example   : Runs example"
	@echo "make travis-ci : Travis CI specific testing"
	@echo "make cpu-pprof : Runs pprof on the cpu profile from make bench
	@echo "make mem-pprof : Runs pprof on the memory profile from make bench

all: test bench example

test:
	go test -race -coverprofile=coverage.txt -covermode=atomic -cover -run=Test ./...

cover:
	go tool cover -html=coverage.txt
bench:
	go test ./... -run=XX -bench=. -test.benchmem -cpuprofile cpu.prof -memprofile mem.prof

example:
	go test ./... -run=Example

travis-ci:
	go test -v ./... -race -coverprofile=coverage.txt -covermode=atomic

cpu-pprof:
	go tool pprof -pdf cpu.prof  > profile_cpu.pdf && open profile_cpu.pdf

mem-pprof:
	go tool pprof -pdf mem.prof  > profile_mem.pdf && open profile_mem.pdf