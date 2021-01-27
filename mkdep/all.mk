dep-all: dep-all-print dep-all-install

# This is the location that all deps use for the Folder to place / install the binary on the OS.
# The actual binary is controlled by each dep make, like goreleaser.mk.
DEP_INSTALL_PREFIX=?
ifeq ($(GOOS),darwin)
	DEP_INSTALL_PREFIX:=/usr/local/bin
endif
# For Windows we use a bash shell, and so Install Prefix is : https://medium.com/prodopsio/3rd-party-binaries-in-windows-git-bash-d882905e470c
ifeq ($(GOOS),windows)
	DEP_INSTALL_PREFIX:=/usr/bin
endif
ifeq ($(GOOS),linux)
	DEP_INSTALL_PREFIX:=/usr/local/bin
endif

DEP_DOWNLOAD_DIR=downloaded


## print all tools
dep-all-print:
	@echo
	@echo -- dep-all-print  -- 
	@echo
	@echo DEP_INSTALL_PREFIX: 	$(DEP_INSTALL_PREFIX)
	# where is home for each OS ?
	@echo HOME_ENV: 		$${HOME}
	$(MAKE) gor-print
	$(MAKE) pro-print
	#$(MAKE) grafana-print

## install all tools
dep-all-install:
	@echo
	@echo -- dep-all-install -- 
	@echo
	$(MAKE) dwn-print
	$(MAKE) gor-dep
	#$(MAKE) grafana-dep