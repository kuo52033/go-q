APP_SERVER_SRC = cmd/server/server.go
APP_WORKER_SRC = cmd/worker/worker.go

BIN_DIR = ./bin
TARGET_SERVER=$(BIN_DIR)/go-q-server
TARGET_WORKER=$(BIN_DIR)/go-q-worker

# tell system that these are not files, instead they are commands
.PHONY: all build run start clean dev run-server run-worker build-server build-worker

all: build

build: build-server build-worker

# @ is used to suppress the output of the command

build-server:
	@echo "Building server..."
	@go build -o $(TARGET_SERVER) $(APP_SERVER_SRC)
	@codesign --sign "go-local" --force --deep $(TARGET_SERVER)

build-worker:
	@echo "Building worker..."
	@go build -o $(TARGET_WORKER) $(APP_WORKER_SRC)
	@codesign --sign "go-local" --force --deep $(TARGET_WORKER)

dev: build
	@pkill -f $(TARGET_SERVER) || true
	@pkill -f $(TARGET_WORKER) || true
	@echo "Running server and worker..."
	@$(TARGET_SERVER) &
	@$(TARGET_WORKER) 

run-server: build-server
	@echo "Running server..."
	@$(TARGET_SERVER)

run-worker: build-worker
	@echo "Running worker..."
	@$(TARGET_WORKER)

clean:
	@echo "Cleaning up..."
	@rm -rf $(BIN_DIR)
