lint: ## Perform linting
	golangci-lint run --enable bodyclose,gofmt,revive,goimports,gosec,gocyclo --exclude-use-default=false --modules-download-mode=vendor --timeout=3m --build-tags=integration,component

test: ## Run unit tests
	go test ./... -mod=vendor -race -count=1

fmt: ## Format the source code
	go fmt ./...

mocks: ## Generate the mocks used from the various tests of this service
	# vendor mock
	mockgen -source mockserver.go -destination test/mock/http_client.go -package mock -mock_names Client=MockHttpClient