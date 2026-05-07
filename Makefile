# --- Project Variables ---
BINARY_NAME=backend_api
MAIN_PATH=./cmd/api
BUILD_DIR=./build

# --- Server Variables (Update these with your Putty info) ---
SERVER_USER=ubuntu
SERVER_IP=123.45.67.89
SERVER_PATH=/home/ubuntu/app

# --- 1. Local Development ---
.PHONY: tidy build clean

release: tidy deploy

tidy:
	go mod tidy

clean:
	@echo "Cleaning old builds..."
	@if [ -d "$(BUILD_DIR)" ]; then rm -rf $(BUILD_DIR); fi
	@if [ -f "$(BINARY_NAME)" ]; then rm $(BINARY_NAME); fi

# --- 2. Production Build ---
# We use GOOS=linux because your server is likely Linux, even if you are on Windows.
build: clean
	@echo "Building static binary for Linux..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@cp config.docker.yaml $(BUILD_DIR)/config.yaml
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# --- 3. Deployment ---
# This target handles the "Filezilla" and "Putty" steps in one go.
deploy: build
	@echo "Step 1: Creating remote directory if it doesn't exist..."
	ssh $(SERVER_USER)@$(SERVER_IP) "mkdir -p $(SERVER_PATH)"

	@echo "Step 2: Uploading binary and config (Replacing Filezilla)..."
	scp -r $(BUILD_DIR)/* $(SERVER_USER)@$(SERVER_IP):$(SERVER_PATH)

	@echo "Step 3: Restarting service (Replacing Putty manual commands)..."
	# Option A: If using Systemd
	# ssh $(SERVER_USER)@$(SERVER_IP) "sudo systemctl restart my-api-service"
	# Option B: If just running the binary in the background (simple version)
	ssh $(SERVER_USER)@$(SERVER_IP) "cd $(SERVER_PATH) && pkill $(BINARY_NAME) || true && nohup ./$(BINARY_NAME) > log.txt 2>&1 &"
	
	@echo "Deployment Successful!"