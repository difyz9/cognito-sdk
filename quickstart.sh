#!/bin/bash

# AWS Cognito Go SDK 快速开始脚本

echo "=========================================="
echo "AWS Cognito Go SDK 使用演示"
echo "=========================================="
echo ""

# 检查环境变量
if [ -z "$COGNITO_USER_POOL_ID" ] || [ -z "$COGNITO_CLIENT_ID" ] || [ -z "$AWS_REGION" ]; then
    echo "错误: 请先设置以下环境变量："
    echo "  export COGNITO_USER_POOL_ID=\"your-pool-id\""
    echo "  export COGNITO_CLIENT_ID=\"your-client-id\""
    echo "  export AWS_REGION=\"ap-east-1\""
    echo "  export AWS_ACCESS_KEY_ID=\"your-access-key\""
    echo "  export AWS_SECRET_ACCESS_KEY=\"your-secret-key\""
    exit 1
fi

echo "✓ 环境变量已配置"
echo "  User Pool ID: $COGNITO_USER_POOL_ID"
echo "  Client ID: $COGNITO_CLIENT_ID"
echo "  Region: $AWS_REGION"
echo ""

# 创建测试项目
PROJECT_DIR="test-cognito-sdk"
echo "创建测试项目: $PROJECT_DIR"

if [ -d "$PROJECT_DIR" ]; then
    rm -rf "$PROJECT_DIR"
fi

mkdir -p "$PROJECT_DIR"
cd "$PROJECT_DIR"

# 初始化Go模块
echo "初始化Go模块..."
go mod init test-cognito-sdk

# 创建main.go
cat > main.go << 'EOF'
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	cognito "github.com/difyz9/cognito-sdk"
)

func main() {
	// 创建客户端
	client, err := cognito.NewClient(cognito.Config{
		UserPoolID: os.Getenv("COGNITO_USER_POOL_ID"),
		ClientID:   os.Getenv("COGNITO_CLIENT_ID"),
		Region:     os.Getenv("AWS_REGION"),
	})
	if err != nil {
		log.Fatalf("❌ 创建客户端失败: %v", err)
	}
	fmt.Println("✓ Cognito客户端创建成功")

	// 生成测试邮箱
	email := fmt.Sprintf("test_%d@example.com", time.Now().Unix())
	password := "Test@123456"

	// 1. 注册用户
	fmt.Println("\n=== 测试注册 ===")
	registerResp, err := client.Register(email, password, "")
	if err != nil {
		log.Printf("注册失败: %v", err)
	} else {
		fmt.Printf("✓ 注册成功! UserSub: %s\n", registerResp.UserSub)
	}

	// 等待用户创建完成
	time.Sleep(2 * time.Second)

	// 2. 登录
	fmt.Println("\n=== 测试登录 ===")
	loginResp, err := client.Login(email, password)
	if err != nil {
		log.Fatalf("❌ 登录失败: %v", err)
	}
	fmt.Println("✓ 登录成功!")
	fmt.Printf("  ID Token: %s...\n", loginResp.IDToken[:50])
	fmt.Printf("  过期时间: %d秒\n", loginResp.ExpiresIn)

	// 3. 验证Token
	fmt.Println("\n=== 测试Token验证 ===")
	claims, err := client.VerifyToken(loginResp.IDToken)
	if err != nil {
		log.Fatalf("❌ Token验证失败: %v", err)
	}
	fmt.Println("✓ Token验证成功!")
	fmt.Printf("  用户名: %v\n", claims["cognito:username"])
	fmt.Printf("  邮箱: %v\n", claims["email"])

	// 4. 获取用户信息
	fmt.Println("\n=== 测试获取用户信息 ===")
	userInfo, err := client.GetUser(email)
	if err != nil {
		log.Fatalf("❌ 获取用户信息失败: %v", err)
	}
	fmt.Println("✓ 获取用户信息成功!")
	fmt.Printf("  用户名: %s\n", userInfo.Username)
	fmt.Printf("  邮箱: %s\n", userInfo.Email)
	fmt.Printf("  状态: %s\n", userInfo.UserStatus)
	fmt.Printf("  启用: %v\n", userInfo.Enabled)

	// 5. 更新用户属性
	fmt.Println("\n=== 测试更新用户属性 ===")
	err = client.UpdateUserAttributes(email, map[string]string{
		"name":         "测试用户",
		"phone_number": "+8613800138000",
	})
	if err != nil {
		log.Printf("❌ 更新用户属性失败: %v", err)
	} else {
		fmt.Println("✓ 用户属性更新成功!")
	}

	// 验证更新
	userInfo, _ = client.GetUser(email)
	fmt.Printf("  姓名: %s\n", userInfo.Attributes["name"])
	fmt.Printf("  电话: %s\n", userInfo.Attributes["phone_number"])

	// 6. 刷新Token
	fmt.Println("\n=== 测试刷新Token ===")
	refreshResp, err := client.RefreshToken(loginResp.RefreshToken)
	if err != nil {
		log.Printf("❌ 刷新Token失败: %v", err)
	} else {
		fmt.Println("✓ Token刷新成功!")
		fmt.Printf("  新的ID Token: %s...\n", refreshResp.IDToken[:50])
	}

	// 7. 清理测试用户
	fmt.Println("\n=== 清理测试数据 ===")
	err = client.DeleteUser(email)
	if err != nil {
		log.Printf("删除用户失败: %v", err)
	} else {
		fmt.Println("✓ 测试用户已删除")
	}

	fmt.Println("\n========================================")
	fmt.Println("✓ SDK测试完成!")
	fmt.Println("========================================")
}
EOF

# 在go.mod中添加本地SDK引用
SDK_PATH="$(cd .. && pwd)"
echo "" >> go.mod
echo "replace github.com/difyz9/cognito-sdk => $SDK_PATH" >> go.mod

# 安装依赖
echo ""
echo "安装依赖..."
go mod tidy

# 运行测试
echo ""
echo "=========================================="
echo "运行SDK测试..."
echo "=========================================="
echo ""

go run main.go

echo ""
echo "测试项目位置: $(pwd)"
echo "你可以查看 main.go 了解如何使用SDK"
