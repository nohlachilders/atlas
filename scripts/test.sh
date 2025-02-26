#!/usr/bin/env bash
pushd sql/schema/
goose up
popd
go test ./...
