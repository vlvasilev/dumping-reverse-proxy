REGISTRY                           := hisshadow85
PROXYE_IMAGE_REPOSITORY            := $(REGISTRY)/proxy
IMAGE_TAG                          := $(shell cat VERSION)


.PHONY: build-proxy
build-proxy:
	@CGO_ENABLED=0 & GOOS=linux & go build  -ldflags "-linkmode external -extldflags -static" -a -o ./bin/proxy ./src/main/main.go

.PHONY: docker-build-proxy
docker-build-proxy:
	@if [ ! -f ./bin/proxy ]; then echo "No binary found. Please run 'make build-proxy'"; false; fi
	@docker build -t $(PROXYE_IMAGE_REPOSITORY):$(IMAGE_TAG) -t $(PROXYE_IMAGE_REPOSITORY):latest -f ./Dockerfile --rm .

.PHONY: docker-push-proxy
docker-push-proxy:
	@if ! docker images $(PROXYE_IMAGE_REPOSITORY) | awk '{ print $$2 }' | grep -q -F $(IMAGE_TAG); then echo "$(PROXYE_IMAGE_REPOSITORY) version $(IMAGE_TAG) is not yet built. Please run 'make docker-build-proxy' or 'make docker-build'"; false; fi
	@docker push $(PROXYE_IMAGE_REPOSITORY):$(IMAGE_TAG)
	@docker push $(PROXYE_IMAGE_REPOSITORY):latest

