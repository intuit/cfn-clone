#!/bin/bash

set -e

OS_TYPES=(darwin linux)
PROJECT_ROOT="$( cd -P "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"
PKG_ROOT="$PROJECT_ROOT/pkg"
BIN_NAME="cfn-clone"

if [ -d "$PKG_ROOT" ]; then
    echo -n "Removing existing binaries..."
    rm -fr $PKG_ROOT
    echo "Done"
fi

mkdir -p $PKG_ROOT

version=$(grep version $BIN_NAME/version.go  | awk '{print $4}' | sed 's/"//g')

echo "Building packages for distribution for version $version."

for os in ${OS_TYPES[*]}; do
    FILE_NAME="$BIN_NAME-$version-$os-amd64"
    echo -n "Building packages for $os/amd64..."
    env GOOS=$os GOARCH=amd64 go build -o $PKG_ROOT/$FILE_NAME ./$BIN_NAME
    cd $PKG_ROOT
    shasum --algorithm 256 --binary $FILE_NAME >> $version-sha256-sums
    cd $PROJECT_ROOT
    echo "Done"
done
