#!/usr/bin/env bash

SLEEP=0.5

set -e
echo "" > coverage.txt

for d in $(go list ./... | grep -v vendor); do
    go test -race -coverprofile=profile.out -covermode=atomic $d
    if [ -f profile.out ]; then
      cat profile.out >> coverage.txt
      rm profile.out
    fi
done

# Plugin related tests can't run under race condition
TEST_PLUGIN=true go test ./vm --run .?Plugin.? -v -coverprofile=profile.out -covermode=atomic $d

if [ -f profile.out ]; then
  cat profile.out >> coverage.txt
  rm profile.out
fi


# Test if libs that require built in Goby script would work.
# TODO: Write a test for this specific case
go install .
goby test_fixtures/server.gb & PID=$!
echo "Sleeping for $SLEEP sec to wait server.gb being ready..."; sleep $SLEEP

ab -n 3000 -c 100 http://localhost:3000/

kill $PID
