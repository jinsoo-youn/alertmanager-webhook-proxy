# Makefile for alertmanager-webhook-proxy

APP_NAME := alertmanager-webhook-proxy
BINARY_NAME := $(APP_NAME)
BUILD_DIR := build
GOOS := linux
GOARCH := amd64
VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
HARBOR_REPO := beta-harbor-kr1.cloud.toastoven.net/library
DOCKER_IMAGE := $(HARBOR_REPO)/$(APP_NAME):$(VERSION)

.PHONY: all build clean docker-build run docker-run

all: build

build-for-linux:
	mkdir -p $(BUILD_DIR)/linux
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BUILD_DIR)/linux/$(BINARY_NAME) ./cmd/$(APP_NAME)

build-for-mac:
	mkdir -p $(BUILD_DIR)/mac
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/mac/$(BINARY_NAME) ./cmd/$(APP_NAME)

clean:
	rm -rf $(BUILD_DIR)

run-linux:
	./$(BUILD_DIR)/linux/$(BINARY_NAME)

run-mac:
	./$(BUILD_DIR)/mac/$(BINARY_NAME)

## 도커 허브 등 원격 레지스트리에 푸시 (멀티 아키텍처 지원)
docker-buildx-push:
	docker buildx build \
  		--platform linux/amd64,linux/arm64 \
  		--push \
  		-t $(DOCKER_IMAGE) \
  		.

## 로컬에 저장 (단일 플랫폼일 경우만 사용 가능):
docker-buildx-load:
	docker buildx build \
  		--platform linux/amd64 \
  		--load \
  		-t $(DOCKER_IMAGE) \
  		.

docker-run:
	docker run --rm -it \
		-e LISTEN_ADDRESS="0.0.0.0:8080" \
		-e LOG_LEVEL="info" \
		-e STAGE="beta" \
		-e REGION="kr2" \
		-e WARD_ENABLE="true" \
		-e WARD_EVENT_URL="https://ward-queue.toastmaker.net/event" \
		-e WARD_ACTOR="buoy" \
		-e DOORAY_ENABLE="false" \
		-e DOORAY_WEBHOOK_URL="https://nhnent.dooray.com/services/xxxxx/yyyyy/zzzzz" \
		-p 8080:8080 \
		--name $(APP_NAME) \
		$(DOCKER_IMAGE)
