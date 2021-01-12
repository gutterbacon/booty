LIB=https://github.com/termbacktime/termbacktime
LIB_FOLDER=$(PWD)/termbacktime


BIN=termshot
BIN_FOLDER=$(PWD)/bin/$(BIN)

print:
		@echo 
	@echo LIB: $(LIB)
	@echo LIB_FOLDER: $(LIB_FOLDER)
	@echo 
	@echo BIN: $(BIN)
	@echo BIN_FOLDER: $(BIN_FOLDER)
	@echo

dep:
	git clone $(LIB)
dep-delete:
	rm -rf $(LIB_FOLDER)

build:
	cd $(LIB_FOLDER) && go build -o $(BIN_FOLDER) .

run:
	$(BIN_FOLDER)

record:
	# needs github auth first
	#$(BIN) auth
	$(BIN_FOLDER) auth --set-token 6f39ab1a713a71cea18b8d0c8e6eefa4b332005a99

	$(BIN_FOLDER) record --open
	
live:
	# does not use github at all
	$(BIN_FOLDER) live




	