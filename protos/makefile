MODULE = github.com/over-eng/monzopanel/protos

PROTO_DIR = proto
GO_OUT_DIR = .

PROTO_FILES = $(shell find $(PROTO_DIR) -name "*.proto")

.PHONY: install-tools
install-tools:
	go get google.golang.org/protobuf/cmd/protoc-gen-go
	go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

.PHONY: generate
generate:
	@echo "Generating Go code from proto files..."
	@for file in $(PROTO_FILES); do \
		protoc \
			--proto_path=$(PROTO_DIR) \
			--go_out=$(GO_OUT_DIR) \
			--go_opt=paths=source_relative \
			--go-grpc_out=$(GO_OUT_DIR) \
			--go-grpc_opt=paths=source_relative \
			$$file; \
	done
