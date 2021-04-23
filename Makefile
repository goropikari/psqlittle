test:
	go test $(shell go list ./...)

mockgen:
	mockgen --source translator/translator.go -destination translator/mock/translator_mock.go --package mock
