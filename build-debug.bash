#!/bin/bash
set -euo pipefail

[ -f ./common.sh ] && . ./common.sh

MODE='Debug'

go build \
	-o ${NAME}d \
	-ldflags " \
		-X main.build_mode=${MODE} \
		-X main.app_version=${VERSION}d \
		-X main.app_name=${NAME}d \
	" \
	${SRCS}
