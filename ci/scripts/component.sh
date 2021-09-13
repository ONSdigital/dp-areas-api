#!/bin/bash -eux

cwd=$(pwd)

pushd $cwd/dp-topic-api
  make test-component
popd
