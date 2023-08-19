.PHONY: all clean test setup release-patch release-minor release-major
module := github.com/duanemay/chatgpt-cli
all: chatgpt-cli
ifeq ($(OS),Windows_NT)
    uname_S := Windows
else
    uname_S := $(shell uname -s)
endif

chatgpt-cli: **/*.go
	goreleaser build --snapshot --clean --single-target --output .

**/*.go:
	@# noop

clean:
	rm -rf dist/ coverage.html .coverage-report.out chatgpt-cli

test:
	ginkgo -r -v

race:
	ginkgo -r -v -race --trace

cov:
	ginkgo -r \
		--coverpkg=${module},${module}/cmd \
 		--coverprofile=.coverage-report.out
	go tool cover -html=./.coverage-report.out -o coverage.html
ifeq ($(uname_S),Darwin)
	open coverage.html
else
	@echo To view coverage: open coverage.html
endif

setup:
	brew install caarlos0/tap/svu
	brew install goreleaser/tap/goreleaser

release-patch:
	git tag "$(svu patch)"
	git push --tags
	goreleaser release --clean

release-minor:
	git tag "$(svu minor)"
	git push --tags
	goreleaser release --clean

release-major:
	git tag "$(svu major)"
	git push --tags
	goreleaser release --clean
