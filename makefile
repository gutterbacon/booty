.DEFAULT_GOAL := all

# OS.mk must be explicitly first...
#BOILERPLATE_FSPATH=./../boot/boilerplate
BOILERPLATE_FSPATH=./../shared/boilerplate

include $(BOILERPLATE_FSPATH)/help.mk
include $(BOILERPLATE_FSPATH)/os.mk
include $(BOILERPLATE_FSPATH)/gitr.mk
include $(BOILERPLATE_FSPATH)/tool.mk
include $(BOILERPLATE_FSPATH)/flu.mk
include $(BOILERPLATE_FSPATH)/go.mk

#include ./os.mk
#include ./help.mk
#include ./gitr.mk
#include ./mkdep/*.mk

all: this-print this-dep this-example-build 
all-ci: all this-example-release-ci
all-release: all this-example-release

#all: print dep

## print all
this-print:
	$(MAKE) os-print
	$(MAKE) gitr-print

### Mage

this-mage-dep:
	go install github.com/magefile/mage

this-mage-run:
	# assumes you have mage installed on path
	mage -l

this-mage-build:
	mkdir -p ./dist
	go build -o ./dist/booty mage.go
	# TODO use goxc to build cross platform. I dont think goreleaser will build mage cross platform.
	# Can then release booty as a binary as part of its own CI.



## dep all
this-dep:
	$(MAKE) dep-all



## example all
this-example-all:
	cd ./example && $(MAKE) all

## example build
this-example-build: 
	cd ./example && $(MAKE) build

## example release ci
this-example-release-ci:
	cd ./example && $(MAKE) release-ci

## example release
this-example-release:
	cd ./example && $(MAKE) release