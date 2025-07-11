# List available tasks
list:
    @just --list

# Build the application binary
build:
    @echo "Building gothemr..."
    @go build -o gothemr ./cmd/gothemr

# Run the application (builds it first)
run: build
    @echo "Running gothemr..."
    @./gothemr

# Format all Go source files
fmt:
    @echo "Formatting code..."
    @go fmt ./...

# Tidy the mod file
tidy:
    @echo "Tidying dependencies..."
    @go mod tidy

# Lint the codebase
lint:
    @echo "Linting code..."
    @golangci-lint run ./...

# Vet the codebase
vet:
    @echo "Checking packages..."
    @go vet ./...

# Test the codebase
test:
    @echo "Testing code..."
    @go test ./...

# Run all checkers and build
dev: lint fmt vet test build

# Clean binaries
clean:
    rm -f gothemr
