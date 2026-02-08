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
	fmt.Println("AWS Cognito 注册方式演示")
	fmt.Println("===========================================\n")

	// 演示方式1：管理员创建用户（无需验证）
	// demoAdminRegister(client)

	// fmt.Println("\n" + strings.Repeat("=", 50) + "\n")

	// 演示方式2：用户自注册（需要邮箱验证）
	demoUserSignUp(client)
}

// 方式1: 管理员创建用户（Register）
func demoAdminRegister(client *cognito.Client) {
	fmt.Println("【方式1】管理员创建用户 - Register")
	fmt.Println("特点：立即可用，无需邮箱验证")
	fmt.Println(strings.Repeat("-", 50))

	email := "admin-created@example.com"
	password := "Admin@123456"

	// 清理旧用户
	fmt.Printf("\n1. 清理旧用户 %s...\n", email)
	err := client.DeleteUser(email)
	if err != nil {
		fmt.Printf("   删除用户失败（可能不存在）: %v\n", err)
	} else {
		fmt.Println("   ✓ 旧用户已删除")
	}

	// 管理员创建用户
	fmt.Printf("\n2. 使用 Register 创建用户...\n")
	resp, err := client.Register(email, password, "")
	if err != nil {
		log.Printf("   ✗ 注册失败: %v\n", err)
		return
	}
	fmt.Printf("   ✓ 用户创建成功!\n")
	fmt.Printf("   UserSub: %s\n", resp.UserSub)

	// 获取用户信息
	fmt.Printf("\n3. 查看用户状态...\n")
	userInfo, err := client.GetUser(email)
	if err != nil {
		log.Printf("   ✗ 获取用户失败: %v\n", err)
		return
	}
	fmt.Printf("   用户状态: %s\n", userInfo.UserStatus)
	fmt.Printf("   邮箱验证: %v\n", userInfo.EmailVerified)
	fmt.Printf("   是否启用: %v\n", userInfo.Enabled)

	// 立即登录测试
	fmt.Printf("\n4. 测试立即登录（无需验证）...\n")
	loginResp, err := client.Login(email, password)
	if err != nil {
		log.Printf("   ✗ 登录失败: %v\n", err)
		return
	}
	fmt.Printf("   ✓ 登录成功!\n")
	fmt.Printf("   Token (前50字符): %s...\n", loginResp.IDToken[:50])

	fmt.Println("\n✅ 管理员创建方式完成: 用户创建后可立即登录，无需任何验证步骤")
}

