#!/usr/bin/env bash
ZC="0dns"
read -p "Provide the docker image tag name: " TAG
read -p "Provide the github organisation name[default:-0chaintest]: " organisation
echo "${organisation:-0chaintest}/${ZC}:$TAG"

if [ -n "$TAG" ]; then
echo " $TAG is the tage name provided"
REGISTRY_IMAGE="${organisation:-0chaintest}/${ZC}"
sudo docker system info | grep -E 'Username' 1>/dev/null
if [[ $? -ne 0 ]]; then
  docker login
fi

sudo docker build -t ${REGISTRY_IMAGE}:${TAG} -f docker.local/Dockerfile .
sudo docker pull ${REGISTRY_IMAGE}:latest
sudo docker tag ${REGISTRY_IMAGE}:latest ${REGISTRY_IMAGE}:stable_latest
echo "Re-tagging the remote latest tag to stable_latest"
sudo docker push ${REGISTRY_IMAGE}:stable_latest
sudo docker tag ${REGISTRY_IMAGE}:${TAG} ${REGISTRY_IMAGE}:latest
echo "Pushing the new latest tag to dockerhub"
sudo docker push ${REGISTRY_IMAGE}:latest
echo "Pushing the new tag to dockerhub tagged as ${REGISTRY_IMAGE}:${TAG}"
sudo docker push ${REGISTRY_IMAGE}:${TAG}
fi
