#!/usr/bin/env bash

set -e
echo "" > coverage.txt

for d in $(go list ./... | grep -v vendor); do
    go test -race -coverprofile=profile.out -covermode=atomic $d
    if [ -f profile.out ]; then
      cat profile.out >> coverage.txt
      rm profile.out
    fi
done


# Test if libs that require built in Goby script would work.
# TODO: Write a test for this specific case
go install .
goby test_fixtures/server.gb & PID=$!

sleep 2

kill $PID