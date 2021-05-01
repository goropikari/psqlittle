build:
	go build -o bin/repl cmd/repl/main.go
	go build -o bin/server cmd/server/main.go

test: mockgen
	go test -v $(shell go list ./...)

mockgen:
	find -name "mock_*.go" | xargs -I {} rm -f {}
	grep -r "go:generate mockgen" | cut -d':' -f1 | xargs -I {} go generate {}
