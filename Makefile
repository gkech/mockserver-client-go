lint: ## Perform linting
	docker run --env=GOFLAGS=-mod=vendor --rm -v $(CURDIR):/app -w /app golangci/golangci-lint:v1.51.2 golangci-lint run

test: ## Run unit tests
	go test ./... -mod=vendor -race -count=1

fmt: ## Format the source code
	go fmt ./...

mocks: ## Generate the mocks used from the various tests of this service
	# vendor mock
	mockgen -source mockserver.go -destination test/mock/http_client.go -package mock -mock_names Client=MockHttpClient

.PHONY: lint test fmt mocks
