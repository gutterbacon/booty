.DEFAULT_GOAL := all

# OS.mk must be explicitly first...
include ./os.mk
include ./help.mk
include ./gitr.mk
include ./mkdep/*.mk

#all: print dep example-build example-release

all: print dep

print:
	$(MAKE) os-print
	$(MAKE) gitr-print

dep:
	$(MAKE) dwn-print
	$(MAKE) gor-dep

example-build:
	cd ./example && $(MAKE) all

example-release:
	cd ./example && $(MAKE) release

