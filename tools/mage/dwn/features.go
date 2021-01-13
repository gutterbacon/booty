package dwn

// Delete deletes all downloaded dependencies
func Delete() {

	// dwn-delete:

	// @echo Deleting dep ...
	// @echo
	// @echo DEP_DOWNLOAD_DIR: 	$(DEP_DOWNLOAD_DIR)
	// @echo DEP_INSTALL_PREFIX: 	$(DEP_INSTALL_PREFIX)
	// @echo
	// @echo DWN_FILENAME: 		$(DWN_FILENAME)
	// @echo DWN_BIN_NAME: 		$(DWN_BIN_NAME)
	// @echo
	// @echo DWN_FILENAME_CALC: 	$(DWN_FILENAME_CALC)
	// @echo DWN_FILENAME_BASE: 	$(DWN_FILENAME_BASE)
	// @echo DWN_FILENAME_EXT:	 	$(DWN_FILENAME_EXT)
	// @echo

	// # deletes the tar and binary
	// if [[ $(GOOS) = darwin || $(GOOS) = linux ]]; then \
	// 	rm -rf $(DEP_DOWNLOAD_DIR)/$(DWN_FILENAME); \
	// 	rm -rf $(DEP_DOWNLOAD_DIR)/$(DWN_BIN_NAME); \
	// fi

	// if [[ $(GOOS) = windows ]]; then \
	// 	rm -rf $(DEP_DOWNLOAD_DIR)/$(DWN_FILENAME); \
	// 	rm -rf $(DEP_DOWNLOAD_DIR)/$(DWN_BIN_NAME); \
	// fi

	// # delete the binary from the install location
	// if [[ $(GOOS) = darwin || $(GOOS) = linux ]]; then \
	// 	sudo rm -rf $(DEP_INSTALL_PREFIX)/$(DWN_BIN_NAME); \
	// fi

	// if [[ $(GOOS) = windows ]]; then \
	// 	rm -rf $(DEP_INSTALL_PREFIX)/$(DWN_BIN_NAME); \
	// fi

}

// Download downloads all dependencies
func Download() {
	// dwn-download: dwn-delete

	// @echo Downloading dep ...
	// @echo
	// @echo DEP_DOWNLOAD_DIR: 	$(DEP_DOWNLOAD_DIR)
	// @echo DEP_INSTALL_PREFIX: 	$(DEP_INSTALL_PREFIX)
	// @echo
	// @echo DWN_FILENAME: 		$(DWN_FILENAME)
	// @echo DWN_BIN_NAME: 		$(DWN_BIN_NAME)
	// @echo
	// @echo DWN_FILENAME_CALC: 	$(DWN_FILENAME_CALC)
	// @echo DWN_FILENAME_BASE: 	$(DWN_FILENAME_BASE)
	// @echo DWN_FILENAME_EXT:	 	$(DWN_FILENAME_EXT)
	// @echo

	// if [[ $(GOOS) = darwin || $(GOOS) = linux ]]; then \
	// 	mkdir -p $(DEP_DOWNLOAD_DIR) && \
	// 	curl -L -o $(DEP_DOWNLOAD_DIR)/$(DWN_FILENAME_CALC) $(DWN_URL) && \
	// 	if [[ $(DWN_FILENAME_EXT) = .gz || $(DWN_FILENAME_EXT) = .tar.gz ]]; then \
	// 		tar -zvxf $(DEP_DOWNLOAD_DIR)/$(DWN_FILENAME_CALC) -C $(DEP_DOWNLOAD_DIR); \
	// 	elif [[ $(DWN_FILENAME_EXT) = .zip ]]; then \
	// 		unzip -d $(DEP_DOWNLOAD_DIR) $(DEP_DOWNLOAD_DIR)/$(DWN_FILENAME_CALC); \
	// 	else \
	// 		mv $(DEP_DOWNLOAD_DIR)/$(DWN_FILENAME_CALC) $(DEP_DOWNLOAD_DIR)/$(DWN_BIN_NAME); \
	// 	fi && \
	// 	chmod +x $(DEP_DOWNLOAD_DIR)/$(DWN_BIN_NAME); \
	// fi

	// if [[ $(GOOS) = windows ]]; then \
	// 	mkdir -p $(DEP_DOWNLOAD_DIR) && \
	// 	curl -L -o $(DEP_DOWNLOAD_DIR)/$(DWN_FILENAME_CALC) $(DWN_URL) && \
	// 	if [[ $(DWN_FILENAME_EXT) = .gz || $(DWN_FILENAME_EXT) = .tar.gz ]]; then \
	// 		tar -zvxf $(DEP_DOWNLOAD_DIR)/$(DWN_FILENAME_CALC) -C $(DEP_DOWNLOAD_DIR); \
	// 	elif [[ $(DWN_FILENAME_EXT) = .zip ]]; then \
	// 		unzip -d $(DEP_DOWNLOAD_DIR) $(DEP_DOWNLOAD_DIR)/$(DWN_FILENAME_CALC); \
	// 	else \
	// 		mv $(DEP_DOWNLOAD_DIR)/$(DWN_FILENAME_CALC) $(DEP_DOWNLOAD_DIR)/$(DWN_BIN_NAME); \
	// 	fi && \
	// 	chmod +x $(DEP_DOWNLOAD_DIR)/$(DWN_BIN_NAME); \
	// fi

	// # copy / install into the install location
	// if [[ $(GOOS) = darwin || $(GOOS) = linux ]]; then \
	// 	sudo install -m755 $(DEP_DOWNLOAD_DIR)/$(DWN_BIN_NAME) $(DEP_INSTALL_PREFIX)/$(DWN_BIN_NAME); \
	// fi

	// if [[ $(GOOS) = windows ]]; then \
	// 	mkdir -p $(DEP_INSTALL_PREFIX); \
	// 	cp $(DEP_DOWNLOAD_DIR)/$(DWN_BIN_NAME) $(DEP_INSTALL_PREFIX)/$(DWN_BIN_NAME); \
	// fi
}
