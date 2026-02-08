package main

import (
	"fmt"
	"log"
	"os"

	cognito "github.com/difyz9/cognito-sdk"
)

func main() {

	client, err := cognito.NewClient(cognito.Config{
		UserPoolID: os.Getenv("COGNITO_USER_POOL_ID"),
		ClientID:   os.Getenv("COGNITO_CLIENT_ID"),
		Region:     os.Getenv("AWS_REGION"),
	})
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	// 删除测试用户
	fmt.Println("删除旧的测试用户...")
	err = client.DeleteUser("user@example.com")
	if err != nil {
		log.Printf("删除用户失败（可能不存在）: %v", err)
	} else {
		fmt.Println("删除成功!")
	}

	// 重新创建用户
	fmt.Println("\n创建新用户...")
	resp, err := client.Register("user@example.com", "Password123!", "")
	if err != nil {
		log.Fatalf("注册失败: %v", err)
	}
	fmt.Printf("注册成功! UserSub: %s\n", resp.UserSub)

	// 尝试登录
	fmt.Println("\n测试登录...")
	loginResp, err := client.Login("user@example.com", "Password123!")
	if err != nil {
		log.Fatalf("登录失败: %v", err)
	}

	fmt.Printf("登录成功!\n")
	fmt.Printf("ID Token: %s...\n", loginResp.IDToken[:50])
	fmt.Printf("Access Token: %s...\n", loginResp.AccessToken[:50])
	fmt.Printf("Refresh Token: %s...\n", loginResp.RefreshToken[:50])
	fmt.Printf("过期时间: %d 秒\n", loginResp.ExpiresIn)
}
