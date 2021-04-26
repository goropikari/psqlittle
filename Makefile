build:
	go build -o bin/repl repl/main.go

test: mockgen
	go test $(shell go list ./...)

mockgen:
	find -name "mock_*.go" | xargs -I {} rm -f {}
	grep -r "go:generate mockgen" | cut -d':' -f1 | xargs -I {} go generate {}
