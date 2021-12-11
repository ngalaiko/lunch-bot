#!/usr/bin/env bash

set -eo pipefail

FUNCTION_NAME="slack-lunch-bot"
FILE_NAME="$1"

UPDATE_CODE_RES=$(aws lambda update-function-code --function-name "${FUNCTION_NAME}" --zip-file "fileb://${FILE_NAME}")
CODE_SHA=$(echo "$UPDATE_CODE_RES" | jq --raw-output '.CodeSha256')

PUBLISH_VERSION_RES=$(aws lambda publish-version --function-name "${FUNCTION_NAME}" --code-sha256 "${CODE_SHA}")
FUNCTION_ARN=$(echo "$PUBLISH_VERSION_RES" | jq --raw-output '.FunctionArn')

echo "${FUNCTION_ARN} updated"
