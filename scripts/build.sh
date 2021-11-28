#!/usr/bin/env bash

set -euo pipefail

NAME="slack-lunch-bot"
DIR="$(dirname $0)"
ZIP_NAME="${DIR}/${NAME}.zip"

SLACK_SIGNING_SECRET="$(cat $DIR/secrets/slack-signing-secret.txt)"

GOOS=linux GOARCH=amd64 go build \
  -o "${NAME}" \
  -ldflags "-X lunch/pkg/http/slack.signingSecret=${SLACK_SIGNING_SECRET}" \
  ./cmd/slack-lunch-bot/main.go

zip -q "${ZIP_NAME}" "${NAME}"

rm "${NAME}"

echo "${ZIP_NAME}"
