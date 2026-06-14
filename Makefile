DEPLOY_HOST ?= mango
DEPLOY_BIN ?= ./target/umed_linux_amd64
DEPLOY_SERVICE ?= deploy/systemd/ume.service
ENV_FILE ?=
DEFAULT_ENV_FILE ?= deploy/ume.env
FALLBACK_ENV_FILE ?= app.env
REMOTE_BIN ?= /usr/local/bin/ume
REMOTE_SERVICE ?= /etc/systemd/system/ume.service
REMOTE_STATE_DIR ?= /var/lib/ume
REMOTE_TMP_DIR ?= /var/umed/tmp
REMOTE_ENV_DIR ?= /etc/ume
REMOTE_ENV ?= /etc/ume/ume.env
TMP_BIN ?= /tmp/umed_linux_amd64
TMP_SERVICE ?= /tmp/ume.service
TMP_ENV ?= /tmp/ume.env

pre:
	go mod tidy
	mkdir -p ./target

test:
	go test ./...

dev:
	go run ./cmd/umed

build: pre
	go build -o ./target ./cmd/umed

build_x64: pre
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o ./target/umed_linux_amd64 ./cmd/umed

deploy-upload: build_x64
	scp $(DEPLOY_BIN) $(DEPLOY_SERVICE) $(DEPLOY_HOST):/tmp/
	if test -n "$(ENV_FILE)"; then scp "$(ENV_FILE)" $(DEPLOY_HOST):$(TMP_ENV); elif test -f "$(DEFAULT_ENV_FILE)"; then scp "$(DEFAULT_ENV_FILE)" $(DEPLOY_HOST):$(TMP_ENV); elif test -f "$(FALLBACK_ENV_FILE)"; then scp "$(FALLBACK_ENV_FILE)" $(DEPLOY_HOST):$(TMP_ENV); else echo "No local env file staged; remote $(REMOTE_ENV) must already exist"; fi

deploy-env:
	if test -n "$(ENV_FILE)"; then scp "$(ENV_FILE)" $(DEPLOY_HOST):$(TMP_ENV); elif test -f "$(DEFAULT_ENV_FILE)"; then scp "$(DEFAULT_ENV_FILE)" $(DEPLOY_HOST):$(TMP_ENV); elif test -f "$(FALLBACK_ENV_FILE)"; then scp "$(FALLBACK_ENV_FILE)" $(DEPLOY_HOST):$(TMP_ENV); else echo "Set ENV_FILE=/path/to/ume.env or create $(DEFAULT_ENV_FILE)" >&2; exit 1; fi
	ssh -t $(DEPLOY_HOST) 'sudo sh -c '"'"'install -d -o root -g root -m 0750 $(REMOTE_ENV_DIR); install -o root -g ume -m 0640 $(TMP_ENV) $(REMOTE_ENV)'"'"''

deploy-install: deploy-upload
	ssh -t $(DEPLOY_HOST) 'sudo sh -c '"'"'set -eu; id -u ume >/dev/null 2>&1 || useradd --system --home-dir $(REMOTE_STATE_DIR) --shell /usr/sbin/nologin ume; install -d -o ume -g ume -m 0750 $(REMOTE_STATE_DIR) $(REMOTE_TMP_DIR); install -d -o root -g root -m 0750 $(REMOTE_ENV_DIR); if test -f $(TMP_ENV); then install -o root -g ume -m 0640 $(TMP_ENV) $(REMOTE_ENV); elif ! test -f $(REMOTE_ENV); then echo "Missing $(REMOTE_ENV); run make deploy-env ENV_FILE=/path/to/ume.env first" >&2; exit 1; fi; install -o root -g root -m 0755 $(TMP_BIN) $(REMOTE_BIN); install -o root -g root -m 0644 $(TMP_SERVICE) $(REMOTE_SERVICE); systemctl daemon-reload; systemctl enable --now ume.service; systemctl restart ume.service; systemctl --no-pager --full status ume.service'"'"''

deploy-status:
	ssh $(DEPLOY_HOST) 'systemctl is-enabled ume.service; systemctl is-active ume.service; systemctl --no-pager --full status ume.service'

deploy: deploy-install

.PHONY: pre test dev build build_x64 deploy-upload deploy-env deploy-install deploy-status deploy
