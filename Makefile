export GO15VENDOREXPERIMENT=1

.PHONY: test build

PACKAGES = $(shell go list ./ct)

build:
	@go install -v github.com/coreos/terraform-provider-ct

test:
	go vet $(PACKAGES)
	go test -v $(PACKAGES)

revendor:
	@glide update --strip-vendor
	@glide-vc --use-lock-file --no-tests --only-code
