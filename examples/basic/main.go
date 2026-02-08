package main

import (
	"fmt"
	"log"
	"os"

	cognito "github.com/difyz9/cognito-sdk"
)



func main() {

	// 方式1: 使用默认的AWS凭证（从环境变量或~/.aws/credentials）
	client, err := cognito.NewClient(cognito.Config{
		UserPoolID: os.Getenv("COGNITO_USER_POOL_ID"),
		ClientID:   os.Getenv("COGNITO_CLIENT_ID"),
		Region:     os.Getenv("AWS_REGION"),
	})
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	// 方式2: 使用指定的AWS凭证
	// client, err := cognito.NewClientWithCredentials(
	// 	cognito.Config{
	// 		UserPoolID: os.Getenv("COGNITO_USER_POOL_ID"),
	// 		ClientID:   os.Getenv("COGNITO_CLIENT_ID"),
	// 		Region:     os.Getenv("AWS_REGION"),
	// 	},
	// 	os.Getenv("AWS_ACCESS_KEY_ID"),
	// 	os.Getenv("AWS_SECRET_ACCESS_KEY"),
	// )

	// 示例：注册用户
	registerExample(client)

	// 示例：用户登录
	loginExample(client)

	// 示例：Token验证
	tokenExample(client)

	// 示例：用户管理
	userManagementExample(client)

	// 示例：组管理
	groupExample(client)
}

// 注册示例
func registerExample(client *cognito.Client) {
	fmt.Println("=== 注册示例 ===")

	resp, err := client.Register("user@example.com", "Password123!", "")
	if err != nil {
		log.Printf("注册失败: %v", err)
		return
	}

	fmt.Printf("注册成功! UserSub: %s\n", resp.UserSub)
}

// 登录示例
func loginExample(client *cognito.Client) {
	fmt.Println("\n=== 登录示例 ===")

	resp, err := client.Login("user@example.com", "Password123!")
	if err != nil {
		log.Printf("登录失败: %v", err)
		return
	}

	fmt.Printf("登录成功!\n")
	fmt.Printf("ID Token: %s...\n", resp.IDToken[:50])
	fmt.Printf("Access Token: %s...\n", resp.AccessToken[:50])
	fmt.Printf("Refresh Token: %s...\n", resp.RefreshToken[:50])
	fmt.Printf("过期时间: %d 秒\n", resp.ExpiresIn)
}

// Token验证示例
func tokenExample(client *cognito.Client) {
	fmt.Println("\n=== Token验证示例 ===")

	// 先登录获取Token
	loginResp, err := client.Login("user@example.com", "Password123!")
	if err != nil {
		log.Printf("登录失败: %v", err)
		return
	}

	// 验证Token
	claims, err := client.VerifyToken(loginResp.IDToken)
	if err != nil {
		log.Printf("Token验证失败: %v", err)
		return
	}

	fmt.Printf("Token验证成功! Claims: %+v\n", claims)

	// 从Token提取用户信息
	userInfo, err := client.GetUserInfoFromToken(loginResp.IDToken)
	if err != nil {
		log.Printf("提取用户信息失败: %v", err)
		return
	}

	fmt.Printf("用户信息: %+v\n", userInfo)

	// 刷新Token
	refreshResp, err := client.RefreshToken(loginResp.RefreshToken)
	if err != nil {
		log.Printf("刷新Token失败: %v", err)
		return
	}

	fmt.Printf("Token刷新成功! 新的ID Token: %s...\n", refreshResp.IDToken[:50])
}

// 用户管理示例
func userManagementExample(client *cognito.Client) {
	fmt.Println("\n=== 用户管理示例 ===")

	username := "user@example.com"

	// 获取用户信息
	userInfo, err := client.GetUser(username)
	if err != nil {
		log.Printf("获取用户信息失败: %v", err)
		return
	}

	fmt.Printf("用户信息: %+v\n", userInfo)

	// 更新用户属性
	err = client.UpdateUserAttributes(username, map[string]string{
		"name":         "张三",
		"phone_number": "+8613800138000",
	})
	if err != nil {
		log.Printf("更新用户属性失败: %v", err)
		return
	}

	fmt.Println("用户属性更新成功!")

	// 修改密码
	err = client.ChangePassword(username, "Password123!", "NewPassword123!")
	if err != nil {
		log.Printf("修改密码失败: %v", err)
		return
	}

	fmt.Println("密码修改成功!")

	// 列出所有用户
	listResp, err := client.ListUsers(cognito.ListUsersFilter{
		Limit: 10,
	})
	if err != nil {
		log.Printf("获取用户列表失败: %v", err)
		return
	}

	fmt.Printf("用户列表 (共%d个):\n", len(listResp.Users))
	for _, user := range listResp.Users {
		fmt.Printf("  - %s (%s)\n", user.Username, user.Email)
	}

	// 禁用用户
	err = client.DisableUser(username)
	if err != nil {
		log.Printf("禁用用户失败: %v", err)
		return
	}

	fmt.Println("用户已禁用!")

	// 启用用户
	err = client.EnableUser(username)
	if err != nil {
		log.Printf("启用用户失败: %v", err)
		return
	}

	fmt.Println("用户已启用!")
}

// 组管理示例
func groupExample(client *cognito.Client) {
	fmt.Println("\n=== 组管理示例 ===")

	username := "user@example.com"

	// 添加用户到组
	err := client.AddUserToGroup(username, "Admins")
	if err != nil {
		log.Printf("添加用户到组失败: %v", err)
		return
	}

	fmt.Println("用户已添加到Admins组!")

	// 列出用户所属的组
	groups, err := client.ListUserGroups(username)
	if err != nil {
		log.Printf("获取用户组列表失败: %v", err)
		return
	}

	fmt.Printf("用户所属的组:\n")
	for _, group := range groups {
		fmt.Printf("  - %s: %s\n", group.GroupName, group.Description)
	}

	// 从组中移除用户
	err = client.RemoveUserFromGroup(username, "Admins")
	if err != nil {
		log.Printf("从组中移除用户失败: %v", err)
		return
	}

	fmt.Println("用户已从Admins组移除!")
}
