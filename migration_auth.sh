#!/bin/bash
source .env

export MIGRATION_DSN="host=pg-auth port=$PORT_TO dbname=$POSTGRES_AUTH_DB user=$POSTGRES_AUTH_USER password=$POSTGRES_AUTH_PASSWORD sslmode=disable"

sleep 2 && goose -dir "${MIGRATION_AUTH_DIR}" postgres "${MIGRATION_DSN}" up -v