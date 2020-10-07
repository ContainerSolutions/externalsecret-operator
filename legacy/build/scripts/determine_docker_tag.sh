#!/bin/bash

# Determines Docker tag based on git hash, branch and tag

GIT_HASH=$1
GIT_BRANCH=$2
GIT_TAG=$3

if [ "$GIT_HASH" == "" ]; then
    echo "ERROR $0 expects GIT_HASH, GIT_BRANCH and GIT_TAG as program arguments"
    exit 1
fi

if [ "$GIT_BRANCH" == "" ]; then
    echo "ERROR $0 expects GIT_HASH, GIT_BRANCH and GIT_TAG as program arguments"
    exit 1
fi

if [ "$GIT_TAG" == "" ]; then
    echo "ERROR $0 expects GIT_HASH, GIT_BRANCH and GIT_TAG as program arguments"
    exit 1
fi

# if we are on master use latest
if [[ "$GIT_BRANCH" == "master" ]]; then
    DOCKER_TAG="latest"

    # if we have a git tag that matches a version use that
    if [[ "$GIT_TAG" == v[0-9]* ]]; then
        DOCKER_TAG="${GIT_TAG//v}"
    fi
else
# if we are on branch use branch name + 'latest'
    DOCKER_TAG="$GIT_BRANCH-latest"
fi

echo "$DOCKER_TAG"
