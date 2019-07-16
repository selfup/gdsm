#!/usr/bin/env bash

set -e

find . -name \*.go -print | entr -r docker-compose up --build --scale workers=2
