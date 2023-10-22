REGISTRY	:= ghcr.io
REPO			:= nicklasfrahm/infrastructure
TARGET		?= metal
SOURCES		:= $(shell find . -name "*.go")
PLATFORM	?= $(shell go version | cut -d " " -f 4)
GOOS		:= $(shell echo $(PLATFORM) | cut -d "/" -f 1)
GOARCH		:= $(shell echo $(PLATFORM) | cut -d "/" -f 2)
SUFFIX		:= $(GOOS)-$(GOARCH)
VERSION		?= $(shell git describe --always --tags --dirty)
BUILD_FLAGS	:= -ldflags="-s -w -X main.version=$(VERSION)"

ifeq ($(GOOS),windows)
SUFFIX	= $(GOOS)-$(GOARCH).exe
endif

BINARY	?= bin/$(TARGET)-$(SUFFIX)

build: bin/$(TARGET)-$(SUFFIX)

bin/$(TARGET)-$(SUFFIX): $(SOURCES)
	@mkdir -p $(@D)
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(BUILD_FLAGS) -o $(BINARY) cmd/$(TARGET)/*


.PHONY: docker
docker:
	docker build \
	  -t $(TARGET):latest \
	  -t $(TARGET):$(VERSION) \
	  -t $(REPO)-$(TARGET):latest \
	  -t $(REPO)-$(TARGET):$(VERSION) \
	  -t $(REGISTRY)/$(REPO)-$(TARGET):latest \
	  -t $(REGISTRY)/$(REPO)-$(TARGET):$(VERSION) \
	  --build-arg TARGET=$(TARGET) \
	  --build-arg VERSION=$(VERSION) \
	  -f build/package/Dockerfile .

.PHONY: docker-push
docker-push: docker
	docker push $(REGISTRY)/$(REPO)-$(TARGET):latest

.PHONY: deploy
deploy:
	sed -i "s|image: .*|image: $(REGISTRY)/$(REPO)-$(TARGET):$(VERSION)|" deploy/kubectl/api/$(TARGET).yaml
	kubectl apply -f deploy/kubectl/api/$(TARGET).yaml
	git reset --hard
