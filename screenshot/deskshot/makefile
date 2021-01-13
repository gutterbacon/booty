#https://github.com/jerson/screenshot-tools

LIB=https://github.com/jerson/screenshot-tools
LIB_FOLDER=$(PWD)/screenshot-tools

BIN=./deskshot
BIN_FOLDER=$(PWD)/bin


print:
	@echo 
	@echo LIB: $(LIB)
	@echo LIB_FOLDER: $(LIB_FOLDER)
	@echo 
	@echo BIN: $(BIN)
	@echo BIN_FOLDER: $(BIN_FOLDER)
	@echo BIN_ARG: $(BIN_ARG)
	@echo

dep:
	git clone $(LIB)

test:
	cd $(LIB_FOLDER) && go test -v .
build:
	# generator embeds various scripts for automator
	cd $(LIB_FOLDER) && go generate .
	cd $(LIB_FOLDER) && go build -o $(BIN) .

run-screenshot:
	# ss = Session name which creates a folder of that name.
	cd $(BIN_FOLDER) && $(BIN) ss "joe" -p macos
run-report:
	# opens TUI screen
	cd $(BIN_FOLDER) && $(BIN) r
	
