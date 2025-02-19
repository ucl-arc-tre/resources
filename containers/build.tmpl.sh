#!/bin/bash
# Templated shell script used to build and save container images

set -x trace

docker build \
  --tag {{ .Image }} \
  {{ .Context }}

docker save \
  {{ .Image }} \
  | gzip > {{ .FileName }}
