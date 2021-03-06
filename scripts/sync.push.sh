#!/usr/bin/env bash

set -e

ORIGIN=git@github.com:selfup/gdsm

git remote -v | grep $ORIGIN

if [[ $? != "0" ]]
then
  git remote add origin $ORIGIN
  git push -u origin master
else
  git push origin master
fi

./scripts/release.sh
