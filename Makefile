TARGET			?= ic
PLATFORM		?= $(shell go version | cut -d " " -f 4)
VERSION			?= $(shell git describe --always --tags --dirty)
REPO				:= $(shell git remote get-url origin | sed 's|https://github.com||g' |  grep -oE '[a-z-]+/[a-z-]+')
SOURCES			:= $(shell find . -name "*.go")
GOOS				:= $(shell echo $(PLATFORM) | cut -d "/" -f 1)
GOARCH			:= $(shell echo $(PLATFORM) | cut -d "/" -f 2)
SUFFIX			:= $(GOOS)-$(GOARCH)
BUILD_FLAGS	:= -ldflags="-s -w -X main.version=$(VERSION)"
BIN_DIR			:= /usr/local/bin
REGISTRY		:= ghcr.io

ifeq ($(GOOS),windows)
SUFFIX	= $(GOOS)-$(GOARCH).exe
endif

BINARY			?= bin/$(TARGET)-$(SUFFIX)

build: $(BINARY)

$(BINARY): $(SOURCES)
	@mkdir -p $(@D)
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(BUILD_FLAGS) -o $@ cmd/$(TARGET)/*.go

install: $(BIN_DIR)/$(TARGET)

$(BIN_DIR)/$(TARGET): $(BINARY)
	sudo install -m 0755 $^ $@

uninstall:
	sudo rm -f $(BIN_DIR)/$(TARGET)

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
	docker push $(REGISTRY)/$(REPO)-$(TARGET):$(VERSION)
	docker push $(REGISTRY)/$(REPO)-$(TARGET):latest

.PHONY: deploy
deploy: docker-push
	sed -i "s|image: .*|image: $(REGISTRY)/$(REPO)-$(TARGET):$(VERSION)|" deploy/kubectl/api/$(TARGET).yaml
	kubectl apply -n api -f deploy/kubectl/api/$(TARGET).yaml
	git reset --hard

.PHONY: edge
edge: bin/nofip-$(SUFFIX)
	@for SITE in alfa bravo charlie ; do \
		echo "\033[0;31m>> $$SITE\033[0m" ; \
    kubectl --context $$SITE create namespace edge --dry-run=client -o yaml | kubectl apply --server-side -f - ; \
		helm --kube-context $$SITE -n edge upgrade --install --atomic edge charts/edge ; \
	done
	@./bin/nofip-$(SUFFIX) -r edge.nicklasfrahm.dev -e alfa.nicklasfrahm.dev,bravo.nicklasfrahm.dev,charlie.nicklasfrahm.dev

.PHONY: kuard
kuard:
	kubectl --context moos create namespace kuard --dry-run=client -o yaml | kubectl apply --server-side -f -
	kubectl --context moos -n kuard apply -f deploy/kubectl/kuard

.PHONY: odance
odance:
	kubectl --context moos create namespace odance-prd --dry-run=client -o yaml | kubectl --context moos apply --server-side -f -
	kubectl --context moos -n odance-prd apply -f secret-odance-prd.yaml
	kubectl --context moos -n odance-prd apply -f deploy/kubectl/odance/prd.yaml
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm repo update bitnami
	helm --kube-context moos -n odance-prd upgrade --install --atomic odance bitnami/wordpress -f deploy/helm/odance.values.yaml
