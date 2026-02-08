package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	cognito "github.com/difyz9/cognito-sdk"
)

func main() {
	// 配置AWS Cognito

	client, err := cognito.NewClient(cognito.Config{
		UserPoolID: os.Getenv("COGNITO_USER_POOL_ID"),
		ClientID:   os.Getenv("COGNITO_CLIENT_ID"),
		Region:     os.Getenv("AWS_REGION"),
	})
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	fmt.Println("===========================================")
	fmt.Println("AWS Cognito 邮箱验证码激活用户示例")
	fmt.Println("===========================================\n")

	reader := bufio.NewReader(os.Stdin)

	// 步骤1: 输入真实邮箱地址
	fmt.Print("请输入您的真实邮箱地址: ")
	emailInput, _ := reader.ReadString('\n')
	email := strings.TrimSpace(emailInput)

	if email == "" {
		log.Fatal("邮箱地址不能为空")
	}

	// 步骤2: 输入密码
	fmt.Print("请设置登录密码（至少8位，包含大小写字母、数字和特殊字符）: ")
	passwordInput, _ := reader.ReadString('\n')
	password := strings.TrimSpace(passwordInput)

	if password == "" {
		log.Fatal("密码不能为空")
	}

	// 步骤3: 清理旧用户（如果存在）
	fmt.Printf("\n正在检查用户 %s 是否已存在...\n", email)
	err = client.DeleteUser(email)
	if err != nil {
		fmt.Printf("删除用户失败（可能不存在）: %v\n", err)
	} else {
		fmt.Println("✓ 旧用户已删除")
	}

	// 步骤4: 注册新用户
	fmt.Printf("\n正在注册新用户...\n")
	resp, err := client.SignUpUser(email, password, "")
	if err != nil {
		log.Fatalf("注册失败: %v", err)
	}

	fmt.Printf("✓ 注册成功!\n")
	fmt.Printf("UserSub: %s\n", resp.UserSub)
	fmt.Printf("用户已确认: %v\n", resp.UserConfirmed)

	if resp.CodeDeliveryDetails != nil {
		fmt.Printf("\n📧 验证码已发送:\n")
		fmt.Printf("   目标地址: %s\n", resp.CodeDeliveryDetails.Destination)
		fmt.Printf("   发送方式: %s\n", resp.CodeDeliveryDetails.DeliveryMedium)
		fmt.Printf("   属性名称: %s\n", resp.CodeDeliveryDetails.AttributeName)
	}

	// 步骤5: 查看未确认状态
	fmt.Printf("\n查看用户当前状态...\n")
	userInfo, err := client.GetUser(email)
	if err != nil {
		log.Printf("获取用户失败: %v\n", err)
	} else {
		fmt.Printf("   用户状态: %s\n", userInfo.UserStatus)
		fmt.Printf("   邮箱验证: %v\n", userInfo.EmailVerified)
		fmt.Printf("   是否启用: %v\n", userInfo.Enabled)
	}

	// 步骤6: 提示用户查收邮件并输入验证码
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("请查收您的邮箱，AWS Cognito 已发送验证码")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Print("\n请输入您收到的6位验证码: ")
	codeInput, _ := reader.ReadString('\n')
	verificationCode := strings.TrimSpace(codeInput)

	if verificationCode == "" {
		log.Fatal("验证码不能为空")
	}

	// 步骤7: 使用验证码确认注册
	fmt.Printf("\n正在验证邮箱...\n")
	err = client.ConfirmSignUp(email, verificationCode)
	if err != nil {
		log.Fatalf("验证失败: %v", err)
	}

	fmt.Printf("✓ 邮箱验证成功!\n")

	// 步骤8: 查看确认后的状态
	fmt.Printf("\n查看验证后的用户状态...\n")
	userInfo, err = client.GetUser(email)
	if err != nil {
		log.Printf("获取用户失败: %v\n", err)
	} else {
		fmt.Printf("   用户状态: %s ✅\n", userInfo.UserStatus)
		fmt.Printf("   邮箱验证: %v ✅\n", userInfo.EmailVerified)
		fmt.Printf("   是否启用: %v\n", userInfo.Enabled)
	}

	// 步骤9: 测试登录
	fmt.Printf("\n正在测试登录...\n")
	loginResp, err := client.Login(email, password)
	if err != nil {
		log.Fatalf("登录失败: %v", err)
	}

	fmt.Printf("✓ 登录成功!\n")
	fmt.Printf("\n认证信息:\n")
	fmt.Printf("   AccessToken (前50字符): %s...\n", loginResp.AccessToken[:50])
	fmt.Printf("   IDToken (前50字符): %s...\n", loginResp.IDToken[:50])
	fmt.Printf("   RefreshToken (前50字符): %s...\n", loginResp.RefreshToken[:50])
	fmt.Printf("   过期时间: %d 秒\n", loginResp.ExpiresIn)

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("✅ 完整流程完成!")
	fmt.Println("用户已成功注册、验证并登录")
	fmt.Println(strings.Repeat("=", 50))
}
