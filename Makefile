VERSION := $(shell git describe --tags --always --dirty)
PACKAGE := github.com/LinPr/lazys3/cmd

GOFILES := $(shell fd --extension go)

.PHONY: default build lint lint-fix test test-unit test-e2e test-e2e-real \
        fmt check-nerd-font fix-nerd-font docs-prepare docs-dev docs-build docs-preview

default: build
	./lazys3 --debug

build:
	go build -ldflags "-X $(PACKAGE).Version=$(VERSION)" .

lint:
	GOEXPERIMENT= golangci-lint run --path-mode=abs --config=".golangci.yml" --timeout=5m

lint-fix:
	GOEXPERIMENT= golangci-lint run --path-mode=abs --config=".golangci.yml" --timeout=5m --fix

test: test-unit

test-unit:
	go test -v $$(go list ./... | grep -v /e2e)

test-e2e:
	go test -v -tags=e2e ./e2e/...

test-e2e-real:
	LAZYS3_E2E_REAL=oss go test -v -tags=e2e ./e2e/...

fmt:
	go fmt ./...

check-nerd-font:
	nerdfix check $(GOFILES)

fix-nerd-font:
	nerdfix fix --format=json $(GOFILES)

docs-prepare:
	cd docs && pnpm i

docs-dev: docs-prepare
	cd docs && pnpm dev

docs-build: docs-prepare
	cd docs && pnpm build

docs-preview: docs-build
	cd docs && pnpm preview
