#!/bin/sh
set -eu
cd "$(dirname "$0")/.."
set -a
. ./.env
set +a
: "${TEST_DATABASE_URL:?Set TEST_DATABASE_URL to the separate laporsiap_test database}"
cd src/laporsiap_BE
exec go test -race -count=1 ./...
