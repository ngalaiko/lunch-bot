#!/bin/sh

lunch_place=$(sort -R ./lunch | head -1)

curl -X POST \
    -H 'Content-type: application/json' \
    --data "{\"text\": \"Today's lunch place is ${lunch_place}\"}" \
    "https://hooks.slack.com/services/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"
