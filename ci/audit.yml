---
platform: linux

image_resource:
  type: docker-image
  source:
    repository: onsdigital/dp-concourse-tools-nancy
    tag: latest

inputs:
  - name: dp-areas-api
    path: dp-areas-api

run:
  path: dp-areas-api/ci/scripts/audit.sh 