
export GO111MODULE=on

.PHONY: bench
bench:
	go test -run=^$$ -bench=. -benchmem -count 5 -cpuprofile=cpu.prof -memprofile=mem.prof

.PHONY: clean
clean:
	go clean
	rm -f cpu.prof mem.prof benchmarks.test
