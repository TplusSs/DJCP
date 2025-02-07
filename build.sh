#!/bin/bash

# 定义目标操作系统和架构的组合
targets=(
    "linux amd64"
    "windows amd64"
    "darwin amd64"
    "darwin arm64"
    "linux arm"
)

# 遍历所有目标组合
for target in "${targets[@]}"; do
    # 分割目标组合为操作系统和架构
    os=$(echo $target | cut -d' ' -f1)
    arch=$(echo $target | cut -d' ' -f2)

    # 生成输出文件名
    output="ZEDB_${os}_${arch}_v1.0.0"
    if [ "$os" == "windows" ]; then
        output="${output}.exe"
    fi

    # 编译程序
    GOOS=$os GOARCH=$arch go build -o $output ./main.go
done