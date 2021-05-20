# Set Shell to bash, otherwise some targets fail with dash/zsh etc.
SHELL := /bin/bash

# Disable built-in rules
MAKEFLAGS += --no-builtin-rules
MAKEFLAGS += --no-builtin-variables
.SUFFIXES:
.SECONDARY:

include Makefile.vars.mk

.PHONY: help
help: ## Show this help
	@grep -E -h '\s##\s' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: test
test: ## Run the unit tests
	go test -v -coverprofile cover.out ./...

integration-test: export RESTIC_PATH = $(restic_path)
integration-test: export RESTIC_BINARY = $(restic_path)
integration-test: export RESTIC_PASSWORD = $(restic_password)
integration-test: export RESTIC_REPOSITORY = s3:http://$(minio_address)/test
integration-test: export RESTORE_S3ENDPOINT = http://$(minio_address)/restore
integration-test: export AWS_ACCESS_KEY_ID = $(minio_root_user)
integration-test: export AWS_SECRET_ACCESS_KEY = $(minio_root_password)
integration-test: export RESTORE_ACCESSKEYID = $(minio_root_user)
integration-test: export RESTORE_SECRETACCESSKEY = $(minio_root_password)
integration-test: export BACKUP_DIR = $(backup_dir)
integration-test: export RESTORE_DIR = $(restore_dir)
integration-test: export STATS_URL = $(stats_url)
integration-test: integration-test-setup ## Run the integration test
	go test -v -tags=integration -coverprofile cover.out ./... ./cmd/wrestic/...

.PHONY: integration-test-setup
integration-test-setup: minio-start restic-download  ## Prepare to run the integration test

.PHONY: clean
clean: minio-stop
	rm -rf $(test_dir)
	rm -f wrestic

.PHONY: minio-address
minio-address: ## Get the address to connect to minio
	@echo "http://$(minio_address)"

.PHONY: minio-reset
minio-reset: minio-stop minio-delete-config minio-start ## Reset minio's configuration and data dirs and restart minio

minio-delete-config:
	rm -rf "$(minio_config)" "$(minio_data)"

.PHONY: minio-restart
minio-restart: minio-stop minio-start ## Restart minio

minio-set-alias: minio-start ## Set the alias 'wrestic' in mc to the minio server
	@mc alias set wrestic "http://$(minio_address)" "$(minio_root_user)" "$(minio_root_password)"

.PHONY: minio-start
minio-start: minio-check $(minio_pid) ## Run minio

.PHONY: minio-check
minio-check: minio-clean ## Check if minio is running
	@test -f "$(minio_pid)" && echo "Minio runs as PID $$(cat $(minio_pid))." || echo "Minio is not running."

.PHONY: minio-clean
minio-clean: ## Remove the minio PID file if the process is not running
	@./clean.sh "$(minio_pid)"

.PHONY: minio-stop
minio-stop: ## Stop minio
	@./kill.sh "$(minio_pid)"

minio-download: .test/minio ## Download minio

restic-download: .test/restic ## Download restic

$(minio_pid): export MINIO_ACCESS_KEY = $(minio_root_user)
$(minio_pid): export MINIO_SECRET_KEY = $(minio_root_password)
$(minio_pid): minio-download
	@mkdir -p "$(minio_data)" "$(minio_config)"
	@./exec.sh "$(minio_pid)" \
		"$(minio_path)" \
			server "$(minio_data)" \
			"--address" "$(minio_address)" \
			"--config-dir" "$(minio_config)"
	@while ! curl --silent "http://$(minio_address)" > /dev/null; do echo "Waiting for server http://$(minio_address) to become ready"; sleep 0.5; done

.test/minio:
	@mkdir -p $(test_dir)
	curl $(curl_args) --output "$(minio_path)" "$(minio_url)"
	chmod +x "$(minio_path)"
	"$(minio_path)" --version

.test/restic:
	@mkdir -p $(test_dir)
	curl $(curl_args) "$(restic_url)" | \
		bunzip2 > "$(restic_path)"
	chmod +x "$(restic_path)"
	"$(restic_path)" version
