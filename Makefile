BINDIR := $(CURDIR)/bin
DIST_DIRS := find * -type d -exec
TARGETS := linux/amd64
TARGET_OBJS ?= linux-amd64.tar.gz
BINNAME := ashhellwiggo

GOPATH = $(shell go env GOPATH)
MOD = $(GOPATH)/bin/dep
GOX = $(GOPATH)/bin/gox
GOIMPORTS = $(GOPATH)/bin/goimports
ARCH = $(shell uname -p)

PKG := ./...
TAGS :=
TESTS := .
TESTFLAGS :=
LDFLAGS := -w -s
GOFLAGS :=
SRC := $(shell find . -type f -name '*.go' -print)

SHELL = /usr/bin/env bash

GIT_COMMIT = $(shell git rev-parse HEAD)
GIT_SHA = $(shell git rev-parse --short HEAD)
GIT_TAG = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
GIT_DIRTY = $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")

ifdef VERSION
	BINARY_VERSION = $(VERSION)
endif
BINARY_VERSION ?= ${GIT_TAG}

ifneq ($(GIT_TAG),)
	VERSION_METADATA =
endif

LDFLAGS += -X github.com/ashellwig/ashhellwig-go/internal/version.metadata=${VERSION_METADATA}
LDFLAGS += -X github.com/ashellwig/ashhellwig-go/internal/version.gitCommit=${GIT_COMMIT}
LDFLAGS += -X github.com/ashellwig/ashhellwig-go/internal/version.gitTreeState=${GIT_DIRTY}

.PHONY: all
all: build

# --- Build ---
.PHONY: build
build: $(BINDIR)/$(BINNAME)

$(BINDIR)/$(BINNAME): $(SRC)
	GO111MODULE=on go build $(GOFLAGS) -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -o $(BINDIR)/$(BINNAME) ./cmd/ashhellwiggo

# --- Test ---


# --- Dependencies ---
$(GOX):
	(cd /; GO111MODULE=on go get -u github.com/mitchellh/gox)
$(MOD):
	(cd /; GO111MODULE=on go get -u golang.org/x/tools/cmd/gomod)
$(GOIMPORTS):
	(cd /; GO111MODULE=on go get -u golang.org/x/tools/cmd/goimports)

# --- Release ---
.PHONY: build-release
build-release: LDFLAGS += -extldflags "-static"
build-release: $(GOX)
	GO111MODULE=on CGO_ENABLED=0 $(GOX) -parallel=3 -output="_dist/linux-amd64/$(BINNAME)" -osarch='$(TARGETS)' $(GOFLAGS) -tags '$(TAGS)' -ldflags '$(LDFLAGS)' ./cmd/ashhellwiggo

.PHONY: dist
dist:
	( \
		cd _dist && \
		$(DIST_DIRS) cp ../LICENSE {} \; && \
		$(DIST_DIRS) cp ../README.md {} \; && \
		$(DIST_DIRS) tar -zcf ashhellwiggo-${VERSION}-{}.tar.gz {} \; && \
		$(DIST_DIRS) zip -r ashhellwiggo-${VERSION}-{}.zip {} \; \
	)

.PHONY: fetch-dist
fetch-dist:
	mkdir -p _dist
	cd _dist && \
	for obj in ${TARGET_OBJS}; do \
		curl -sSL -o ashhellwiggo-${VERSION}-$${obj} https://api.github.com/repos/ashellwig/ashhellwig-go/releases/latest | grep "browser_download_url.*tar.gz" | cut -d '"' -f 4 | wget -qi - ; \
	done

.PHONY: sign
sign:
	for f in _dist/*.{gz,zip,sha256,sha256sum} ; do \
		gpg --armor --detach-sign $${f} ; \
	done

.PHONY: checksum
checksum:
	for f in _dist/*.{gz,zip} ; do \
		shasum -a 256 "$${f}" | sed 's/_dist\///' > "$${f}.sha256sum" ; \
		shasum -a 256 "$${f}" | awk '{print $$1}' > "$${f}.sha256" ; \
	done

# --- Clean ---
.PHONY: clean
clean:
	@rm -rf $(BINDIR) ./_dist

.PHONY: release-notes
release-notes:
		@if [ ! -d "./_dist" ]; then \
			echo "please run 'make fetch-release' first" && \
			exit 1; \
		fi
		@if [ -z "${PREVIOUS_RELEASE}" ]; then \
			echo "please set PREVIOUS_RELEASE environment variable" \
			&& exit 1; \
		fi

		@./scripts/release-notes.sh ${PREVIOUS_RELEASE} ${VERSION}


.PHONY: info
info:
	 @echo "Version:           ${VERSION}"
	 @echo "Git Tag:           ${GIT_TAG}"
	 @echo "Git Commit:        ${GIT_COMMIT}"
	 @echo "Git Tree State:    ${GIT_DIRTY}"