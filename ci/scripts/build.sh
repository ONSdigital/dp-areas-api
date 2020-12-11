#!/bin/bash -eux

pushd dp-areas-api
  make build
  cp build/dp-areas-api Dockerfile.concourse ../build
popd
