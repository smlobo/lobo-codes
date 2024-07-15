#/bin/bash

set -ex

PROJECT_NAME="h-tree"

BUILD_DIR=${1}

mkdir install
cp ${BUILD_DIR}/${PROJECT_NAME}.js .
cp ${BUILD_DIR}/${PROJECT_NAME}.wasm .
cp ${BUILD_DIR}/${PROJECT_NAME}.html index.html
