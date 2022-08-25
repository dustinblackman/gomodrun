#!/usr/bin/env bash

set -e

VERSION="$1"
DIST="$2"
TMPDIR="$(mktemp -d)"

cd "$TMPDIR"
git clone --depth 1 git@github.com:dustinblackman/apt.git .
./add-debs.sh gomodrun "$VERSION" "$DIST"
cd ~
rm -rf "$TMPDIR"
