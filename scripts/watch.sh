#!/usr/bin/env bash

set -e

find . -name \*.go -print | entr -r go run cmd/daemon/main.go
