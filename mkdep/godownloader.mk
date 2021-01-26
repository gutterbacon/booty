# godownloader (gor)
# https://github.com/goreleaser/godownloader

### BIN
GOD_BIN=godownloader

GOD_BIN_VERSION=0.149.0
ifeq ($(GOOS),darwin)
    GOD_BIN_PLATFORM:=Darwin_x86_64.tar.gz
endif
ifeq ($(GOOS),windows)
    echo 'WINDOWS IS NOT SUPPORTED'
endif
ifeq ($(GOOS),linux)
    GOD_BIN_PLATFORM:=Linux_x86_64.tar.gz
endif
GOD_BIN_FILE=godownloader_$(GOD_BIN_PLATFORM)
GOD_BIN_URL=https://github.com/goreleaser/godownloader/releases/download/v$(GOD_BIN_VERSION)/$(GOD_BIN_FILE)


gor-print:
	@echo -- GORELEASE Dep --
	@echo GOD_BIN_VERSION: 	$(GOD_BIN_VERSION)
	@echo GOD_BIN_URL: 		$(GOD_BIN_URL)
	@echo GOD_BIN_URLFILE: 	$(GOD_BIN_FILE)
	@echo GOD_BIN: 			$(GOD_BIN)
	@echo
	
gor-dep: gor-dep-delete
	$(MAKE) DWN_URL=$(GOD_BIN_URL) \
		DWN_FILENAME=$(GOD_BIN_FILE) \
 		DWN_BIN_NAME=$(GOD_BIN) \
 		DWN_BIN_OUTPUT_DIR=$(INSTALL_DOWNLOAD_DIR) dwn-download

gor-dep-delete:
	$(MAKE) DWN_URL=$(GOD_BIN_URL) \
		DWN_FILENAME=$(GOD_BIN_FILE) \
		DWN_BIN_NAME=$(GOD_BIN) \
		DWN_BIN_OUTPUT_DIR=$(INSTALL_DOWNLOAD_DIR) dwn-delete


	