.PHONY: test test-unit test-integration test-integration-docker test-e2e test-e2e-docker test-docker fmt tidy

TESTFLAGS ?=

ifneq ($(JUNO_TEST_LOG),)
TESTFLAGS += -v
endif

test-unit:
	go test $(TESTFLAGS) ./...

test: test-unit test-integration test-e2e

test-integration:
	go test $(TESTFLAGS) -tags=integration ./...

test-integration-docker:
	go test $(TESTFLAGS) -tags=integration,docker ./...

test-e2e:
	go test $(TESTFLAGS) -tags=e2e ./...

test-e2e-docker:
	go test $(TESTFLAGS) -tags=e2e,docker ./...

test-docker: test-integration-docker test-e2e-docker

fmt:
	gofmt -w .

tidy:
	go mod tidy
