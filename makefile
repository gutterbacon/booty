.DEFAULT_GOAL := all

# OS.mk must be explicitly first...
include ./os.mk
include ./help.mk
include ./gitr.mk
include ./mkdep/*.mk

all: this-print this-dep this-example-build this-example-release

#all: print dep

## print all
this-print:
	$(MAKE) os-print
	$(MAKE) gitr-print

## dep all
this-dep:
	$(MAKE) dep-all

## example all
this-example-all:
	cd ./example && $(MAKE) all

## example build
this-example-build: 
	cd ./example && $(MAKE) build

## example release
this-example-release:
	cd ./example && $(MAKE) release

