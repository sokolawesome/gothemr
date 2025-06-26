build:
    go build -o gothemr ./cmd/gothemr

test:
    go test ./...

fmt:
    go fmt ./...

vet:
    go vet ./...

dev: fmt vet test build

install:
    go install ./cmd/gothemr

clean:
    rm -f gothemr