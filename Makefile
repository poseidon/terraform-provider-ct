.PHONY: test build

PACKAGES = $(shell go list ./fuze)

build:
	@go install -v .

test:
	go vet $(PACKAGES)
	go test -v $(PACKAGES)
