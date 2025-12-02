#!/bin/bash
set -euo pipefail

# BUILD_VERSION
eval "export BUILD_$(grep ^VERSION Tupfile)"

export REMOTE_HOST='ghcr.io'
export NOW="$(date -Iseconds | cut -c 1-19 | tr ' :' '__')Z"
export TAG_DT="$(<<<${NOW} cut -c 1-10)"
export COMMON_BUILD_FLAGS="\
	--build-context=build-release=build-linux-release \
	--http-proxy=false --no-cache=false --compress=true \
	--layers=true --format=oci --log-level=info"

# DIRECTOR_NAME
eval "export $(grep DIRECTOR_NAME buildconfigs/linux-release.conf \
	| sed 's/^CONFIG_//')"
export DIRECTOR_URL="${REMOTE_HOST}/dreamtrack-net/${DIRECTOR_NAME}"

# CLIENT_NAME
eval "export $(grep CLIENT_NAME buildconfigs/linux-release.conf \
	| sed 's/^CONFIG_//')"
export CLIENT_URL="${REMOTE_HOST}/dreamtrack-net/${CLIENT_NAME}"

export DIRECTOR_FLAGS="\
	${COMMON_BUILD_FLAGS} --tag="${DIRECTOR_URL}:latest" \
	--tag="${DIRECTOR_URL}:${TAG_DT}" \
	--tag="${DIRECTOR_URL}:${NOW}" \
	--tag="${DIRECTOR_URL}:v${BUILD_VERSION}" \
	"

export CLIENT_FLAGS="\
	${COMMON_BUILD_FLAGS} --tag="${CLIENT_URL}:latest" \
	--tag="${CLIENT_URL}:${TAG_DT}" \
	--tag="${CLIENT_URL}:${NOW}" \
	--tag="${CLIENT_URL}:v${BUILD_VERSION}" \
	"

# ❯ <<<"FROM ghcr.io/dreamtrack-net/tlsthing-director:latest" \
#   podman build --output type=tar,dest=d.tar - && rm -rf tars/* \
#   && tar xf d.tar -C tars/

# ❯ <<<"FROM ghcr.io/dreamtrack-net/tlsthing-client:latest" \
#   podman build --output type=tar,dest=c.tar - && rm -rf tars/* \
#   && tar xf c.tar -C tars/

cat                                      \
	container/Containerfile.head     \
	container/director/Containerfile \
	container/Containerfile.tail     \
	>/dev/shm/Containerfile.director

cat                                    \
	container/Containerfile.head   \
	container/client/Containerfile \
	container/Containerfile.tail   \
	>/dev/shm/Containerfile.client

podman build                               \
	${DIRECTOR_FLAGS}                  \
	-f=/dev/shm/Containerfile.director \
	--tag=${DIRECTOR_NAME}

podman build                             \
	${CLIENT_FLAGS}                  \
	-f=/dev/shm/Containerfile.client \
	--tag=${CLIENT_NAME}

# --runtime-flag debug --log-level debug

# $ </dev/shm/pw more | podman secret create github-containers-pat -
podman login                                       \
	-u FatmanUK --secret github-containers-pat \
	${REMOTE_HOST}

podman push ${DIRECTOR_URL}:latest
podman push ${DIRECTOR_URL}:${TAG_DT}
podman push ${DIRECTOR_URL}:${NOW}
podman push ${DIRECTOR_URL}:v${BUILD_VERSION}

podman push ${CLIENT_URL}:latest
podman push ${CLIENT_URL}:${TAG_DT}
podman push ${CLIENT_URL}:${NOW}
podman push ${CLIENT_URL}:v${BUILD_VERSION}
