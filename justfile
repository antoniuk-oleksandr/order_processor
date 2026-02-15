alias cov := coverage
alias doc := documentation

# Run tests with coverage
coverage:
    go test -coverprofile=coverage.out ./internal/...
    go tool cover -html=coverage.out -o coverage.html

# Run all tests
test:
    ginkgo -r -race -v

# Run only end-to-end tests
e2e:
    ginkgo -r -v -race --label-filter="e2e"

# Run only unit tests
unit:
    ginkgo -r -v --label-filter="unit"

# Run the application
run:
    go run cmd/order_processor/main.go

# Generate documentation and open it in the browser
documentation:
    pkgsite -open .
