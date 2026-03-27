SHELL := /bin/bash

# ⚙️ Configuration
APP             ?= watersystem-data-pipeline
DOCKER_COMPOSE  := COMPOSE_BAKE=true docker compose

DOCKERFILE ?= Dockerfile

IMAGE_REG  ?= ghcr.io/bruli
IMAGE_NAME := $(IMAGE_REG)/$(APP)
VERSION    ?= 0.5.0
CURRENT_IMAGE := $(IMAGE_NAME):$(VERSION)

GOLANGCI_LINT_VERSION ?= v2.11.4

# Default goal
.DEFAULT_GOAL := help


# ────────────────────────────────────────────────────────────────
# 🐳 Docker
# ────────────────────────────────────────────────────────────────
.PHONY: docker-up
docker-up:
	@set -euo pipefail; \
	echo "🚀 Starting services with Docker Compose..."; \
	$(DOCKER_COMPOSE) up -d --build

.PHONY: docker-down
docker-down:
	@set -euo pipefail; \
	echo "🛑 Stopping and removing Docker Compose services..."; \
	$(DOCKER_COMPOSE) down

.PHONY: docker-ps
docker-ps:
	@set -euo pipefail; \
	echo "📋 Active services:"; \
	$(DOCKER_COMPOSE) ps

.PHONY: docker-exec
docker-exec:
	@set -euo pipefail; \
	echo "🔎 Opening shell inside ..."; \
	$(DOCKER_COMPOSE) exec $(APP) sh

.PHONY: docker-logs
docker-logs:
	@set -euo pipefail; \
	echo "👀 Showing logs for container $(APP) (CTRL+C to exit)..."; \
	docker logs -f $(APP)

# ────────────────────────────────────────────────────────────────
# 🧹 Code quality: format, lint, tests
# ────────────────────────────────────────────────────────────────
.PHONY: fmt
fmt:
	@set -euo pipefail; \
	echo "👉 Formatting code with gofumpt..."; \
	go tool gofumpt -w .

.PHONY: security
security:
	@set -euo pipefail; \
	echo "👉 Check security"; \
	go tool govulncheck ./...

.PHONY: install-lint
install-lint:
	@set -euo pipefail; \
    echo "🔧 Installing golangci-lint $(GOLANGCI_LINT_VERSION)..."; \
    	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

.PHONY: lint
lint: install-lint
	@set -euo pipefail; \
	echo "🚀 Executing golangci-lint..."; \
    golangci-lint run ./...

.PHONY: test
test:
	@set -euo pipefail; \
	echo "🧪 Running unit tests (race, JSON → tparse)..."; \
	go test -race ./... -json -cover -coverprofile=coverage.out| go tool tparse -all

.PHONY: check
check: fmt security lint test

.PHONY: pub-weather
pub-weather:
	@set -euo pipefail; \
    echo "🚀 Publishing weather events..."; \
    docker exec nats-box /scripts/publish.sh weather

.PHONY: pub-logs
pub-logs:
	@set -euo pipefail; \
    echo "🚀 Publishing logs events..."; \
    docker exec nats-box /scripts/publish.sh log

.PHONY: stream-clean
stream-clean:
	@set -euo pipefail; \
    echo "🧹 Cleaning events..."; \
    docker exec nats-box /scripts/publish.sh clean

.PHONY: clean
clean:
	@set -euo pipefail; \
	echo "🧹 Cleaning local artifacts..."; \
	rm -rf bin dist coverage .*cache || true; \
	go clean -testcache

.PHONY: docker-login
docker-login:
	echo "🔐 Logging into Docker registry...";
	echo "$$CR_PAT" | docker login ghcr.io -u bruli --password-stdin

.PHONY: docker-push-image
docker-push-image: docker-login
	echo "🐳 Building and pushing Docker image $(CURRENT_IMAGE) ...";
	docker buildx build \
		--platform linux/arm64 \
		-t $(CURRENT_IMAGE) \
		-f $(DOCKERFILE) \
		--push \
		.
	 echo "✅ Image $(CURRENT_IMAGE) pushed successfully."

# ────────────────────────────────────────────────────────────────
# ℹ️ Help
# ────────────────────────────────────────────────────────────────
help:
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:' Makefile | awk -F':' '{print "  - " $$1}'
