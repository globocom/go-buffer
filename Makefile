test:
	go run github.com/onsi/ginkgo/ginkgo -keepGoing -progress -timeout 1m -race --randomizeAllSpecs --randomizeSuites

bench:
	go test -bench=. -run=Benchmark
