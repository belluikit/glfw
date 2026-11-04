#!/bin/sh

# This script computes the hash of the directory called "glfw", relative to $PWD
# at time of invocation, and writes it into the upstreamTreeSHA constant in
# glfw_tree_rebuild.go. Invoke this in a directory containing
# glfw_tree_rebuild.go (go generate does, from the repository root).
#
# TestUpstreamTreeSHAIsCurrent computes the same hash and fails when the
# constant is stale.

set -e

GIT_INDEX_FILE=$(mktemp)
export GIT_INDEX_FILE
rm "$GIT_INDEX_FILE"
git add glfw
HASH=$(git write-tree)
rm "$GIT_INDEX_FILE"

# -i.bak works with both GNU and BSD sed; a bare -i does not.
sed -E -i.bak 's/const upstreamTreeSHA = "[0-9a-f]+"/const upstreamTreeSHA = "'"${HASH}"'"/' glfw_tree_rebuild.go
rm glfw_tree_rebuild.go.bak
