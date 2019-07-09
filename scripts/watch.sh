#!/usr/bin/env bash

set -e

echo 'main.go' | entr -r go run main.go
