#!/usr/bin/env bash

SERVER=127.0.0.1:8081

if [[ $1 != "" ]]
then
  SERVER=$1
fi

go run integration/ping.go $SERVER
