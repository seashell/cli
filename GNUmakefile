SHELL = sh
PROJECT_ROOT := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))

GIT_COMMIT := $(shell git rev-parse HEAD)
GIT_DIRTY := $(if $(shell git status --porcelain),+CHANGES)

GO_LDFLAGS = "-s -w -X=github.com/seashell/cli/version.GitCommit=$(GIT_COMMIT)$(GIT_DIRTY)"

INSTALLATION_PATH ?= /usr/local/bin

default: help
					
bin/seashell: $(SOURCE_FILES) ## Build seashell CLI binary
	@echo "==> Building the seashell CLI executable..."
	@GOOS=linux GOARCH=amd64 \
			go build \
			-trimpath \
	 		-ldflags $(GO_LDFLAGS) \
			-o "$@"


.PHONY: install
install: clean bin/seashell $(SOURCE_FILES) ## Build and install the seashell CLI binary on the system
	@echo "==> Installing CLI executable at $(INSTALLATION_PATH)/seashell..."
	@sudo cp $(PROJECT_ROOT)/bin/seashell $(INSTALLATION_PATH)/seashell				

.PHONY: clean
clean: ## Remove build artifacts
	@echo "==> Cleaning build artifacts from $(PROJECT_ROOT)/bin/ ..."
	@go mod tidy
	@rm -rf "$(PROJECT_ROOT)/bin/"

HELP_FORMAT="    \033[36m%-25s\033[0m %s\n"
EG_FORMAT="    \033[36m%s\033[0m %s\n"
.PHONY: help
help: ## Display this usage information
	@echo "Valid targets:"
	@grep -E '^[^ ]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		sort | \
		awk 'BEGIN {FS = ":.*?## "}; \
			{printf $(HELP_FORMAT), $$1, $$2}'
	@echo ""
	@echo "Examples:"
	@printf $(EG_FORMAT) "~${PWD}" "$$ make bin/seashell"
	@printf $(EG_FORMAT) "~${PWD}" "$$ make install"
