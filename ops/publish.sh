#!/usr/bin/env bash

set -o errexit

# Retrieve latest tag
LAST_TAG=`git describe --abbrev=0 --tags`

# Use semver to bump the version
NEW_TAG=`./semver bump $1 $LAST_TAG`

# Pull the last version
echo "Fetching master..."
git fetch origin master
git checkout origin/master

# Tag
echo "Tagging with $NEW_TAG..."
git tag -a $NEW_TAG
git push origin $NEW_TAG
