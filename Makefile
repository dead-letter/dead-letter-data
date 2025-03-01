PROTO_DIR=../dead-letter-manifests/protos
OUT_DIR=./pb

PROTOC=$(shell which protoc)
PROTOC_GEN_GO=$(shell which protoc-gen-go)
PROTOC_GEN_GO_GRPC=$(shell which protoc-gen-go-grpc)

.PHONY: all
all: data

.PHONY: check
check:
	@which protoc > /dev/null || { echo "❌ protoc not found"; exit 1; }
	@which protoc-gen-go > /dev/null || { echo "❌ protoc-gen-go not found"; exit 1; }
	@which protoc-gen-go-grpc > /dev/null || { echo "❌ protoc-gen-go-grpc not found"; exit 1; }

.PHONY: data
data: check
	@mkdir -p $(OUT_DIR)
	@$(PROTOC) -I $(PROTO_DIR) \
			--go_out=$(OUT_DIR) --go_opt=paths=source_relative \
	        --go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
	        $(PROTO_DIR)/data.proto
	@echo "✅ Generated gRPC code for data.proto"

.PHONY: clean
clean:
	@rm -f $(OUT_DIR)/*.pb.go
	@echo "🗑️  Cleaned generated files"
