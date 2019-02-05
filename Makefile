export CGO_ENABLED:=0
export GO111MODULE=on
export GOFLAGS=-mod=vendor

VERSION=$(shell git describe --tags --match=v* --always --dirty)

.PHONY: all
all: build test vet lint fmt

.PHONY: build
build: clean bin/terraform-provider-ct

bin/terraform-provider-ct:
	@go build -o $@ github.com/coreos/terraform-provider-ct

.PHONY: install
install: bin/terraform-provider-ct
	@cp $< $(GOPATH_BIN)

.PHONY: test
test:
	@go test ./... -cover

.PHONY: vet
vet:
	@go vet -all ./...

.PHONY: lint
lint:
	@golint -set_exit_status `go list ./...`

.PHONY: fmt
fmt:
	@test -z $$(go fmt ./...)

.PHONY: update
update:
	@GOFLAGS="" go get -u
	@go mod tidy

.PHONY: vendor
vendor:
	@go mod vendor

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
