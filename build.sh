#!/usr/bin/env bash

set -eo pipefail

NAME="slack-lunch-bot"

GOOS=linux GOARCH=amd64 go build \
    -o "${NAME}" \
    ./cmd/slack-lunch-bot/main.go

zip "${NAME}.zip" "${NAME}"

rm "${NAME}"
