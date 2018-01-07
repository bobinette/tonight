#!/bin/sh
#
# Pre-commit hook running the tests
# Tips&tricks gotten from http://codeinthehole.com/tips/tips-for-using-a-git-pre-commit-hook/

STASH_NAME="pre-commit-$(date +%s)"
git stash save -q --keep-index $STASH_NAME


echo "Running tests. This may take a while. Use --no-verify to skip"
# Test prospective commit
go test ./...
RESULT=$?

STASHES=$(git stash list -n 1)
(echo $STASHES | grep $STASH_NAME && git stash pop -q) || echo "Failed to restore stash"

# Fail if tests failed
[ $RESULT -ne 0 ] && exit 1
exit 0
