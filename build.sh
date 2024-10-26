#!/bin/bash

VERSION=$(git describe --tags --always)
COMMIT=$(git rev-parse --short HEAD)
BUILD_DATE=$(date +%Y-%m-%d)

DIST_DIR="dist"
mkdir -p "$DIST_DIR"

build_and_package() {
  local os=$1
  local arch=$2
  local output_name="nexora-${os}-${arch}"

  GOOS=$os GOARCH=$arch go build -ldflags "-X github.com/rdmnl/nexora/version.Version=$VERSION -X github.com/rdmnl/nexora/version.BuildDate=$BUILD_DATE -X github.com/rdmnl/nexora/version.Commit=$COMMIT" -o nexora ./cmd

  mkdir -p "$output_name"
  mv nexora "$output_name/"
  cp config.yaml "$output_name/"
  tar -czvf "$DIST_DIR/${output_name}-${VERSION}.tar.gz" "$output_name/"
  rm -rf "$output_name"
  
  sed -i '' "s|https://github.com/rdmnl/nexora/releases/download/.*/nexora.tar.gz|https://github.com/rdmnl/nexora/releases/download/$VERSION/nexora-${os}-${arch}-${VERSION}.tar.gz|" README.md
}

build_and_package linux amd64
build_and_package linux arm64
build_and_package darwin amd64
build_and_package darwin arm64

echo "Build and packaging completed for Linux and macOS (amd64 and arm64)."