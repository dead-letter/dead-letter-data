PROTO_DIR=../dead-letter-manifests/protos
PB_OUT_DIR=./pb
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


.PHONY: pb/check
pb/check:
	@which protoc > /dev/null || { echo "❌ protoc not found"; exit 1; }
	@which protoc-gen-go > /dev/null || { echo "❌ protoc-gen-go not found"; exit 1; }
	@which protoc-gen-go-grpc > /dev/null || { echo "❌ protoc-gen-go-grpc not found"; exit 1; }


.PHONY: pb/data
pb/data: pb/check
	@mkdir -p $(PB_OUT_DIR)
	@$(PROTOC) -I $(PROTO_DIR) \
			--go_out=$(PB_OUT_DIR) --go_opt=paths=source_relative \
	        --go-grpc_out=$(PB_OUT_DIR) --go-grpc_opt=paths=source_relative \
	        $(PROTO_DIR)/data.proto
	@echo "✅ Generated gRPC code for data.proto"


.PHONY: pb/clean
pb/clean:
	@rm -f $(PB_OUT_DIR)/*.pb.go
	@echo "🗑️  Cleaned generated files"


## db/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	psql ${DATABASE_URL}


## db/migrations/new label=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo "Creating migration files for ${label}..."
	go tool goose create ${label} sql


# confirmation dialog helper
.PHONY: confirm
confirm:
	@echo -n "Are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]
