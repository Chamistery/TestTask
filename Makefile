include .env
LOCAL_BIN:=$(CURDIR)/bin

MAIN_FILES := cmd/auth_service/main.go

# Правило для запуска main.go
run:
	@for file in $(MAIN_FILES); do \
		echo "Running $$file..."; \
		go run $$file & \
	done; \
	wait
install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.20.0
	GOBIN=$(LOCAL_BIN) go install github.com/gojuno/minimock/v3/cmd/minimock@v3

local-migration-auth-status:
	$(LOCAL_BIN)/goose -dir ${MIGRATION_AUTH_DIR} postgres ${PG_AUTH_DSN} status -v

local-migration-auth-up:
	$(LOCAL_BIN)/goose -dir ${MIGRATION_AUTH_DIR} postgres ${PG_AUTH_DSN} up -v

local-migration-auth-down:
	$(LOCAL_BIN)/goose -dir ${MIGRATION_AUTH_DIR} postgres ${PG_AUTH_DSN} down -v

generate-test:
	go generate ./...

test-coverage:
	go clean -testcache
	go test ./... -coverprofile=coverage.tmp.out -covermode count -coverpkg=github.com/Chamistery/TestTask/internal/auth/auth_http/... -count 5
	grep -v 'mocks\|config' coverage.tmp.out  > coverage.out
	rm coverage.tmp.out
	go tool cover -html=coverage.out;
	go tool cover -func=./coverage.out | grep "total";
	grep -sqFx "/coverage.out" .gitignore || echo "/coverage.out" >> .gitignore