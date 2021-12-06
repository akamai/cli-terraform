# Tools versions
golangci-lint-version = v1.41.1

.PHONY: all
all: fmt lint vet coverage

.PHONY: test
test:
	go test -count=1 -race ./...

.PHONY: coverage-ui
coverage-ui:
	go test -covermode=count -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

.PHONY: coverage
coverage:
	go test -coverprofile coverage.out ./...
	go tool cover -func coverage.out | grep total

.PHONY: lint
lint:
	@echo "==> Checking source code against golangci-lint"
	@$$(go env GOPATH)/bin/golangci-lint run

.PHONY: vet
vet:
	go vet ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: init
init:
	@echo Installing golangci-lint
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin $(golangci-lint-version)

.PHONY: validate-testdata
validate-testdata:
	# look up all unique directories containing .tf files and execute 'terraform validate'
	@for dir in $(shell find . -type f -name "*.tf" -exec dirname "{}" \; |sort -u); do \
		pushd $${dir} && \
		echo Validating directory: $(shell pwd) && \
		terraform init && \
		terraform validate && \
		rm -r .terraform* ; \
		popd ; \
	done