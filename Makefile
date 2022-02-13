TARGETS		:= $(shell ls -d cmd/* | sed -e 's!cmd/!bin/!')
SOURCES		:= $(shell find . -name "*.go")
PLATFORM	?= $(shell go version | cut -d " " -f 4)
GOOS		:= $(shell echo $(PLATFORM) | cut -d "/" -f 1)
GOARCH		:= $(shell echo $(PLATFORM) | cut -d "/" -f 2)

build: $(TARGETS)

bin/%: TARGET=$*
bin/%: $(SOURCES)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $@-$(GOOS)-$(GOARCH) cmd/$(TARGET)/*

.PHONY: clean
clean:
	@rm -rvf bin
