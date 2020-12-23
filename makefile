.DEFAULT_GOAL := all

include *.mk

all: print dep example-build example-release

print:
	$(MAKE) gitr-print

	$(MAKE) os-print

dep:
	$(MAKE) gor-dep

example-build:
	cd ./example && $(MAKE) all

example-release:
	cd ./example && $(MAKE) release

