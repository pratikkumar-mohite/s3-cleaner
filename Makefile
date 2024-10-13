.PHONY: mod build run test

mod:
	@go mod download
	@go mod verify

build: bench
	@CGO_ENABLED=0 go build -v -o s3-cleanup ./cmd/s3-cleanup/

run:
	@CGO_ENABLED=0 go run ./cmd/s3-cleanup/

test:
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run -v --timeout=10m
	@go test -v -timeout 20m -covermode=atomic -coverprofile=coverage.txt ./...

bench:
	@go test -v -bench=. -count=6 -cpuprofile=cpu.out -memprofile=mem.out ./cmd/s3-cleanup/ -run=^Benchmark$ > benchmark.txt

bench-heavy:
	@go test -v -bench=. -count=10 -cpu=1,2,4,8 -cpuprofile=cpu.out -memprofile=mem.out ./cmd/s3-cleanup/ -run=^Benchmark$ > benchmark.txt
	@benchstat benchmark.txt
	@echo top | go tool pprof mem.out
	@echo top | go tool pprof cpu.out

result:
	@cp coverage.txt benchmark.txt mem.out cpu.out coverage
	@go tool cover -func coverage.txt
	@benchstat benchmark.txt
	@echo top | go tool pprof mem.out
	@echo top | go tool pprof cpu.out
