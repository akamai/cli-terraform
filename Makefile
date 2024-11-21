BIN = $(CURDIR)/bin
GOCMD = go
GOTEST = $(GOCMD) test
GOBUILD = $(GOCMD) build
M = $(shell echo ">")

# Until v0.25.0 is not fixed, we have to use previous version. To install it, we must enable module aware mode.
GOIMPORTS = $(BIN)/goimports
GOIMPORTS_VERSION = v0.24.0
# Rule to install goimports with version pinning
$(GOIMPORTS): | $(BIN) ; $(info $(M) Installing goimports $(GOIMPORTS_VERSION)...)
	$Q env GO111MODULE=on GOBIN=$(BIN) $(GOCMD) install golang.org/x/tools/cmd/goimports@$(GOIMPORTS_VERSION)

$(BIN):
	@mkdir -p $@
$(BIN)/%: | $(BIN) ; $(info $(M) Installing $(PACKAGE)...)
	env GOBIN=$(BIN) $(GOCMD) install $(PACKAGE)

GOCOV = $(BIN)/gocov
$(BIN)/gocov: PACKAGE=github.com/axw/gocov/gocov@v1.1.0

GOCOVXML = $(BIN)/gocov-xml
$(BIN)/gocov-xml: PACKAGE=github.com/AlekSi/gocov-xml@v1.1.0

GOJUNITREPORT = $(BIN)/go-junit-report
$(BIN)/go-junit-report: PACKAGE=github.com/jstemmer/go-junit-report/v2@v2.1.0

TFLINT = $(BIN)/tflint
$(BIN)/tflint: $(BIN) ; $(info $(M) Installing tflint...)
	@export TFLINT_INSTALL_PATH=$(BIN); \
	curl -sSfL https://raw.githubusercontent.com/terraform-linters/tflint/master/install_linux.sh  | bash

GOLANGCILINT = $(BIN)/golangci-lint
GOLANGCI_LINT_VERSION = v1.55.2
$(BIN)/golangci-lint: ; $(info $(M) Installing golangci-lint...)
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(BIN) $(GOLANGCI_LINT_VERSION)

COVERAGE_MODE = atomic
COVERAGE_DIR = $(CURDIR)/test/coverage
COVERAGE_PROFILE = $(COVERAGE_DIR)/profile.out
COVERAGE_XML = $(COVERAGE_DIR)/coverage.xml
COVERAGE_HTML = $(COVERAGE_DIR)/index.html

.PHONY: all
all: clean tidy fmt-check lint terraform-fmt terraform-lint coverage create-junit-report create-coverage-files clean-tools

.PHONY: test
test: ; $(info $(M) Running tests...) ## Run all unit tests
	$(GOTEST) -v -count=1 ./...

.PHONY: coverage
coverage: ; $(info $(M) Running tests with coverage...) @ ## Run tests and generate coverage profile
	@mkdir -p $(COVERAGE_DIR)
	@$(GOTEST) -v -covermode=$(COVERAGE_MODE) \
               -coverprofile="$(COVERAGE_PROFILE)" ./... | tee test/tests.output

.PHONY: create-junit-report
create-junit-report: | $(GOJUNITREPORT) ; $(info $(M) Creating juint xml report) @ ## Generate junit-style coverage report
	@cat $(CURDIR)/test/tests.output | $(GOJUNITREPORT) > $(CURDIR)/test/tests.xml
	@sed -i -e 's/skip=/skipped=/g' $(CURDIR)/test/tests.xml

.PHONY: create-coverage-files
create-coverage-files: | $(GOCOV) $(GOCOVXML); $(info $(M) Creating coverage files...) @ ## Generate coverage report files
	@$(GOCMD) tool cover -html=$(COVERAGE_PROFILE) -o $(COVERAGE_HTML)
	@$(GOCOV) convert $(COVERAGE_PROFILE) | $(GOCOVXML) > $(COVERAGE_XML)

.PHONY: tidy
tidy: ; $(info $(M) Running go mod tidy...) @
	@$(GOMODTIDY)

.PHONY: lint
lint: | $(GOLANGCILINT); $(info $(M) Running linter...) @ ## Run golangci-lint on all source files
	@$(BIN)/golangci-lint run

.PHONY: fmt
fmt:  | $(GOIMPORTS); $(info $(M) Running goimports...) @ ## Run goimports on all source files
	$Q $(GOIMPORTS) -w .

.PHONY: fmt-check
fmt-check: | $(GOIMPORTS); $(info $(M) Running format and imports check...) @ ## Run goimports on all source files (do not modify files)
	$(eval OUTPUT = $(shell $(GOIMPORTS) -l .))
	@if [ "$(OUTPUT)" != "" ]; then\
		echo "Found following files with incorrect format and/or imports:";\
		echo "$(OUTPUT)";\
		false;\
	fi

.PHONY: terraform-fmt
terraform-fmt: ; $(info $(M) Running terraform fmt check...) @ ## Check formatting of all HCL files in the project
	terraform fmt -recursive -check

.PHONY: terraform-lint
terraform-lint: | $(TFLINT) ; $(info $(M) Checking source code against tflint...) @ ## Run tflint on all HCL files in the project
	@find ./pkg -type f -name "*.tf" | xargs -I % dirname % | sort -u | xargs -I @ sh -c "echo @ && $(TFLINT) --filter @"

.PHONY: validate-testdata
validate-testdata: ; $(info $(M) Validating testdata agains terraform-provider-akamai...) @ ## terraform init & terraform validate on all .tf files
	# look up all unique directories containing .tf files and execute 'terraform validate'
	@for dir in $(shell find . -type f -name "*.tf" -exec dirname "{}" \; |sort -u); do \
		pushd $${dir} && \
		echo Validating directory: $(shell pwd) && \
		terraform init -upgrade -no-color && \
		terraform validate -no-color && \
		rm -r .terraform* ; \
		popd ; \
	done

.PHONY: release
release: ; $(info $(M) Generating release binaries and signatures...) @ ## Generate release binaries
	@./scripts/build.sh

.PHONY: clean
clean: ; $(info $(M) Removing 'tools' directory and test results...) @ ## Cleanup installed packages and test reports
	@rm -rf $(BIN)
	@rm -rf $(BIN)/test/tests.* $(BIN)/test/coverage

clean-tools: ## Cleanup installed packages
	@rm -rf $(BIN)/go*

.PHONY: help
help: ## List all make targets
	echo $(MAKEFILE_LIST)
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}'
