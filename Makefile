# prepbot build & deploy. Override REMOTE_* on the command line, e.g.
#   make deploy REMOTE_HOST=1.2.3.4 REMOTE_USER=stepan
BINARY      := prepbot
REMOTE_USER ?= stepan
REMOTE_HOST ?= your-vps
REMOTE_DIR  := /opt/prepbot
REMOTE_BIN  := /usr/local/bin/prepbot

.PHONY: build test vet deploy

build: ## static linux/amd64 binary
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -o bin/$(BINARY) ./cmd/prepbot

test: ## unit tests with the race detector
	go test -race ./...

vet:
	go vet ./...

deploy: build ## build, ship binary + data/prompts, restart the service
	scp bin/$(BINARY) $(REMOTE_USER)@$(REMOTE_HOST):/tmp/$(BINARY)
	rsync -az --delete data prompts $(REMOTE_USER)@$(REMOTE_HOST):$(REMOTE_DIR)/
	ssh $(REMOTE_USER)@$(REMOTE_HOST) '\
		sudo install -m 0755 /tmp/$(BINARY) $(REMOTE_BIN) && \
		sudo chown -R prepbot:prepbot $(REMOTE_DIR) && \
		sudo systemctl restart prepbot && \
		sudo systemctl --no-pager status prepbot | head -5'
