#!/bin/bash

# 定义输出的基础目录
OUTPUT_BASE="./bin"

# 定义不同的平台
PLATFORMS=("windows/amd64" "linux/amd64" "darwin/amd64" "darwin/arm64" "linux/arm64")

# 为每个平台编译并指定输出路径
for PLATFORM in "${PLATFORMS[@]}"; do
    GOOS=${PLATFORM%/*}
    GOARCH=${PLATFORM#*/}
    OUTPUT_DIR="${OUTPUT_BASE}/${GOOS}/${GOARCH}"
    mkdir -p $OUTPUT_DIR  # 确保输出目录存在
    OUTPUT="${OUTPUT_DIR}/warp-plus"
    if [ $GOOS = "windows" ]; then
        OUTPUT+='.exe'
    fi
    env GOOS=$GOOS GOARCH=$GOARCH go build -o $OUTPUT
    echo 'Built:' $OUTPUT
done
