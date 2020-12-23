.DEFAULT_GOAL := all

# OS.mk must be explicitly first...
include ./os.mk
include ./help.mk
include ./gitr.mk
include ./mkdep/*.mk

all: print dep example-build example-release

#all: print dep

## print all
print:
	$(MAKE) os-print
	$(MAKE) gitr-print

## dep all
dep:
	$(MAKE) dep-all-devtime

## example all
example-all:
	cd ./example && $(MAKE) all

## example build
example-build: 
	cd ./example && $(MAKE) build

## example release
example-release:
	cd ./example && $(MAKE) release

