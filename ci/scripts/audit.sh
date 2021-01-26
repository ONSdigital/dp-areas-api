#!/bin/bash -eux

export cwd=$(pwd)

pushd $cwd/dp-topic-api
  make audit
popd
