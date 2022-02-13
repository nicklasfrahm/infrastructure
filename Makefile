TARGETS		:= $(shell ls -d cmd/* | sed -e 's!cmd/!bin/!')
SOURCES		:= $(shell find . -name "*.go")
PLATFORM	?= $(shell go version | cut -d " " -f 4)
GOOS		:= $(shell echo $(PLATFORM) | cut -d "/" -f 1)
GOARCH		:= $(shell echo $(PLATFORM) | cut -d "/" -f 2)
SUFFIX		:= $(GOOS)-$(GOARCH)
BUILD_FLAGS	:= -ldflags="-s -w"

# Adjust the binary name on Windows.
ifeq ($(GOOS),windows)
SUFFIX	= $(GOOS)-$(GOARCH).exe
endif

build: $(TARGETS)

bin/%: TARGET=$*
bin/%: $(SOURCES)
	@mkdir -p $(@D)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(BUILD_FLAGS) -o $@-$(SUFFIX) cmd/$(TARGET)/*.go
ifdef UPX
	upx -qq $(UPX) $@-$(SUFFIX)
endif

.PHONY: clean
clean:
	@rm -rvf bin
