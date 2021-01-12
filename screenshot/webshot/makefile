LIB=https://github.com/sensepost/gowitness
LIB_FOLDER=$(PWD)/gowitness

BIN=webshot
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
	cd $(LIB_FOLDER) && $(MAKE) generate
	cd $(LIB_FOLDER) && go build -o $(BIN_FOLDER) .


TEST_URL_00=https://news.ycombinator.com/
TEST_URL_01='https://gallery.flutter.dev/#/reply'
TEST_URL_01=https://flutter.github.io/samples/web/form_app/#/

TEST_ARG=single --disable-db --screenshot-path ./screen_test --delay 1

run-test:

	$(BIN_FOLDER) $(TEST_ARG) $(TEST_URL_00)
	$(BIN_FOLDER) $(TEST_ARG) https://gallery.flutter.dev/#/reply
	$(BIN_FOLDER) $(TEST_ARG) https://flutter.github.io/samples/web/form_app/#/
	$(BIN_FOLDER) $(TEST_ARG) https://flutter.github.io/samples/web/chart

FLY_ARG=single --disable-db --screenshot-path ./screen_fly.dev --delay 1

run-fly:
	$(BIN_FOLDER) $(FLY_ARG) https://org-y.fly.dev/#/projects
	$(BIN_FOLDER) $(FLY_ARG) https://org-y.fly.dev/#/settings
 
LOCAL_ARG=scan --threads 20 --ports 80,443,8080 --cidr 192.168.0.0/24 --disable-db --screenshot-path ./screen_local

run-local:
	$(BIN_FOLDER) $(LOCAL_ARG)