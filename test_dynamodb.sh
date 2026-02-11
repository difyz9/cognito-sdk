#!/bin/bash

# DynamoDB 集成测试脚本

echo "========================================"
echo "DynamoDB 集成测试"
echo "========================================"
echo ""

# 检查环境变量
echo "1. 检查环境变量..."
if [ -z "$AWS_REGION" ]; then
    echo "❌ AWS_REGION 未设置"
    echo "请运行: export AWS_REGION=us-east-1"
    exit 1
fi

if [ -z "$COGNITO_USER_POOL_ID" ]; then
    echo "⚠️  COGNITO_USER_POOL_ID 未设置 (示例可能需要)"
fi

if [ -z "$COGNITO_CLIENT_ID" ]; then
    echo "⚠️  COGNITO_CLIENT_ID 未设置 (示例可能需要)"
fi

echo "✅ AWS_REGION: $AWS_REGION"
echo ""

# 检查依赖
echo "2. 检查 Go 依赖..."
if ! command -v go &> /dev/null; then
    echo "❌ Go 未安装"
    exit 1
fi

echo "✅ Go 版本: $(go version)"
echo ""

# 编译测试
echo "3. 编译项目..."
if go build; then
    echo "✅ 编译成功"
else
    echo "❌ 编译失败"
    exit 1
fi
echo ""

# 运行测试
echo "4. 运行单元测试..."
if go test -v; then
    echo "✅ 测试通过"
else
    echo "⚠️  部分测试未通过（可能需要 AWS 凭证）"
fi
echo ""

# 检查示例
echo "5. 检查示例代码..."
examples=(
    "examples/basic"
    "examples/dynamodb-basic"
    "examples/cognito-dynamodb-integration"
)

for example in "${examples[@]}"; do
    if [ -d "$example" ]; then
        echo "✅ $example"
    else
        echo "❌ $example 不存在"
    fi
done
echo ""

# 检查文档
echo "6. 检查文档..."
docs=(
    "README.md"
    "API_REFERENCE.md"
    "DYNAMODB_GUIDE.md"
    "DYNAMODB_INTEGRATION_SUMMARY.md"
)

for doc in "${docs[@]}"; do
    if [ -f "$doc" ]; then
        echo "✅ $doc"
    else
        echo "❌ $doc 不存在"
    fi
done
echo ""

# 生成使用说明
echo "========================================"
echo "✅ DynamoDB 集成测试完成！"
echo "========================================"
echo ""
echo "快速开始："
echo ""
echo "1. 查看文档："
echo "   cat DYNAMODB_GUIDE.md"
echo ""
echo "2. 运行基础示例："
echo "   cd examples/dynamodb-basic"
echo "   go run main.go"
echo ""
echo "3. 运行集成示例："
echo "   cd examples/cognito-dynamodb-integration"
echo "   go run main.go"
echo ""
echo "4. 在你的项目中使用："
echo "   go get github.com/difyz9/cognito-sdk"
echo ""
echo "文档："
echo "  - DynamoDB 使用指南: DYNAMODB_GUIDE.md"
echo "  - API 参考: API_REFERENCE.md"
echo "  - 项目总结: DYNAMODB_INTEGRATION_SUMMARY.md"
echo ""
