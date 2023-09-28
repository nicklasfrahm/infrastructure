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

.PHONY: inception-up
inception-up:
	kubectl -n inception apply -f configs/inception/inception.yaml
	kubectl -n inception wait --for=condition=Ready pod -l app.kubernetes.io/name=k3s-poc
	kubectl -n inception get pods -l app.kubernetes.io/name=k3s-poc -o jsonpath='{.items[0].metadata.name}'
	kubectl -n inception cp $$(kubectl -n inception get pods -l app.kubernetes.io/name=k3s-poc -o jsonpath='{.items[0].metadata.name}'):/etc/rancher/k3s/k3s.yaml kubeconfig.yaml
	sed -i "s/127.0.0.1/k3s-poc.moos.nicklasfrahm.dev/" kubeconfig.yaml

.PHONY: inception-down
inception-down:
	rm kubeconfig.yaml || true
	kubectl delete -f configs/inception/inception.yaml
