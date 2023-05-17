#!/bin/bash

set -e

GV="$1"

rm -rf generated
./hack/generate_group.sh "client,lister,informer" github.com/x893675/demo-crd/generated github.com/x893675/demo-crd/apis "${GV}" --output-base=./  -h "$PWD/hack/boilerplate.go.txt"
mv github.com/x893675/demo-crd/generated ./
rm -rf ./github.com