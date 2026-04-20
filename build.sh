#!/bin/bash

set -euo pipefail

APP_NAME="isopod"
MAIN_PATH="cmd/app/main.go"
BIN_DIR="bin"
OUTPUT="$BIN_DIR/$APP_NAME"

mkdir -p "$BIN_DIR"

echo "> Running tests..."
go test -v ./...

echo ""> Building $APP_NAME..."
go build -o "$OUTPUT" "$MAIN_PATH"

echo "> Execution started: $OUTPUT"
"./$OUTPUT"

