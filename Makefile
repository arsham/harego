help: ## Show help messages.
	@grep -E '^[0-9a-zA-Z_-]+:(.*?## .*)?$$' $(MAKEFILE_LIST) | sed 's/^Makefile://' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


rabbitmq_image=rabbitmq:3.8-management-alpine
rabbitmq_container=harego_rabbit_1
rabbitmq_data=harego_rabbit_data_1

run="."
dir="./..."
short="-short"
flags=""
timeout=40s

include ./config/dev.env
export $(shell sed 's/=.*//' ./config/dev.env)


.PHONY: unittest
unittest: ## Run unit tests in watch mode. You can set: [run, timeout, short, dir, flags]. Example: make unittest flags="-race".
	@echo "running tests on $(run). waiting for changes..."
	@-zsh -c "go test -trimpath --timeout=$(timeout) $(short) $(dir) -run $(run) $(flags); repeat 100 printf '#'; echo"
	@reflex -d none -r "(\.go$$)|(go.mod)" -- zsh -c "go test -trimpath --timeout=$(timeout) $(short) $(dir) -run $(run) $(flags); repeat 100 printf '#'"


.PHONY: integration_test
integration_test: ## Run integration in watch mode. You can set: [run, timeout, short, dir, flags]. Example: make integration_test flags="-race".
	@-docker start $(rabbitmq_container)
	@echo "running tests on $(run). waiting for changes..."
	@-zsh -c "go test -trimpath --timeout=$(timeout) -tags=integration $(short) $(dir) -run $(run) $(flags); repeat 100 printf '#'; echo"
	@reflex -d none -r "(\.go$$)|(go.mod)" -- zsh -c "go test -trimpath -failfast --timeout=$(timeout) -tags=integration $(short) $(dir) -run $(run) $(flags); repeat 100 printf '#'"


.PHONY: ci_tests
ci_tests: ## Run tests for CI.
	go fmt ./...
	go vet ./...
	golangci-lint run ./...
	go test -trimpath --timeout=5m -failfast -v -tags=integration -race ./...
	go test -trimpath --timeout=5m -failfast -tags=integration -covermode=count -coverprofile=coverage.out ./...


.PHONY: integration_deps
integration_deps: ## Install integration test databases. It removes every existing setup.
	@-docker pull $(rabbitmq_image)
	@-docker network create harego
	@-docker run -d --net harego -p $(RABBITMQ_PORT):5672 -p $(RABBITMQ_ADMIN_PORT):15672 --name $(rabbitmq_container) --hostname $(rabbitmq_container) -e RABBITMQ_ERLANG_COOKIE='harego' -e RABBITMQ_DEFAULT_USER=$(RABBITMQ_USER) -e RABBITMQ_DEFAULT_PASS=$(RABBITMQ_PASSWORD) $(rabbitmq_image)
	@docker exec -it $(rabbitmq_container) rabbitmqctl wait /var/lib/rabbitmq/mnesia/rabbit@$(rabbitmq_container).pid
	@-docker exec -it $(rabbitmq_container) rabbitmqctl set_permissions "$(RABBITMQ_USER)" ".*" ".*" ".*"


.PHONY: reset_docker
reset_docker: ## Reset containers and delete their data.
	@-docker rm -f $(rabbitmq_container)
	@-docker network rm harego


.PHONY: dependencies
dependencies: ## Install dependencies requried for development operations.
	@go get -u github.com/cespare/reflex
	@go get -u github.com/git-chglog/git-chglog/cmd/git-chglog
	@go get github.com/stretchr/testify/mock
	@go get github.com/vektra/mockery/.../
	@go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.29.0
	@go get -u golang.org/x/tools/cmd/stringer
	@go mod tidy


.PHONY: mocks
mocks: ## Generate mocks in all packages.
	@go generate ./...


.PHONY: changelog
changelog: ## Update the changelog.
	@git-chglog > CHANGELOG.md
	@echo "Changelog has been updated."


.PHONY: changelog_release
changelog_release: ## Update the changelog with a release tag.
	@git-chglog --next-tag $(tag) > CHANGELOG.md
	@echo "Changelog has been updated."


.PHONY: clean
clean: ## Clean test caches and tidy up modules.
	@go clean -testcache
	@go mod tidy


.PHONY: coverage
coverage: ## Show the test coverage on browser.
	go test -covermode=count -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | tail -n 1
	go tool cover -html=coverage.out
