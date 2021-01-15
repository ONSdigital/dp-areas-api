#!/bin/bash -eux

export cwd=$(pwd)

pushd $cwd/dp-areas-api
  make audit
popd 