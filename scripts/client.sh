#!/usr/bin/env bash

set -e

echo 'client/main.go' | entr -r go run client/main.go
