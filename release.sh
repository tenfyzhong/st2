#!/usr/bin/env bash
#################################################################
#
#    file: release.sh
#   brief: release.sh is a script to create a github release.
#          It build the release asset and the upload it to the release.
#  author: tenfyzhong
#   email: tenfy@tenfy.cn
# created: 2023-01-06 11:43:12
#
#################################################################

set -x

CWD=$(pwd)
COMMAND=st2
CMD_PATH=cmd/st2
OUTPUT=${CWD}/output
VERSION_KEY=github.com/tenfyzhong/st2/cmd/st2/config.Version

VERSION=$(git describe --tags --abbrev=0 2>/dev/null)
if [ $? -ne 0 ]; then
    echo 'Version not found'
    exit 1
fi
GIT_TAG_MESSAGE=$(git tag -l --format='%(contents)' "$VERSION" 2>/dev/null)

gobuild() {
    OS=$1
    ARCH=$2
    CGO_ENABLED=0 GOOS="$OS" GOARCH="$ARCH" go build -ldflags "-X '${VERSION_KEY}=${VERSION}'" -o "${OUTPUT}/${COMMAND}-${OS}-${ARCH}"
}

build() {
    declare -a OS_ARRAY=(linux windows)
    declare -a ARCH_ARRAY=(amd64 386 arm)

    cd "${CWD}/${CMD_PATH}" || exit

    rm -rf "$OUTPUT"
    mkdir -p "$OUTPUT"

    go install

    gobuild darwin amd64

    for OS in "${OS_ARRAY[@]}"; do
        for ARCH in "${ARCH_ARRAY[@]}"; do
            gobuild "$OS" "$ARCH"
        done
    done
    cd "$CWD" || exit
}

release() {
    gh release create "$VERSION" -d --verify-tag -n "$GIT_TAG_MESSAGE" "$OUTPUT"/*
}

build
release
