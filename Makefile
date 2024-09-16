GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_TEST=$(GO_CMD) test

VERSION?=1.0.0
COMMIT?=$(shell git rev-parse --short HEAD)
DATE?=$(shell date +%Y-%m-%dT%H:%M:%SZ)

OUT_DIR = ./build
PROTO_DIR = ./pkg/pb

MAIN_FILE=main.go
PROTO_FILE=service.proto

all: build

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
	$(call BUILD_BINARY,windows,amd64,$(OUT_DIR)/golink.exe,$(MAIN_FILE))

linux: gen
	$(call BUILD_BINARY,linux,amd64,$(OUT_DIR)/golink,$(MAIN_FILE))

mac: gen
	$(call BUILD_BINARY,darwin,amd64,$(OUT_DIR)/golink_darwin,$(MAIN_FILE))

test: gen
	$(GO_TEST) ./...

clean:
	rm -rf $(OUT_DIR) $(PROTO_DIR)/*.pb.go

define BUILD_BINARY
	GOOS=$(1) GOARCH=$(2) $(GO_BUILD) -ldflags="-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildDate=$(DATE)" -o $(3) $(4)
endef
