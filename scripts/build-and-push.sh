#!/usr/bin/env bash

set -euo pipefail

REPOSITORY_URI="794038378494.dkr.ecr.eu-north-1.amazonaws.com/lunch"
VERSION="$(date +%Y-%m-%d-%H-%M-%S)"

VERSION_IMAGE="${REPOSITORY_URI}:${VERSION}"
LATEST_IMAGE="${REPOSITORY_URI}:latest"

aws ecr get-login-password --region eu-north-1 |
  docker login --username AWS --password-stdin "${REPOSITORY_URI}"

docker buildx build \
  --platform linux/arm64 \
  --tag "${VERSION_IMAGE}" \
  --tag "${LATEST_IMAGE}" \
  --push \
  .

echo "Pushed ${VERSION_IMAGE}"
echo "Pushed ${LATEST_IMAGE}"
