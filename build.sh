#!/bin/bash

APP_NAME=fourclover
APP_VERSION=v010

rm -rf build
mkdir build

declare -a platforms=(
  "linux/amd64"
  "linux/386"
  "darwin/amd64"
  "windows/amd64"
  "windows/386"
)

for platform in "${platforms[@]}"; do
  GOOS=${platform%/*}
  GOARCH=${platform#*/}
  output_name=${APP_NAME}_${APP_VERSION}_${GOOS}_${GOARCH}
  if [ $GOOS = "windows" ]; then
    output_name+=.exe
  fi

  env CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build -ldflags '-extldflags=-fno-PIC' -o build/$output_name
done

# Change directory to build folder
cd build

# Loop through all files in the folder
for file in *; do
  # Calculate sha256 checksum and store in variable
  checksum=$(sha256sum "$file" | awk '{print $1}')
  
  # Write checksum and filename to checksum.txt
  echo "$checksum, $file" >> checksums.txt
done

# Exit build folder
cd ..

echo "Build complete"