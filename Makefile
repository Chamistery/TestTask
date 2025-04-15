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

local-migration-status:
	$(LOCAL_BIN)/goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} status -v

local-migration-hotel-up:
	$(LOCAL_BIN)/goose -dir ${MIGRATION_HOTEL_DIR} postgres ${PG_HOTEL_DSN} up -v

local-migration-booking-up:
	$(LOCAL_BIN)/goose -dir ${MIGRATION_BOOKING_DIR} postgres ${PG_BOOKING_DSN} up -v

local-migration-booking-down:
	$(LOCAL_BIN)/goose -dir ${MIGRATION_BOOKING_DIR} postgres ${PG_BOOKING_DSN} down -v

local-migration-auth-up:
	$(LOCAL_BIN)/goose -dir ${MIGRATION_AUTH_DIR} postgres ${PG_AUTH_DSN} up -v
local-migration-auth-down:
	$(LOCAL_BIN)/goose -dir ${MIGRATION_AUTH_DIR} postgres ${PG_AUTH_DSN} down -v

local-migration-hotel-down:
	$(LOCAL_BIN)/goose -dir ${MIGRATION_HOTEL_DIR} postgres ${PG_HOTEL_DSN} down -v

local-migration-up:
	local-migration-hotel-up
	local-migration-booking-up
	local-migration-auth-up

local-migration-down:
	local-migration-hotel-down
	local-migration-booking-down
	local-migration-auth-down