# Atlas Makefile
# All developer commands are defined here.
# Run `make help` to see available commands.

# ── Variables ──────────────────────────────────────────────────────────────
BINARY      := atlas
CMD         := ./cmd/atlas
BUILD_DIR   := ./build
VERSION     ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS     := -ldflags "-s -w -X main.version=$(VERSION) -X github.com/Suke2004/atlas-go/internal/health.version=$(VERSION)"
CGO_ENABLED := 1

# Colours for output
GREEN  := \033[0;32m
YELLOW := \033[0;33m
RESET  := \033[0m

.DEFAULT_GOAL := help

# ── Help ───────────────────────────────────────────────────────────────────
.PHONY: help
help: ## Show this help message
	@echo ""
	@echo "Atlas — Developer Commands"
	@echo "────────────────────────────────────────"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(RESET) %s\n", $$1, $$2}'
	@echo ""

# ── Development ────────────────────────────────────────────────────────────
.PHONY: dev
dev: ## Start dev server with hot-reload (Air)
	@echo "$(YELLOW)Starting Atlas in development mode...$(RESET)"
	air

.PHONY: build
build: templ ## Build production binary
	@echo "$(YELLOW)Building Atlas $(VERSION)...$(RESET)"
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=$(CGO_ENABLED) go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY) $(CMD)
	@echo "$(GREEN)Built: $(BUILD_DIR)/$(BINARY)$(RESET)"

.PHONY: run
run: build ## Build and run the binary directly (no hot-reload)
	$(BUILD_DIR)/$(BINARY)

# ── Code Generation ────────────────────────────────────────────────────────
.PHONY: templ
templ: ## Regenerate Templ templates
	@echo "$(YELLOW)Generating Templ templates...$(RESET)"
	templ generate

.PHONY: sqlc
sqlc: ## Regenerate sqlc database code from queries/
	@echo "$(YELLOW)Generating sqlc code...$(RESET)"
	sqlc generate

.PHONY: generate
generate: templ sqlc ## Run all code generators

# ── Database ───────────────────────────────────────────────────────────────
.PHONY: migrate
migrate: ## Apply all pending database migrations
	@echo "$(YELLOW)Applying migrations...$(RESET)"
	goose -dir ./migrations sqlite3 ./data/db/atlas.db up

.PHONY: migrate-down
migrate-down: ## Roll back the last migration
	@echo "$(YELLOW)Rolling back last migration...$(RESET)"
	goose -dir ./migrations sqlite3 ./data/db/atlas.db down

.PHONY: migrate-status
migrate-status: ## Show current migration status
	goose -dir ./migrations sqlite3 ./data/db/atlas.db status

.PHONY: migrate-create
migrate-create: ## Create a new migration file (usage: make migrate-create name=add_users)
	@if [ -z "$(name)" ]; then echo "Usage: make migrate-create name=your_migration_name"; exit 1; fi
	goose -dir ./migrations create $(name) sql

.PHONY: seed
seed: ## Run the demo data seed script
	@echo "$(YELLOW)Seeding demo data...$(RESET)"
	go run ./scripts/seed.go

# ── Testing ────────────────────────────────────────────────────────────────
.PHONY: test
test: ## Run all tests with race detector and coverage
	@echo "$(YELLOW)Running tests...$(RESET)"
	CGO_ENABLED=$(CGO_ENABLED) go test ./... -race -coverprofile=coverage.out -covermode=atomic
	@echo "$(GREEN)Tests passed.$(RESET)"

.PHONY: test-unit
test-unit: ## Run unit tests only
	CGO_ENABLED=$(CGO_ENABLED) go test ./tests/unit/... -race -v

.PHONY: test-integration
test-integration: ## Run integration tests only
	CGO_ENABLED=$(CGO_ENABLED) go test ./tests/integration/... -race -v

.PHONY: coverage
coverage: test ## Generate and open HTML coverage report
	go tool cover -html=coverage.out

# ── Code Quality ───────────────────────────────────────────────────────────
.PHONY: lint
lint: ## Run golangci-lint
	@echo "$(YELLOW)Running linter...$(RESET)"
	golangci-lint run ./...

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: fmt
fmt: ## Format all Go code
	gofmt -s -w .

.PHONY: check
check: vet lint test ## Run vet + lint + tests (CI equivalent)
	@echo "$(GREEN)All checks passed.$(RESET)"

# ── Docker ─────────────────────────────────────────────────────────────────
.PHONY: docker-build
docker-build: ## Build the Docker image
	docker compose -f docker/docker-compose.yml build

.PHONY: docker-up
docker-up: ## Start Atlas in Docker
	docker compose -f docker/docker-compose.yml up -d

.PHONY: docker-down
docker-down: ## Stop Atlas Docker containers
	docker compose -f docker/docker-compose.yml down

.PHONY: docker-logs
docker-logs: ## Tail Atlas Docker logs
	docker compose -f docker/docker-compose.yml logs -f atlas

# ── Setup ──────────────────────────────────────────────────────────────────
.PHONY: setup-dirs
setup-dirs: ## Create the /data directory structure
	@echo "$(YELLOW)Creating data directories...$(RESET)"
	@mkdir -p data/db data/uploads data/backups data/logs
	@echo "$(GREEN)Data directories ready.$(RESET)"

.PHONY: setup-env
setup-env: ## Copy .env.example to .env (safe — does not overwrite)
	@if [ -f .env ]; then \
		echo "$(YELLOW).env already exists. Edit it manually.$(RESET)"; \
	else \
		cp .env.example .env; \
		echo "$(GREEN).env created. Fill in ATLAS_SESSION_SECRET before running.$(RESET)"; \
	fi

.PHONY: setup
setup: setup-dirs setup-env ## Full first-time setup (dirs + .env)
	@echo "$(GREEN)Setup complete. Run: make migrate && make dev$(RESET)"

# ── Cleanup ────────────────────────────────────────────────────────────────
.PHONY: clean
clean: ## Remove build artefacts, coverage files, and Air temp dir
	@rm -rf $(BUILD_DIR) tmp/ coverage.out coverage.html
	@echo "$(GREEN)Clean.$(RESET)"