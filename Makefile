# ====================================================================================
# Project Variables
# ====================================================================================

# --- General ---
APP_NAME := card-manage
GO_MAIN := ./cmd/server/main.go

# --- Local Environment ---
# This DB_URL is used by the migrate tool for the local Docker container.
# It matches the settings in docker-compose.yml.
LOCAL_DB_URL := "postgres://root:secret@localhost:5432/card_manage?sslmode=disable"

# --- GCP Environment ---
# !!! IMPORTANT !!!
# Please replace these placeholder values with your actual GCP project details.
GCP_PROJECT_ID := "your-gcp-project-id"
GCP_REGION := "your-gcp-region" # e.g., asia-east1
GCP_ARTIFACT_REPO := "$(GCP_REGION)-docker.pkg.dev/$(GCP_PROJECT_ID)/$(APP_NAME)"
GCP_SERVICE_ACCOUNT := "your-service-account@$(GCP_PROJECT_ID).iam.gserviceaccount.com"

# This is the full identifier for the Cloud SQL instance.
GCP_CLOUD_SQL_INSTANCE := "$(GCP_PROJECT_ID):$(GCP_REGION):your-cloud-sql-instance-name"

# This DB_URL is for connecting to GCP Cloud SQL.
# It requires you to provide the password via an environment variable `GCP_DB_PASSWORD`.
# The Cloud SQL Auth Proxy must be running for this to work.
GCP_DB_URL := "postgres://your-db-user:$(GCP_DB_PASSWORD)@localhost:5433/your-db-name?sslmode=disable"


# ====================================================================================
# Helper Targets
# ====================================================================================

.PHONY: help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "------------------[ Local Development ]------------------"
	@echo "  db-up                Start local PostgreSQL container using Docker Compose."
	@echo "  db-down              Stop and remove local PostgreSQL container."
	@echo "  migrate-up           Apply all up migrations to the local database."
	@echo "  migrate-down         Apply all down migrations to the local database."
	@echo "  run                  Run the Go application locally."
	@echo ""
	@echo "------------------[ GCP Deployment ]------------------"
	@echo "  gcp-auth             Configure gcloud Docker credential helper."
	@echo "  gcp-build-push       Build the Docker image and push it to Google Artifact Registry."
	@echo "  gcp-deploy           Deploy the application to Google Cloud Run."
	@echo "  cloud-sql-proxy      Start the Cloud SQL Auth Proxy (requires GCP_CLOUD_SQL_INSTANCE to be set)."
	@echo "  gcp-migrate-up       Apply all up migrations to the GCP Cloud SQL database (requires proxy to be running)."
	@echo "  gcp-migrate-down     Apply all down migrations to the GCP Cloud SQL database (requires proxy to be running)."


# ====================================================================================
# Local Development Targets
# ====================================================================================

.PHONY: db-up
db-up:
	@echo "Starting local PostgreSQL database..."
	@docker-compose up -d

.PHONY: db-down
db-down:
	@echo "Stopping local PostgreSQL database..."
	@docker-compose down

.PHONY: migrate-up
migrate-up:
	@echo "Applying database migrations to local DB..."
	@migrate -path db/migration -database "$(LOCAL_DB_URL)" -verbose up

.PHONY: migrate-down
migrate-down:
	@echo "Reverting database migrations from local DB..."
	@migrate -path db/migration -database "$(LOCAL_DB_URL)" -verbose down

.PHONY: run
run:
	@echo "Running the application..."
	@go run $(GO_MAIN)


# ====================================================================================
# GCP Deployment Targets
# ====================================================================================

.PHONY: gcp-auth
gcp-auth:
	@echo "Configuring Docker to authenticate with GCP Artifact Registry..."
	@gcloud auth configure-docker $(GCP_REGION)-docker.pkg.dev

.PHONY: gcp-build-push
gcp-build-push:
	@echo "Building and pushing Docker image to $(GCP_ARTIFACT_REPO)/$(APP_NAME):latest..."
	@docker build -t $(GCP_ARTIFACT_REPO)/$(APP_NAME):latest .
	@docker push $(GCP_ARTIFACT_REPO)/$(APP_NAME):latest

.PHONY: gcp-deploy
gcp-deploy:
	@echo "Deploying to Cloud Run in $(GCP_REGION)..."
	@gcloud run deploy $(APP_NAME) \
		--image=$(GCP_ARTIFACT_REPO)/$(APP_NAME):latest \
		--platform=managed \
		--region=$(GCP_REGION) \
		--service-account=$(GCP_SERVICE_ACCOUNT) \
		--add-cloudsql-instances=$(GCP_CLOUD_SQL_INSTANCE) \
		--allow-unauthenticated

.PHONY: cloud-sql-proxy
cloud-sql-proxy:
	@echo "Starting Cloud SQL Auth Proxy..."
	@echo "Instance: $(GCP_CLOUD_SQL_INSTANCE)"
	@echo "Use Ctrl+C to exit."
	@cloud_sql_proxy -instances=$(GCP_CLOUD_SQL_INSTANCE)=tcp:5433

.PHONY: gcp-migrate-up
gcp-migrate-up:
	@echo "Applying database migrations to GCP Cloud SQL..."
	@if [ -z "$(GCP_DB_PASSWORD)" ]; then \
		echo "Error: GCP_DB_PASSWORD environment variable is not set."; \
		echo "Please run: export GCP_DB_PASSWORD='your-password'"; \
		exit 1; \
	fi
	@migrate -path db/migration -database "$(GCP_DB_URL)" -verbose up

.PHONY: gcp-migrate-down
gcp-migrate-down:
	@echo "Reverting database migrations from GCP Cloud SQL..."
	@if [ -z "$(GCP_DB_PASSWORD)" ]; then \
		echo "Error: GCP_DB_PASSWORD environment variable is not set."; \
		echo "Please run: export GCP_DB_PASSWORD='your-password'"; \
		exit 1; \
	fi
	@migrate -path db/migration -database "$(GCP_DB_URL)" -verbose down
