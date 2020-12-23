# os tools

export GO_OS 		:= $(shell go env GOOS)
export GO_ARCH		?= $(shell go env GOARCH)

export OS_detected	:=

# other
ifeq '$(findstring ;,$(PATH))' ';'
    OS_detected := Windows
else
    OS_detected := $(shell uname 2>/dev/null || echo Unknown)
    OS_detected := $(patsubst CYGWIN%,Cygwin,$(OS_detected))
    OS_detected := $(patsubst MSYS%,MSYS,$(OS_detected))
    OS_detected := $(patsubst MINGW%,MSYS,$(OS_detected))
endif


## Prints the OS settings
os-print:
	@echo
	@echo -- OS --
	@echo GO_OS: $(GO_OS)
	@echo GO_ARCH: $(GO_ARCH)
	@echo
	@echo -- OS --
	@echo OS_detected: $(GO_ARCH)
	@echo




