export CGO_ENABLED:=0

VERSION=$(shell git describe --tags --match=v* --always)
SEMVER=$(shell git describe --tags --match=v* --always | cut -c 2-)

.PHONY: all
all: build test vet fmt

.PHONY: build
build:
	@go build -o $@ github.com/poseidon/terraform-provider-ct

.PHONY: test
test:
	@go test ./... -cover

.PHONY: vet
vet:
	@go vet -all ./...

.PHONY: fmt
fmt:
	@test -z $$(go fmt ./...)

.PHONY: lint
lint:
	@golangci-lint run ./...

.PHONY: clean
clean:
	@rm -rf bin
	@rm -rf _output

.PHONY: release
release: \
	clean \
	_output/plugin-linux-amd64.zip \
	_output/plugin-linux-arm64.zip \
	_output/plugin-darwin-amd64.zip \
	_output/plugin-darwin-arm64.zip \
	_output/plugin-windows-amd64.zip

_output/plugin-%.zip: NAME=terraform-provider-ct_$(SEMVER)_$(subst -,_,$*)
_output/plugin-%.zip: DEST=_output/$(NAME)
_output/plugin-%.zip: LOCAL=$(HOME)/.terraform.d/plugins/terraform.localhost/poseidon/ct/$(SEMVER)
_output/plugin-%.zip: _output/%/terraform-provider-ct
	@mkdir -p $(DEST)
	@cp _output/$*/terraform-provider-ct $(DEST)/terraform-provider-ct_$(VERSION)
	@zip -j $(DEST).zip $(DEST)/terraform-provider-ct_$(VERSION)
	@mkdir -p $(LOCAL)/$(subst -,_,$*)
	@cp _output/$*/terraform-provider-ct $(LOCAL)/$(subst -,_,$*)/terraform-provider-ct_$(VERSION)

_output/linux-amd64/terraform-provider-ct: GOARGS = GOOS=linux GOARCH=amd64
_output/linux-arm64/terraform-provider-ct: GOARGS = GOOS=linux GOARCH=arm64
_output/darwin-amd64/terraform-provider-ct: GOARGS = GOOS=darwin GOARCH=amd64
_output/darwin-arm64/terraform-provider-ct: GOARGS = GOOS=darwin GOARCH=arm64
_output/windows-amd64/terraform-provider-ct: GOARGS = GOOS=windows GOARCH=amd64
_output/%/terraform-provider-ct:
	$(GOARGS) go build -o $@ github.com/poseidon/terraform-provider-ct

release-sign:
	cd _output; sha256sum *.zip > terraform-provider-ct_$(SEMVER)_SHA256SUMS
	gpg --default-key 0x8F515AD1602065C8 --detach-sign _output/terraform-provider-ct_$(SEMVER)_SHA256SUMS

release-verify: NAME=_output/terraform-provider-ct
release-verify:
	gpg --verify $(NAME)_$(SEMVER)_SHA256SUMS.sig $(NAME)_$(SEMVER)_SHA256SUMS

