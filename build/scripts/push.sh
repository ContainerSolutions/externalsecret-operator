#!/bin/bash

set -ex

GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD | sed 's/\//-/')
GIT_HASH=$(git rev-parse --short HEAD) 
GIT_TAG=$(git describe --tags --abbrev=0 || true)

tag_and_push() {
    docker tag $1 $1:$2
    docker push $1:$2
}

# tag and push with its GIT_HASH
tag_and_push $DOCKER_IMAGE $GIT_HASH

# if we are on master, tag and push with latest
if [[ "$GIT_BRANCH" == "master" ]] ; then
    tag_and_push $DOCKER_IMAGE latest

    # if we are on master, and there's a tag, tag and push with the version tag
    if [[ "$GIT_TAG" == v[0-9]* ]] ; then
        tag_and_push $DOCKER_IMAGE ${GIT_TAG//v}
    fi
else
    tag_and_push $DOCKER_IMAGE $GIT_BRANCH-latest
fi
