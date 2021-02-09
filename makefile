.PHONY: all
.DEFAULT_GOAL := all
VERSION_GITHASH = $(shell git rev-parse master | tr -d '\n')
GO_LDFLAGS = CGO_ENABLED=0 go build -ldflags "-X main.build=${VERSION_GITHASH}" -a -tags netgo

# Should not depend on any other repo / makefiles.
# should be self-contained since this is the bootstrapper of all things.
include ./includes/os.mk
include ./includes/gitr.mk
include ./includes/help.mk

all: this-print this-build
all-ci: all this-dep this-mock-release
all-release: all this-dep this-release

## print all
this-print:
	$(MAKE) os-print
	$(MAKE) gitr-print

this-build:
	$(GO_LDFLAGS) -o bin/booty main.go

this-dep:
	./bin/booty -d -c ./components_version.json install

this-mock-release:
	./bin/booty release --rm-dist --skip-publish --snapshot
	# rm -rf dist

this-release:
	./bin/booty release release
