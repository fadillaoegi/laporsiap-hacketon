#!/bin/sh
set -eu
cd "$(dirname "$0")/.."
set -a
. ./.env
if [ -f ../../.apertus-api.env ]; then
  . ../../.apertus-api.env
fi
set +a
cd src/laporsiap_BE
exec go run ./cmd/api
