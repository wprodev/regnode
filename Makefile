.PHONY: all clean

BINARY_NAME=regnode
BUILD_DIR=build
SRC_FILES=cmd/regnode/main.go

# Build for Linux amd64
build-linux:
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/linux-amd64/$(BINARY_NAME) $(SRC_FILES)

# Build for macOS amd64
build-macos:
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/darwin-amd64/$(BINARY_NAME) $(SRC_FILES)

# Build for macOS arm64 (Apple Silicon)
build-macos-arm:
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/darwin-arm64/$(BINARY_NAME) $(SRC_FILES)

all: build-linux build-macos build-macos-arm

clean:
	rm -rf $(BUILD_DIR)
