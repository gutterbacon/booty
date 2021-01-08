# protoc (pro)
# https://github.com/protocolbuffers/protobuf/releases

### BIN
PRO_BIN=protoc
# https://github.com/protocolbuffers/protobuf/releases/tag/v3.14.0
# https://github.com/protocolbuffers/protobuf/releases/download/v3.14.0/protoc-3.14.0-win64.zip
# https://github.com/protocolbuffers/protobuf/releases/download/v3.14.0/protoc-3.14.0-osx-x86_64.zip
# https://github.com/protocolbuffers/protobuf/releases/download/v3.14.0/protoc-3.14.0-linux-x86_64.zip
PRO_BIN_VERSION=3.14.0
ifeq ($(GOOS),darwin)
    PRO_BIN_PLATFORM:=osx-x86_64.zip
endif
ifeq ($(GOOS),windows)
    PRO_BIN_PLATFORM:=win64.zip
endif
ifeq ($(GOOS),linux)
    PRO_BIN_PLATFORM:=linux-x86_64.zip
endif
PRO_BIN_FILE=protoc_$(PRO_BIN_VERSION)-$(PRO_BIN_PLATFORM)
PRO_BIN_URL=https://github.com/goreleaser/goreleaser/releases/download/v$(PRO_BIN_VERSION)/$(PRO_BIN_FILE)


pro-print:
	@echo -- PROTOC Dep --
	@echo PRO_BIN_VERSION: 	$(PRO_BIN_VERSION)
	@echo PRO_BIN_URL: 		$(PRO_BIN_URL)
	@echo PRO_BIN_URLFILE: 	$(PRO_BIN_FILE)
	@echo PRO_BIN: 			$(PRO_BIN)
	@echo
	
pro-dep: pro-dep-delete
	$(MAKE) DWN_URL=$(PRO_BIN_URL) \
		DWN_FILENAME=$(PRO_BIN_FILE) \
 		DWN_BIN_NAME=$(PRO_BIN) \
 		DWN_BIN_OUTPUT_DIR=$(INSTALL_DOWNLOAD_DIR) dwn-download

pro-dep-delete:
	$(MAKE) DWN_URL=$(PRO_BIN_URL) \
		DWN_FILENAME=$(PRO_BIN_FILE) \
		DWN_BIN_NAME=$(PRO_BIN) \
		DWN_BIN_OUTPUT_DIR=$(INSTALL_DOWNLOAD_DIR) dwn-delete


	