// 方式2: 用户自注册（SignUpUser）
func demoUserSignUp(client *cognito.Client) {
	fmt.Println("【方式2】用户自注册 - SignUpUser")
	fmt.Println("特点：需要邮箱验证码确认")
	fmt.Println(strings.Repeat("-", 50))

	email := "admin@126.com"
	password := "Ab@123456"

	// 清理旧用户
	fmt.Printf("\n1. 清理旧用户 %s...\n", email)
	err := client.DeleteUser(email)
	if err != nil {
		fmt.Printf("   删除用户失败（可能不存在）: %v\n", err)
	} else {
		fmt.Println("   ✓ 旧用户已删除")
	}

	// 用户自注册
	fmt.Printf("\n2. 使用 SignUpUser 进行自注册...\n")
	resp, err := client.SignUpUser(email, password, "")
	if err != nil {
		log.Printf("   ✗ 注册失败: %v\n", err)
		return
	}
	fmt.Printf("   ✓ 注册请求成功!\n")
	fmt.Printf("   UserSub: %s\n", resp.UserSub)
	fmt.Printf("   用户已确认: %v\n", resp.UserConfirmed)

	if resp.CodeDeliveryDetails != nil {
		fmt.Printf("\n   📧 验证码发送信息:\n")
		fmt.Printf("      目标地址: %s\n", resp.CodeDeliveryDetails.Destination)
		fmt.Printf("      发送方式: %s\n", resp.CodeDeliveryDetails.DeliveryMedium)
		fmt.Printf("      属性名称: %s\n", resp.CodeDeliveryDetails.AttributeName)
	}

	// 获取用户信息
	fmt.Printf("\n3. 查看未确认用户的状态...\n")
	userInfo, err := client.GetUser(email)
	if err != nil {
		log.Printf("   ✗ 获取用户失败: %v\n", err)
		return
	}
	fmt.Printf("   用户状态: %s ⚠️\n", userInfo.UserStatus)
	fmt.Printf("   邮箱验证: %v\n", userInfo.EmailVerified)
	fmt.Printf("   是否启用: %v\n", userInfo.Enabled)

	// 尝试未验证前登录
	fmt.Printf("\n4. 尝试在验证前登录...\n")
	_, err = client.Login(email, password)
	if err != nil {
		fmt.Printf("   ✗ 登录失败（符合预期）: %v\n", err)
	} else {
		fmt.Printf("   ✓ 登录成功（意外）\n")
	}

	// 模拟用户输入验证码
	fmt.Printf("\n5. 模拟验证流程...\n")
	fmt.Printf("   在实际应用中，用户需要:\n")
	fmt.Printf("   - 检查邮箱收件箱\n")
	fmt.Printf("   - 获取验证码（6位数字）\n")
	fmt.Printf("   - 调用 ConfirmSignUp 完成验证\n")

	// 由于我们没有实际邮件，使用管理员确认
	// fmt.Printf("\n6. 使用管理员权限确认用户（模拟验证）...\n")
	// err = client.ConfirmUser(email)
	// if err != nil {
	// 	log.Printf("   ✗ 确认失败: %v\n", err)
	// 	return
	// }
	// fmt.Printf("   ✓ 用户已确认\n")

	// 确认后查看状态
	fmt.Printf("\n7. 查看确认后的用户状态...\n")
	userInfo, err = client.GetUser(email)
	if err != nil {
		log.Printf("   ✗ 获取用户失败: %v\n", err)
		return
	}
	fmt.Printf("   用户状态: %s ✅\n", userInfo.UserStatus)
	fmt.Printf("   邮箱验证: %v\n", userInfo.EmailVerified)
	fmt.Printf("   是否启用: %v\n", userInfo.Enabled)

	// 确认后登录
	fmt.Printf("\n8. 确认后尝试登录...\n")
	loginResp, err := client.Login(email, password)
	if err != nil {
		log.Printf("   ✗ 登录失败: %v\n", err)
		return
	}
	fmt.Printf("   ✓ 登录成功!\n")
	fmt.Printf("   Token (前50字符): %s...\n", loginResp.IDToken[:50])

	fmt.Println("\n✅ 用户自注册方式完成: 需要验证步骤，但更适合公开注册场景")
}

// 交互式演示（可选）
func interactiveDemo(client *cognito.Client) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n=== 交互式注册演示 ===")
	fmt.Print("请输入邮箱: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	fmt.Print("请输入密码: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	fmt.Println("\n选择注册方式:")
	fmt.Println("1. 管理员创建（立即可用）")
	fmt.Println("2. 用户自注册（需要验证）")
	fmt.Print("请选择 (1/2): ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		resp, err := client.Register(email, password, "")
		if err != nil {
			log.Fatalf("注册失败: %v", err)
		}
		fmt.Printf("注册成功! UserSub: %s\n", resp.UserSub)
		fmt.Println("用户可以立即登录!")

	case "2":
		resp, err := client.SignUpUser(email, password, "")
		if err != nil {
			log.Fatalf("注册失败: %v", err)
		}
		fmt.Printf("注册成功! UserSub: %s\n", resp.UserSub)
		if resp.CodeDeliveryDetails != nil {
			fmt.Printf("验证码已发送至: %s\n", resp.CodeDeliveryDetails.Destination)
		}

		fmt.Print("\n请输入收到的验证码: ")
		code, _ := reader.ReadString('\n')
		code = strings.TrimSpace(code)

		err = client.ConfirmSignUp(email, code)
		if err != nil {
			log.Fatalf("验证失败: %v", err)
		}
		fmt.Println("验证成功! 现在可以登录了!")

	default:
		fmt.Println("无效选择")
	}
}
