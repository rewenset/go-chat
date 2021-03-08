APP=go-chat
BUILD_DIR=bin

.PHONY: build
build: clean
	@echo "Building..."
	go build -o ${BUILD_DIR}/${APP} .

.PHONY: run
run:
	@echo "Starting..."
	go run .

.PHONY: clean
clean:
	@echo "Cleaning..."
	rm -rf ${BUILD_DIR}
