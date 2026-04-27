#!/bin/bash

set -euo pipefail

APP_NAME="isopod"
MIGRATE_NAME="migrate"
MAIN_PATH="cmd/app/main.go"
MIGRATE_PATH="cmd/migrate/main.go"
BIN_DIR="bin"
OUTPUT="$BIN_DIR/$APP_NAME"
MIGRATE_OUTPUT="$BIN_DIR/$MIGRATE_NAME"

mkdir -p "$BIN_DIR"

echo "> Running tests..."
go test -v ./...

echo
echo "> Building $MIGRATE_NAME..."
go build -o "$MIGRATE_OUTPUT" "$MIGRATE_PATH"

echo
echo "> Building $APP_NAME..."
go build -o "$OUTPUT" "$MAIN_PATH"

echo
echo "> Running database migrations: $MIGRATE_OUTPUT"
"./$MIGRATE_OUTPUT"

echo
echo "> Execution started: $OUTPUT"
"./$OUTPUT"

