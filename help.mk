# help

# Note assumes AWK is installed

.DEFAULT_GOAL       := help-print
HELP_TARGET_MAX_CHAR_NUM := 20

HELP_GREEN  := $(shell tput -Txterm setaf 2)
HELP_YELLOW := $(shell tput -Txterm setaf 3)
HELP_WHITE  := $(shell tput -Txterm setaf 7)
HELP_RESET  := $(shell tput -Txterm sgr0)

.PHONY: help help-*

# Print help
help-print:

	@echo ''
	@echo 'Usage:'
	@echo '  ${HELP_YELLOW}make${HELP_RESET} ${HELP_GREEN}<target>${HELP_RESET}'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  ${HELP_YELLOW}%-$(HELP_TARGET_MAX_CHAR_NUM)s${HELP_RESET} ${HELP_GREEN}%s${HELP_RESET}\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

define help-print-target
    @printf "Executing target: \033[36m$@\033[0m\n"
endef

help-log:
	@echo '${HELP_YELLOW} $(NAME) ${HELP_RESET} ${HELP_GREEN} $(VALUE) ${HELP_RESET}'

