#!/bin/bash

if [ $# -lt 1 ]; then
  echo "$0: takes at least one argument: 'release <version_number>'"
  exit 1
fi

RELEASE_VERSION=$1

go build
go test

rm -rf releases/$RELEASE_VERSION
mkdir -p releases/$RELEASE_VERSION

echo "Building for Mac"
GOOS=darwin GOARCH=amd64 go build -o fuse
tar -zcvf releases/$RELEASE_VERSION/fuse-v$RELEASE_VERSION-darwin-amd64.tar.gz fuse

echo "Building for Linux"
GOOS=linux GOARCH=amd64 go build -o fuse
tar -zcvf releases/$RELEASE_VERSION/fuse-v$RELEASE_VERSION-linux-amd64.tar.gz fuse

echo "Building for Windows"
GOOS=windows GOARCH=amd64 go build -o fuse
tar -zcvf releases/$RELEASE_VERSION/fuse-v$RELEASE_VERSION-windows-amd64.tar.gz fuse

echo "Release artifact files can be found in releases/$RELEASE_VERSION"
