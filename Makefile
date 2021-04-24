test: mockgen
	go test $(shell go list ./...)

mockgen:
	grep -r "go:generate mockgen" | cut -d':' -f1 | xargs -I {} go generate {}
