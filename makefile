.DEFAULT_GOAL := all

include *.mk

all: print example-build example-release

print:
	$(MAKE) gitr-print

	$(MAKE) os-print

dep:
	@echo fake dep ...

example-build:
	cd ./example && $(MAKE) all

example-release:
	cd ./example && $(MAKE) release

