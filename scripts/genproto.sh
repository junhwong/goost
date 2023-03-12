#!/usr/bin/env bash
#
# Generate all protobuf bindings.
# Run from repository root.
set -e
set -u

if ! [[ "$0" =~ "scripts/genproto.sh" ]]; then
	echo "must be run from repository root"
	exit 255
fi
# https://protobuf.dev/downloads/
# go install github.com/golang/protobuf/protoc-gen-go
# 3.12
# if ! [[ $(protoc --version) =~ "3.15.8" ]]; then
# 	echo "could not find protoc 3.15.8, is it installed + in PATH?"
# 	exit 255
# fi
PROM_ROOT="${PWD}"

protoc --experimental_allow_proto3_optional --proto_path=${PROM_ROOT}/apm/field/pb --go_out=${PROM_ROOT} types.proto