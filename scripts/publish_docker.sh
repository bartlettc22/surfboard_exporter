#!/bin/bash

set -eo pipefail

VERSION=${1:-dev}
DOCKER_REPO="bartlettc/surfboard-exporter"

if [[ "${VERSION}" == "master" ]]; then
  docker tag ${DOCKER_REPO}:master ${DOCKER_REPO}:latest
  VERSION="latest"
fi

echo "${DOCKER_PASSWORD}" | docker login -u "${DOCKER_USERNAME}" --password-stdin
docker push ${DOCKER_REPO}:${VERSION}
