#!/bin/bash
set -e

# This script is meant to run on Travis-CI only
if [ -z "$TRAVIS_BRANCH" ]; then 
  echo "ABORTING: this script runs on Travis-CI only"
  exit 1
fi
# Check essential envs
if [ -z "$GITHUB_TOKEN" ]; then
  echo "ABORTING: env GITHUB_TOKEN is missing"
  exit 1
fi
# verbose logging
set -x

# create a build number
export BUILD_NR="$(date '+%Y%m%d-%H%M%S')"
echo "BUILD_NR=$BUILD_NR"

# run build
mkdir builds
mkdir builds/$BUILD_NR/
./build.sh $BUILD_NR
# deploy to GitHub releases
export GIT_TAG=v$BUILD_NR
export GIT_RELTEXT="Auto-released by [Travis-CI build #$TRAVIS_BUILD_NUMBER](https://travis-ci.org/$TRAVIS_REPO_SLUG/builds/$TRAVIS_BUILD_ID)"
curl -sSL https://github.com/tcnksm/ghr/releases/download/v0.13.0/ghr_v0.13.0_linux_amd64.tar.gz > ghr.tar.gz
tar xvf ghr.tar.gz   
./ghr_v0.13.0_linux_amd64/ghr --version
./ghr_v0.13.0_linux_amd64/ghr --debug -u xjx00 -b "$GIT_RELTEXT" $GIT_TAG builds/$BUILD_NR/