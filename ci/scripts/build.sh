#!/bin/bash -eux

pushd dp-topic-api
  make build
  cp build/dp-topic-api Dockerfile.concourse ../build
popd
