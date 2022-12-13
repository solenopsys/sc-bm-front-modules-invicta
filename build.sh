#!/bin/sh

build_push(){
  docker buildx build  --platform ${ARCHS} -t ${REGISTRY}/${NAME}:latest   --push .
}

helm_build_push(){
  FN=${NAME}-${VER}.tgz
  helm package ./install --version ${VER}
  curl --data-binary "@${FN}" http://helm.alexstorm.solenopsys.org/api/charts
  rm ${FN}
}

REGISTRY=registry.alexstorm.solenopsys.org
NAME=alexstorm-front-modules
ARCHS="linux/amd64,linux/arm64"
VER=0.1.18


build_push
helm_build_push





