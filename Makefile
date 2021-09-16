SHELL = bash

GOCMD = go
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test -v $(shell $(GOCMD) list ./... | grep -v /vendor/)
GOFMT = go fmt

APP := $(shell basename $(PWD) | tr '[:upper:]' '[:lower:]')
DATE := $(shell date -u +%Y-%m-%d%Z%H:%M:%S)
VERSION := $(shell git describe --tags 2>/dev/null || echo v0.0.1)
COVERAGE_DIR = coverage

CIRCLE_BUILD_NUM ?= 1
CIRCLE_SHA1 ?= $(shell git rev-parse HEAD)

BUILD_NUMBER ?= $(CIRCLE_BUILD_NUM)
BUILD_VERSION := $(VERSION)-$(BUILD_NUMBER)
DEPLOY_VERSION := $(shell echo $${GITHUB_REF/refs\/tags\//})
GIT_COMMIT_HASH ?= $(CIRCLE_SHA1)
JOB_ID ?= local

DOCKER_REPO ?= gomicro
DOCKER_IMAGE_NAME ?= $(APP)
DOCKER_IMAGE_LABEL ?= latest

DOCKER_LOCAL_SERVICE := $(shell if [[ -z "$$NO_LOCAL_SERVICE" ]]; then echo "-p 4567:4567"; fi)
DOCKER_SERVICE_NAME := $(APP)-app-$(JOB_ID)


.PHONY: all
all: test

.PHONY: clean
clean: ## Clean out all generated items
	-@$(GOCLEAN)
	-@rm -f output.txt
	-@rm -rf coverage
	-@docker rm -f $(shell docker ps -a -q -f name=$(APP)-app)

.PHONY: coverage
coverage:  ## Generates the code coverage from all the tests
	@echo "Total Coverage: $$(make coverage_compfriendly)%"

.PHONY: coverage_compfriendly
coverage_compfriendly:  ## Generates the code coverage in a computer friendly manner
	-@rm -rf coverage
	-@mkdir -p $(COVERAGE_DIR)/tmp
	@for j in $$(go list ./... | grep -v '/vendor/' | grep -v '/ext/'); do go test -covermode=count -coverprofile=$(COVERAGE_DIR)/$$(basename $$j).out $$j > /dev/null 2>&1; done
	@echo 'mode: count' > $(COVERAGE_DIR)/tmp/full.out
	@tail -q -n +2 $(COVERAGE_DIR)/*.out >> $(COVERAGE_DIR)/tmp/full.out
	@$(GOCMD) tool cover -func=$(COVERAGE_DIR)/tmp/full.out | tail -n 1 | sed -e 's/^.*statements)[[:space:]]*//' -e 's/%//'

.PHONY: deploy
deploy: ## Deploy the artifacts
	@echo "Logging into Docker Hub"
	-@echo "$(DOCKER_PASSWORD)" | docker login -u "$(DOCKER_USERNAME)" --password-stdin
	@VERSION=$(DEPLOY_VERSION) goreleaser release

.PHONY: help
help: ## Show This Help
	@for line in $$(cat Makefile | grep "##" | grep -v "grep" | sed  "s/:.*##/:/g" | sed "s/\ /!/g"); do verb=$$(echo $$line | cut -d ":" -f 1); desc=$$(echo $$line | cut -d ":" -f 2 | sed "s/!/\ /g"); printf "%-30s--%s\n" "$$verb" "$$desc"; done

.PHONY: test
test: unit_test ## Run all available tests

.PHONY: unit_test
unit_test: ## Run all available unit tests
	$(GOTEST)

.PHONY: fmt
fmt: ## Run gofmt
	@echo "checking formatting..."
	@$(GOFMT) $(shell $(GOLIST) ./... | grep -v '/vendor/')

.PHONY: vet
vet: ## Run go vet
	@echo "vetting..."
	@$(GOVET) $(shell $(GOLIST) ./... | grep -v '/vendor/')

.PHONY: lint
lint: ## Run golint
	@echo "linting..."
	@$(GOLINTCMD) -set_exit_status $(shell $(GOLIST) ./... | grep -v '/vendor/')
