#!/bin/bash

set -eo pipefail

VERSION=${1:-dev}
DOCKER_REPO="bartlettc/surfboard-exporter"

echo "${DOCKER_PASSWORD}" | docker login -u "${DOCKER_USERNAME}" --password-stdin
docker push ${DOCKER_REPO}:${VERSION}
