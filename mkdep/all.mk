dep-all: dep-all-print dep-all-install

# This is the location that all deps use for the Folder to place / install the binary on the OS.
# The actual binary is controlled by each dep make, like goreleaser.mk.
INSTALL_PREFIX=?
ifeq ($(GOOS),darwin)
	INSTALL_PREFIX:=/usr/local/bin
endif
ifeq ($(GOOS),windows)
	INSTALL_PREFIX:=%APPDATA%
endif
ifeq ($(GOOS),linux)
	INSTALL_PREFIX:=/usr/local/bin
endif

## print all tools
dep-all-print:
	@echo
	@echo -- dep-all-print  -- 
	@echo
	@echo INSTALL_PREFIX: 	$(INSTALL_PREFIX)
	$(MAKE) gor-print
	#$(MAKE) grafana-print

## install all tools
dep-all-install:
	@echo
	@echo -- dep-all-install -- 
	@echo
	$(MAKE) dwn-print
	$(MAKE) gor-dep
	#$(MAKE) grafana-dep