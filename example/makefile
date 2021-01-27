# boot repo makefiles
BOOTY_FSPATH=./..
include $(BOOTY_FSPATH)/*.mk

all: print dep build 

all-ci : all release-ci

all-release: all release

print:
	@echo -- OS --
	$(MAKE) os-print
	@echo

dep:
	@echo dep
	# TODO: should call dep and do nothing because parent make already did all deps.
	# - can just check if download exists.

build:
	cd ./go/cmd/cli && go build .
	cd ./go/cmd/server && go build .

release-ci:
	#git tag -a $(VERSION) -m "Release" || true
	#git push origin $(VERSION)
	#goreleaser --rm-dist
	goreleaser --rm-dist --skip-publish --snapshot

release:
	goreleaser release	


server-run:
	cd ./go/cmd/server && ./server
cli-run:
	cd ./go/cmd/cli && ./cli





