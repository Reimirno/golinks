all: build

VERSION?=1.0.0
COMMIT?=$(shell git rev-parse --short HEAD)
DATE?=$(shell date +%Y-%m-%dT%H:%M:%SZ)

##### Commands for go server #####
GO_CMD=go

OUT_DIR = ./build
PROTO_DIR = ./pkg/pb
COVERAGE_DIR = ./coverage

MAIN_FILE=main.go
PROTO_FILE=service.proto

gen:
	protoc --version
	protoc --proto_path=$(PROTO_DIR) \
       --go_out=$(PROTO_DIR) \
       --go_opt=paths=source_relative \
       --go-grpc_out=$(PROTO_DIR) \
       --go-grpc_opt=paths=source_relative \
       $(PROTO_DIR)/$(PROTO_FILE)

build: windows linux mac

windows: gen
	$(call BUILD_GO_BINARY,windows,amd64,$(OUT_DIR)/golink.exe,$(MAIN_FILE))

linux: gen
	$(call BUILD_GO_BINARY,linux,amd64,$(OUT_DIR)/golink,$(MAIN_FILE))

mac: gen
	$(call BUILD_GO_BINARY,darwin,amd64,$(OUT_DIR)/golink_darwin,$(MAIN_FILE))

define BUILD_GO_BINARY
	GOOS=$(1) GOARCH=$(2) $(GO_CMD) build -ldflags="-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildDate=$(DATE)" -o $(3) $(4)
endef

test:
	$(GO_CMD) test ./...

cover:
	mkdir -p $(COVERAGE_DIR)
	$(GO_CMD) test -v -coverprofile=$(COVERAGE_DIR)/coverage.out -covermode=atomic ./...
	$(GO_CMD) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html

clean:
	rm -rf $(OUT_DIR) $(PROTO_DIR)/*.pb.go $(COVERAGE_DIR)

lint:
	golangci-lint run --config .golangci.yaml --sort-results


##### Commands for browser extension #####
NPM_CMD=npm

EXT_DIR=browser
CHROME_DIR=$(EXT_DIR)/chrome

test-chrome:
	cd $(CHROME_DIR) && $(NPM_CMD) test


##### Commands for web app #####
WEB_DIR=web

dev-web:
	cd $(WEB_DIR) && $(NPM_CMD) run dev

build-web:
	cd $(WEB_DIR) && $(NPM_CMD) run build

test-web:
	cd $(WEB_DIR) && $(NPM_CMD) run test
