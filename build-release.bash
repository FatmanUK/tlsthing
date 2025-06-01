#!/bin/bash
set -euo pipefail

[ -f ./common.sh ] && . ./common.sh

MODE='Release'

go build \
	-o ${NAME} \
	-ldflags " \
		-X main.build_mode=${MODE} \
		-X main.app_version=${VERSION} \
		-X main.app_name=${NAME} \
	" \
	${SRCS}
