# os tools

#export GO_OS 		:= $(shell go env GOOS)
#export GO_ARCH		?= $(shell go env GOARCH)

# SOURCE: https://github.com/containerd/containerd/blob/master/Makefile#L28
ifneq "$(strip $(shell command -v go 2>/dev/null))" ""
	GOOS ?= $(shell go env GOOS)
	GOARCH ?= $(shell go env GOARCH)
else
	ifeq ($(GOOS),)
		# approximate GOOS for the platform if we don't have Go and GOOS isn't
		# set. We leave GOARCH unset, so that may need to be fixed.
		ifeq ($(OS),Windows_NT)
			GOOS = windows
		else
			UNAME_S := $(shell uname -s)
			ifeq ($(UNAME_S),Linux)
				GOOS = linux
			endif
			ifeq ($(UNAME_S),Darwin)
				GOOS = darwin
			endif
			ifeq ($(UNAME_S),FreeBSD)
				GOOS = freebsd
			endif
		endif
	else
		GOOS ?= $$GOOS
		GOARCH ?= $$GOARCH
	endif
endif

## Prints the OS settings
os-print:
	@echo
	@echo -- OS --
	@echo GOOS: $(GOOS)
	@echo GOARCH: $(GOARCH)
	@echo




