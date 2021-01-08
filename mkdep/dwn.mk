# downloads (dwn) releases from github to local
# This is designed to be called from other makefiles that provide the variables.
# THis then just does the work..

# TODO: Cache: Make it check is the download is already there, and is so, dont download it.

SHELL=bash

# variables to overide.
DWN_URL:=getcourage.org	# Github URL to the file
DWN_FILENAME:=hello		# Github FileName
DWN_BIN_NAME:=?			# Local filename (the actual bin)

# calculated private variables
DWN_FILENAME_CALC=$(notdir $(DWN_URL)) # todo use this, so we dont need to pass in this anymore :)
DWN_FILENAME_BASE=$(shell basename -- $(DWN_FILENAME))
DWN_FILENAME_EXT := $(suffix $(DWN_FILENAME))
ifeq ($(DWN_FILENAME_EXT),)
	DWN_FILENAME_EXT += NONE
endif


dwn-print:
	@echo
	@echo -- DWN Downloader --
	@echo DWN_URL: $(DWN_URL)
	@echo DWN_FILENAME: $(DWN_FILENAME)
	@echo DWN_BIN_NAME: $(DWN_BIN_NAME)
	@echo
	@echo -- DWN Downloader calculated variables --
	@echo DWN_FILENAME_CALC: $(DWN_FILENAME_CALC)
	@echo DWN_FILENAME_BASE: $(DWN_FILENAME_BASE)
	@echo DWN_FILENAME_EXT: $(DWN_FILENAME_EXT)
	@echo

dwn-download: dwn-delete

	@echo Downloading dep ...
	@echo
	@echo DEP_DOWNLOAD_DIR: 	$(DEP_DOWNLOAD_DIR)
	@echo DEP_INSTALL_PREFIX: 	$(DEP_INSTALL_PREFIX)
	@echo
	@echo DWN_FILENAME: 		$(DWN_FILENAME)
	@echo DWN_BIN_NAME: 		$(DWN_BIN_NAME)
	@echo
	@echo DWN_FILENAME_CALC: 	$(DWN_FILENAME_CALC)
	@echo DWN_FILENAME_BASE: 	$(DWN_FILENAME_BASE)
	@echo DWN_FILENAME_EXT:	 	$(DWN_FILENAME_EXT)
	@echo

	if [[ $(GOOS) = darwin || $(GOOS) = linux ]]; then \
		mkdir -p $(DEP_DOWNLOAD_DIR) && \
		curl -L -o $(DEP_DOWNLOAD_DIR)/$(DWN_FILENAME_CALC) $(DWN_URL) && \
		if [[ $(DWN_FILENAME_EXT) = .gz || $(DWN_FILENAME_EXT) = .tar.gz ]]; then \
			tar -zvxf $(DEP_DOWNLOAD_DIR)/$(DWN_FILENAME_CALC) -C $(DEP_DOWNLOAD_DIR); \
		elif [[ $(DWN_FILENAME_EXT) = .zip ]]; then \
			unzip -d $(DEP_DOWNLOAD_DIR) $(DEP_DOWNLOAD_DIR)/$(DWN_FILENAME_CALC); \
		else \
			mv $(DEP_DOWNLOAD_DIR)/$(DWN_FILENAME_CALC) $(DEP_DOWNLOAD_DIR)/$(DWN_BIN_NAME); \
		fi && \
		chmod +x $(DEP_DOWNLOAD_DIR)/$(DWN_BIN_NAME); \
	fi

	
	if [[ $(GOOS) = windows ]]; then \
		mkdir -p $(DEP_DOWNLOAD_DIR) && \
		curl -L -o $(DEP_DOWNLOAD_DIR)/$(DWN_FILENAME_CALC) $(DWN_URL) && \
		if [[ $(DWN_FILENAME_EXT) = .gz || $(DWN_FILENAME_EXT) = .tar.gz ]]; then \
			tar -zvxf $(DEP_DOWNLOAD_DIR)/$(DWN_FILENAME_CALC) -C $(DEP_DOWNLOAD_DIR); \
		elif [[ $(DWN_FILENAME_EXT) = .zip ]]; then \
			unzip -d $(DEP_DOWNLOAD_DIR) $(DEP_DOWNLOAD_DIR)/$(DWN_FILENAME_CALC); \
		else \
			mv $(DEP_DOWNLOAD_DIR)/$(DWN_FILENAME_CALC) $(DEP_DOWNLOAD_DIR)/$(DWN_BIN_NAME); \
		fi && \
		chmod +x $(DEP_DOWNLOAD_DIR)/$(DWN_BIN_NAME); \
	fi


	# copy / install into the install location
	if [[ $(GOOS) = darwin || $(GOOS) = linux ]]; then \
  		sudo install -m755 $(DEP_DOWNLOAD_DIR)/$(DWN_BIN_NAME) $(DEP_INSTALL_PREFIX)/$(DWN_BIN_NAME); \
  	fi

	if [[ $(GOOS) = windows ]]; then \
		mkdir -p $(DEP_INSTALL_PREFIX); \
		cp $(DEP_DOWNLOAD_DIR)/$(DWN_BIN_NAME) $(DEP_INSTALL_PREFIX)/$(DWN_BIN_NAME); \
	fi

dwn-delete:

	@echo Deleting dep ...
	@echo
	@echo DEP_DOWNLOAD_DIR: 	$(DEP_DOWNLOAD_DIR)
	@echo DEP_INSTALL_PREFIX: 	$(DEP_INSTALL_PREFIX)
	@echo
	@echo DWN_FILENAME: 		$(DWN_FILENAME)
	@echo DWN_BIN_NAME: 		$(DWN_BIN_NAME)
	@echo
	@echo DWN_FILENAME_CALC: 	$(DWN_FILENAME_CALC)
	@echo DWN_FILENAME_BASE: 	$(DWN_FILENAME_BASE)
	@echo DWN_FILENAME_EXT:	 	$(DWN_FILENAME_EXT)
	@echo

	# deletes the tar and binary
	if [[ $(GOOS) = darwin || $(GOOS) = linux ]]; then \
		rm -rf $(DEP_DOWNLOAD_DIR)/$(DWN_FILENAME); \
		rm -rf $(DEP_DOWNLOAD_DIR)/$(DWN_BIN_NAME); \
	fi

	if [[ $(GOOS) = windows ]]; then \
		rm -rf $(DEP_DOWNLOAD_DIR)/$(DWN_FILENAME); \
		rm -rf $(DEP_DOWNLOAD_DIR)/$(DWN_BIN_NAME); \
	fi


	# delete the binary from the install location
	if [[ $(GOOS) = darwin || $(GOOS) = linux ]]; then \
		sudo rm -rf $(DEP_INSTALL_PREFIX)/$(DWN_BIN_NAME); \
	fi

	if [[ $(GOOS) = windows ]]; then \
		rm -rf $(DEP_INSTALL_PREFIX)/$(DWN_BIN_NAME); \
	fi
