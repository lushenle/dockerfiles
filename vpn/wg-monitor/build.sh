#!/bin/bash
#
REGISTRY="harbor.internal.moqi.ai"
PROJECT="itimages"
REPOSITORY="wireguard_exporter"
TAG=("go" "rust")

for i in "${TAG[@]}"; do
    echo -e "Building image: ${REGISTRY}/${PROJECT}/${REPOSITORY}:${i}"
    #docker build -t ${REGISTRY}/${PROJECT}/${REPOSITORY}:${i} -f Dockerfile-${i}
    echo ""
    echo -e "Push image: ${REGISTRY}/${PROJECT}/${REPOSITORY}:${i}"
    #docker push ${REGISTRY}/${PROJECT}/${REPOSITORY}:${i}
    echo ""
done

