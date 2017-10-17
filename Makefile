COMMIT = $$(git describe --always)
PKG = github.com/yuuki/mkr-check
PKGS = $$(go list ./... | grep -v vendor)

all: build

.PHONY: build
build:
	go build -ldflags "-X main.GitCommit=\"$(COMMIT)\"" $(PKG)

.PHONY: test
test: vet
	go test -v $(PKGS)

.PHONY: vet
vet:
	go vet $(PKGS)

.PHONY: lint
lint:
	golint $(PKGS)

.PHONY: patch
patch:
	scripts/bump_version.sh patch

.PHONY: minor
minor:
	scripts/bump_version.sh minor
