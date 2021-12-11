#!/usr/bin/env bash

set -euo pipefail

DIR="$(dirname $0)"

$DIR/deploy.sh $( "$DIR/build.sh" )
