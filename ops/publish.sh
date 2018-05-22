#!/usr/bin/env bash

set -o errexit

# Pull the last version
echo "Fetching master..."
git fetch origin master -q

LAST_TAG=`git describe --abbrev=0 --tags`
echo "Last tag is $LAST_TAG"
NEW_TAG=`./semver bump $1 $LAST_TAG`
echo "Tagging with $NEW_TAG..."

# Tag
git checkout origin/master -q
git tag -a $NEW_TAG
git push -q origin $NEW_TAG
