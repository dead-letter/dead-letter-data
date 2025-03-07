PB_DIR_IN=../dead-letter-manifests/protos
PB_DIR_OUT=./pkg/pb
PROTOC=$(shell which protoc)

## help: print this help message
.PHONY: help
help:
	@echo "Usage:"
	@sed -n "s/^##//p" ${MAKEFILE_LIST} | column -t -s ":" |  sed -e "s/^/ /"


## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit:
	@echo "Tidying and verifying module dependencies..."
	go mod tidy
	go mod verify
	@echo "Formatting code..."
	go fmt ./...
	@echo "Vetting code..."
	go vet ./...
	go tool staticcheck ./...
	@echo "Running tests..."
	go test -race -vet=off ./...


# test/check: ensure pgtestdb container is running
.PHONY: test/check
test/check:
	@docker info > /dev/null 2>&1 || (echo "docker is not running"; exit 1)
	@docker compose ps -q pgtestdb | grep -q . || (echo -e "pgtestdb container is not runnning\nrun: docker compose up -d"; exit 1)

## test: run verbose tests
.PHONY: test
test: test/check
	go test -v ./...


# proto/check: check for necessary build tools and directories
.PHONY: proto/check
proto/check:
	@which protoc > /dev/null || { echo "❌ protoc not found"; exit 1; }
	@which protoc-gen-go > /dev/null || { echo "❌ protoc-gen-go not found"; exit 1; }
	@which protoc-gen-go-grpc > /dev/null || { echo "❌ protoc-gen-go-grpc not found"; exit 1; }
	@[ -d "$(PB_DIR_IN)" ] || { echo "❌ PB_DIR_IN $(PB_DIR_IN) does not exist"; exit 1; }


## proto/gen: generate protoc stubs
.PHONY: proto/gen
proto/gen: proto/check
	@mkdir -p $(PB_DIR_OUT)
	@$(PROTOC) -I $(PB_DIR_IN) \
			--go_out=$(PB_DIR_OUT) --go_opt=paths=source_relative \
	        --go-grpc_out=$(PB_DIR_OUT) --go-grpc_opt=paths=source_relative \
	        $(PB_DIR_IN)/data.proto
	@echo "✅ Generated gRPC code for data.proto"


## proto/clean: clean generated files
.PHONY: proto/clean
proto/clean:
	@rm -f $(PB_DIR_OUT)/*.pb.go
	@echo "🗑️  Cleaned generated files"


## migrations/new label=$1: create a new database migration
.PHONY: migrations/new
migrations/new:
	@echo "Creating migration files for ${label}..."
	go tool goose -dir migrations -s create ${label} sql


# confirmation dialog helper
.PHONY: confirm
confirm:
	@echo -n "Are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]
