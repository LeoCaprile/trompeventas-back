#! /bin/sh

set -a
export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING=postgres://store:store@localhost:5432/store?sslmode=disable
export GOOSE_MIGRATION_DIR=./db/migrations/

set +a
