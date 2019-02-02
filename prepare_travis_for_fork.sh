#!/usr/bin/env bash

set -o errexit

# On forked repositories, the code will be cloned under `$GOPATH/src/github.com/<fork_owner>/<fork_name>`,
# which is not compatible with the Gooby defined imports.
# The Travis `go_import_path` seems not to solve the problem, as it will end up creating
# `$GOPATH/src/github.com/gooby-lang/gooby/<go_import_path>`.
# In order to solve this problem, we simply move the data in the project where in the intended location,
# which is `$GOBY_ROOT` (`$GOPATH/src/github.com/gooby-lang/gooby`).

if [[ ${TRAVIS_REPO_SLUG} != "gooby-lang/gooby" ]]; then
  mkdir -p "$(dirname ${GOBY_ROOT})"
  mv $HOME/gopath/src/github.com/${TRAVIS_REPO_SLUG} ${GOBY_ROOT}
  cd ${GOBY_ROOT}
fi
