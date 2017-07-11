.PHONY: test build vendor

PACKAGES = $(shell go list ./ct)

build:
	@go install -v github.com/coreos/terraform-provider-ct

test:
	go vet $(PACKAGES)
	go test -v $(PACKAGES)

vendor:
	@glide update --strip-vendor
	@glide-vc --use-lock-file --no-tests --only-code
