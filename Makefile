.PHONY: test test-integration test-e2e fmt tidy

test:
	go test ./...

test-integration:
	go test -tags=integration ./...

test-e2e:
	go test -tags=e2e ./...

fmt:
	gofmt -w .

tidy:
	go mod tidy

