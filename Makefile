GOOS		?= linux
GOARCH		?= amd64

.PHONY: all
all: bin/microddns

bin/microddns-$(GOOS)-$(GOARCH):
	@mkdir -p bin
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $@ cmd/microddns/*
