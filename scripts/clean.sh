#!/usr/bin/env bash

set -ex

rm -rf .git
git init
git add .
git commit -m "initial commit"
git log
git status
