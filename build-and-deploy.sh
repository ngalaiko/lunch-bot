#!/usr/bin/env bash

set -eo pipefail

./deploy.sh $( ./build.sh )
