# grafana
# https://grafana.com/grafana/download

## TODO: We should try to do what tiup does ? 

### BIN
PREFIX=/usr/local/bin
PREFIX=$(GOBIN)

GRAFANA_OUTPUT_DIR=downloaded
GRAFANA_BIN=grafana
# DARWIN: https://dl.grafana.com/oss/release/grafana-7.3.6.darwin-amd64.tar.gz
# WINDOWS: https://dl.grafana.com/oss/release/grafana-7.3.6.windows-amd64.zip
# LINUX: https://dl.grafana.com/oss/release/grafana-7.3.6.linux-amd64.tar.gz
GRAFANA_BIN_VERSION=7.3.6
ifeq ($(GOOS),darwin)
    GRAFANA_BIN_PLATFORM:=darwin-amd64.tar.gz
endif
ifeq ($(GOOS),windows)
    GRAFANA_BIN_PLATFORM:=windows-amd64.zip
endif
ifeq ($(GOOS),linux)
    GRAFANA_BIN_PLATFORM:=linux-amd64.tar.gz
endif
GRAFANA_BIN_FILE=$(GRAFANA_BIN_PLATFORM)
GRAFANA_BIN_URL=https://dl.grafana.com/oss/release/grafana-$(GRAFANA_BIN_VERSION).$(GRAFANA_BIN_FILE)


grafana-print:
	@echo -- Grafana Dep --
	@echo GRAFANA_BIN_VERSION: 	$(GRAFANA_BIN_VERSION)
	@echo GRAFANA_BIN_URL: 		$(GRAFANA_BIN_URL)
	@echo GRAFANA_BIN_URLFILE: 	$(GRAFANA_BIN_FILE)
	@echo GRAFANA_BIN: 			$(GRAFANA_BIN)
	@echo
	
grafana-dep: grafana-dep-delete
	$(MAKE) DWN_URL=$(GRAFANA_BIN_URL) \
		DWN_FILENAME=$(GRAFANA_BIN_FILE) \
 		DWN_BIN_NAME=$(GRAFANA_BIN) \
 		DWN_BIN_OUTPUT_DIR=$(GRAFANA_OUTPUT_DIR) dwn-download

	if [[ $(GOOS) = darwin || $(GOOS) = linux ]]; then \
  		sudo install -m755 $(GRAFANA_OUTPUT_DIR)/$(GRAFANA_BIN) $(PREFIX)/$(GRAFANA_BIN); \
  	fi

	if [[ $(GOOS) = windows ]]; then \
		sudo install -m755 $(GRAFANA_OUTPUT_DIR)/$(GRAFANA_BIN) $(PREFIX)/$(GRAFANA_BIN); \
	fi

grafana-dep-delete:
	$(MAKE) DWN_URL=$(GRAFANA_BIN_URL) \
		DWN_FILENAME=$(GRAFANA_BIN_FILE) \
		DWN_BIN_NAME=$(GRAFANA_BIN) \
		DWN_BIN_OUTPUT_DIR=$(GRAFANA_OUTPUT_DIR) dwn-delete	

	if [[ $(GOOS) = darwin || $(GOOS) = linux ]]; then \
		sudo rm -rf $(PREFIX)/$(GRAFANA_BIN); \
	fi

	if [[ $(GOOS) = windows ]]; then \
		sudo rm -rf $(PREFIX)/$(GRAFANA_BIN); \
	fi


	