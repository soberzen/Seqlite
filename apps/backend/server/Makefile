APP_ENV ?= development
export APP_ENV
APP_PACKAGE := ./cmd
OUTPUT := ./dist/server
COMPOSE_FLAGS = --env-file ".env.${APP_ENV}" -p "seqlite-${APP_ENV}"

build:
	@echo "📦 Building $(OUTPUT)..."
	@mkdir -p $(dir $(OUTPUT))
	@go build -o $(OUTPUT) $(APP_PACKAGE)
	@echo "✅ $(OUTPUT) build successfully"

dev:
	@echo "🚀 Running in development mode..."
	@APP_ENV=development go run $(APP_PACKAGE)
	@echo "✅ $(OUTPUT) dev successfully"

prod: build
	@echo "🚀 Running in production mode..."
	@APP_ENV=production ./$(OUTPUT)
	@echo "✅ $(OUTPUT) prod successfully"

docker-run:
	@if docker compose version >/dev/null 2>&1; then \
		docker compose -f docker-compose.yml $(COMPOSE_FLAGS) up --build; \
	elif command -v docker-compose >/dev/null 2>&1; then \
		echo "Using Docker Compose V1"; \
		docker-compose -f docker-compose.yml $(COMPOSE_FLAGS) up --build; \
	else \
		echo "error: please ensure Docker Compose is installed"; \
		exit 1; \
	fi

docker-down:
	@if docker compose version >/dev/null 2>&1; then \
		docker compose -f docker-compose.yml $(COMPOSE_FLAGS) down; \
	elif command -v docker-compose >/dev/null 2>&1; then \
		echo "Using Docker Compose V1"; \
		docker-compose -f docker-compose.yml $(COMPOSE_FLAGS) down; \
	else \
		echo "error: please ensure Docker Compose is installed"; \
		exit 1; \
	fi


migrate:
	@goose -env ".env.$(APP_ENV)" -dir ./migrations up

migrate-status:
	@goose -env ".env.$(APP_ENV)" -dir ./migrations status

migrate-down:
	@goose -env ".env.$(APP_ENV)" -dir ./migrations down

clean:
	@echo "Cleaning..."
	@rm -rf dist
