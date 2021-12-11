#!/usr/bin/env bash

set -euo pipefail

NAME="slack-lunch-bot"
DIR="$(dirname $0)"
ZIP_NAME="${DIR}/${NAME}.zip"

GOOS=linux GOARCH=amd64 go build \
  -o "${NAME}" \
  ./cmd/lambda/main.go

zip -q "${ZIP_NAME}" "${NAME}"

rm "${NAME}"

echo "${ZIP_NAME}"
