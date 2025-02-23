deps:
	go install github.com/daixiang0/gci@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install mvdan.cc/gofumpt@latest
	go install github.com/bombsimon/wsl/v4/cmd...@latest
	go install go.uber.org/mock/mockgen@latest
.PHONY: deps

run:
	go mod tidy
	go run cmd/main.go
.PHONY: run

ready: update-deps test lint
.PHONY: ready

update-deps:
	go get -u ./...
	go mod tidy
.PHONY: update-deps

test:
	go test ./...
.PHONY: test

lint:
	go mod tidy
	go vet ./...
	gci write --skip-generated -s standard -s default -s "prefix(surfe-actions)" .
	gofumpt -l -w .
	wsl -fix ./... 2>/dev/null || true # Ignore errors because wsl fixes them
	golangci-lint run $(p)
	go fmt ./...
.PHONY: lint

### Mocks
mock:
	go generate ./...
.PHONY: mock