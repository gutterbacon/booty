include *.mk

all: print example-build

print:
	$(MAKE) gitr-print

	$(MAKE) os-print

dep:
	@echo fake dep ...

example-build:
	cd ./example && $(MAKE) all