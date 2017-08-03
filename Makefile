export CGO_ENABLED:=0

VERSION=$(shell ./scripts/git-version)
PACKAGES = $(shell go list ./ct)

.PHONY: all
all: bin/terraform-provider-ct

bin/terraform-provider-ct:
	@go build -o $@ -v github.com/coreos/terraform-provider-ct

.PHONY: install
install: bin/terraform-provider-ct
	@cp $< $(GOPATH_BIN)

.PHONY: test
test:
	go vet $(PACKAGES)
	go test -v $(PACKAGES)

.PHONY: vendor
vendor:
	@glide update --strip-vendor
	@glide-vc --use-lock-file --no-tests --only-code

.PHONY: clean
clean:
	@rm -rf bin
	@rm -rf _output

.PHONY: release
release: \
	clean \
	_output/plugin-linux-amd64.tar.gz \
	_output/plugin-darwin-amd64.tar.gz

_output/plugin-%.tar.gz: NAME=terraform-provider-ct-$(VERSION)-$*
_output/plugin-%.tar.gz: DEST=_output/$(NAME)
_output/plugin-%.tar.gz: _output/%/terraform-provider-ct
	@mkdir -p $(DEST)
	@cp _output/$*/terraform-provider-ct $(DEST)
	@tar zcvf $(DEST).tar.gz -C _output $(NAME)

_output/linux-amd64/terraform-provider-ct: GOARGS = GOOS=linux GOARCH=amd64
_output/darwin-amd64/terraform-provider-ct: GOARGS = GOOS=darwin GOARCH=amd64
_output/%/terraform-provider-ct:
	$(GOARGS) go build -o $@ github.com/coreos/terraform-provider-ct
