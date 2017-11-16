#!/usr/bin/env bash

SLEEP=0.5

set -e
echo "" > coverage.txt

for d in $(go list ./...); do
    if [ $d == "github.com/goby-lang/goby/vm" ]; then
        # Test vm's code without running race detection because that breaks plugin tests.
        # This can generate full coverage report of vm package.
        # Test that need to run without race detection include NoRaceDetection in the name,
        # otherwise, they will run twice (in the run below).
        TEST_PLUGIN=true go test -coverprofile=profile.out -covermode=atomic $d -v -run NoRaceDetection
        if [ -f profile.out ]; then
          cat profile.out >> coverage.txt
          rm profile.out
        fi

        # Then we test the race condition without plugin tests.
        go test -race $d -v cover
        continue
    fi
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
echo "Sleeping for $SLEEP sec to wait server.gb being ready..."; sleep $SLEEP

ab -n 3000 -c 100 http://localhost:3000/

kill $PID
