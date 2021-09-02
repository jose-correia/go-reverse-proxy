#!/bin/bash

set -ex

mkdir -p bin
GO_BIN_DIR="bin/"

echo "[$(date)] Discovering main packages..."
MAIN_PACKAGES="$(go list -f '{{if eq .Name "main" }}{{ .ImportPath }}{{ end }}' ./...)"
if [[ -z "$MAIN_PACKAGES" ]]; then
    exit 0
fi

echo "Found main packages -"
echo "$MAIN_PACKAGES"
echo

echo "[$(date)] Installing..."
GIT_COMMIT_COUNT=$(git rev-list HEAD --count)
GIT_COMMIT_HASH=$(git rev-parse HEAD)
while read -r package; do
    echo "Building $package"
    go build \
        -o ${GO_BIN_DIR}/$(basename "$package") \
        -ldflags "-X main.version=${GIT_COMMIT_COUNT} -X main.gitHash=${GIT_COMMIT_HASH}" \
        "$package"
done <<< "$MAIN_PACKAGES"
