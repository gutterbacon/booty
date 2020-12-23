# boot repo makefiles
BOOTY_FSPATH=./../
include $(BOOTY_FSPATH)/*.mk

all: print dep build

print:
	@echo -- OS --
	$(MAKE) os-print
	@echo

dep:
	@echo dep
	# TODO: should call dep and do nothing because parent make already did all deps.

build:
	cd ./main/cmd/cli && go build .
	cd ./main/cmd/server && go build .
server-run:
	cd ./main/cmd/server && ./server
cli-run:
	cd ./main/cmd/cli && ./cli


