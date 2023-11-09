.PHONY: all clean test docs install setup release-patch release-minor release-major
module := github.com/duanemay/chatgpt-cli
all: chatgpt-cli

ifeq ($(OS),Windows_NT)
    uname_S := Windows
else
    uname_S := $(shell uname -s)
endif

chatgpt-cli: **/*.go go.mod go.sum
	goreleaser build --snapshot --clean --single-target --output .

**/*.go:
	@# noop

clean:
	rm -rf dist/ coverage.html .coverage-report.out chatgpt-cli

test:
	go run github.com/onsi/ginkgo/v2/ginkgo -r

docs:
	docs/generate-demos.sh

install:
	cp chatgpt-cli $(which chatgpt-cli)

race:
	go run github.com/onsi/ginkgo/v2/ginkgo -r -race --trace

cov:
	go run github.com/onsi/ginkgo/v2/ginkgo -r \
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

re-release-novalidate:
	goreleaser release --clean --skip-validate

release-patch:
	git tag "$(shell svu patch)"
	git push --tags
	goreleaser release --clean

release-minor:
	git tag "$(shell svu minor)"
	git push --tags
	goreleaser release --clean

release-major:
	git tag "$(shell svu major)"
	git push --tags
	goreleaser release --clean
