EXECUTABLE_NAME := chatgpt-cli
MODULE := github.com/duanemay/chatgpt-cli

ifeq ($(GOHOSTOS),)
	GOHOSTOS:=$(shell uname | tr A-Z a-z | sed 's/mingw/windows/; s/.*windows.*/windows/')
endif

.PHONY: all
all: ${EXECUTABLE_NAME}

${EXECUTABLE_NAME}: **/*.go go.mod go.sum
	goreleaser build --snapshot --clean --single-target --output .

**/*.go:
	@# noop

.PHONY: clean
clean: ## Clean up artifacts
	rm -rf dist/ coverage.html .coverage-report.out chatgpt-cli

.PHONY: test
test: ## Run tests
	go run github.com/onsi/ginkgo/v2/ginkgo -r

.PHONY: docs
docs: ## Generate still/animated images used for documentation
	docs/generate-demos.sh

.PHONY: install
install: ${EXECUTABLE_NAME}  ## Build and Install the binary
	cp chatgpt-cli $(which chatgpt-cli)

.PHONY: lint
lint:  ## Lint the plugin
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run -v

.PHONY: race
race: ## Detect race conditions during test
	go run github.com/onsi/ginkgo/v2/ginkgo -r --race --trace

.PHONY: cov
cov: ## Generate test coverage report
	go run github.com/onsi/ginkgo/v2/ginkgo -r \
 		--coverprofile=.coverage-report.out
	go tool cover -html=./.coverage-report.out -o coverage.html
ifeq ($(GOHOSTOS),darwin)
	open coverage.html
else
	@echo To view coverage, open: coverage.html
endif

.PHONY: update
update: ## Update Dependencies
	go get -t -u ./...
	go mod tidy
	go run github.com/onsi/ginkgo/v2/ginkgo -r

.PHONY: setup
setup:  ## Setup packages needed for release
	brew install caarlos0/tap/svu
	brew install goreleaser/tap/goreleaser

.PHONY: re-release-novalidate
re-release-novalidate:  ## Recreate a release with current tag
	goreleaser release --clean --skip=validate

.PHONY: release-patch
release-patch: ## Create a new patch release
	git tag "$(shell svu patch)"
	git push --tags
	goreleaser release --clean

.PHONY: release-minor
release-minor: ## Create a new minor release
	git tag "$(shell svu minor)"
	git push --tags
	goreleaser release --clean

.PHONY: release-major
release-major: ## Create a new major release
	git tag "$(shell svu major)"
	git push --tags
	goreleaser release --clean

help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
