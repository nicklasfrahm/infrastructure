REGISTRY	:= ghcr.io
REPO		:= nicklasfrahm/infrastructure
TARGET		?= cloudapi
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
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(BUILD_FLAGS) -o $(BINARY) cmd/$(TARGET)/main.go


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

.PHONY: cluster-addons
cluster-addons: ingress dashboard

.PHONY: ingress
ingress:
	kubectl create ns ingress --dry-run=client -o yaml | kubectl apply -f -
	helm upgrade ingress charts/ingress --namespace ingress --install --atomic

.PHONY: dashboard
dashboard:
	kubectl create ns dashboard --dry-run=client -o yaml | kubectl apply -f -
	helm upgrade dashboard charts/dashboard --namespace dashboard --install --atomic --values secret-dashboard.yaml
