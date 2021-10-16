#!/usr/bin/env bash

set -eo pipefail

NAME="slack-lunch-bot"
ZIP_NAME="${NAME}.zip"

GOOS=linux GOARCH=amd64 go build \
    -o "${NAME}" \
    ./cmd/slack-lunch-bot/main.go

zip -q "${ZIP_NAME}" "${NAME}"

rm "${NAME}"

echo ${ZIP_NAME}
