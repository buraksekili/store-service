SERVICES = emailer products users
CGO_ENABLED ?= 0
GOARCH ?= amd64
BUILD_DIR = build
DOCKERS = $(addprefix docker_,$(SERVICES))

define compile_service
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) go build -mod=vendor -ldflags "-s -w" -o ${BUILD_DIR}/store-$(1) cmd/$(1)/main.go
endef

define make_docker
	$(eval svc=$(subst docker_,,$(1)))
	docker build \
		--no-cache \
		--build-arg SVC=$(svc) \
		--tag=store-$(svc) \
		-f docker/Dockerfile ./build
endef

all: $(SERVICES)

.PHONY: $(SERVICES)

$(SERVICES):
	$(call compile_service,$(@))

run:
	docker-compose -f docker/docker-compose.yaml up

down:
	docker-compose -f docker/docker-compose.yaml down

build_clean:
	rm -rf $(BUILD_DIR)

$(DOCKERS):
	$(call make_docker,$(@))

images: $(DOCKERS)

containers: $(DOCKERS)

