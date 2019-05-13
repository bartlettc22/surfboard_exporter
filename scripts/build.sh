#!/bin/bash

set -eo pipefail

VERSION=${1:-dev}
DOCKER_REPO="bartlettc/surfboard-exporter"

# Directory to house our binaries
mkdir -p bin

# Build the binary in Docker
docker build -t ${DOCKER_REPO}:${VERSION} ./

# Run the container in the background in order to extract the binary
docker run --rm --entrypoint "" --name build-bin -d ${DOCKER_REPO}:${VERSION} sh -c "sleep 120"

docker cp build-bin:/usr/bin/surfboard-exporter bin
docker stop build-bin

# Zip up the binary
cd bin
tar -cvzf surfboard-exporter-${VERSION}.tar.gz surfboard-exporter

# Get us back to the root
cd ..
