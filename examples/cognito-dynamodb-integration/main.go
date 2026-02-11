package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	cognito "github.com/difyz9/cognito-sdk"
)

// UserProfile 用户配置文件（存储在DynamoDB）
type UserProfile struct {
	UserID      string `dynamodbav:"user_id"`
	CognitoSub  string `dynamodbav:"cognito_sub"`
	Email       string `dynamodbav:"email"`
	Username    string `dynamodbav:"username"`
	FullName    string `dynamodbav:"full_name"`
	PhoneNumber string `dynamodbav:"phone_number"`
	Bio         string `dynamodbav:"bio"`
	AvatarURL   string `dynamodbav:"avatar_url"`
	CreatedAt   string `dynamodbav:"created_at"`
	UpdatedAt   string `dynamodbav:"updated_at"`
	LastLogin   string `dynamodbav:"last_login"`
	Status      string `dynamodbav:"status"`
}

// Activity 用户活动记录
type Activity struct {
	ActivityID  string `dynamodbav:"activity_id"`
	UserID      string `dynamodbav:"user_id"`
	Type        string `dynamodbav:"type"`
	Description string `dynamodbav:"description"`
	Timestamp   string `dynamodbav:"timestamp"`
	IPAddress   string `dynamodbav:"ip_address"`
}

func main() {
	// 初始化客户端
	client, err := cognito.NewClient(cognito.Config{
		UserPoolID: os.Getenv("COGNITO_USER_POOL_ID"),
		ClientID:   os.Getenv("COGNITO_CLIENT_ID"),
		Region:     os.Getenv("AWS_REGION"),
	})
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	dbClient := client.NewDynamoDBClient()
	ctx := context.Background()

	fmt.Println("=== Cognito + DynamoDB 集成示例 ===\n")

	// 场景1: 用户注册流程
	fmt.Println("场景1: 用户注册并创建 Profile")
	fmt.Println("----------------------------------------")
	
	username := "testuser_" + time.Now().Format("20060102150405")
	password := "TestPass123!"
	email := fmt.Sprintf("%s@example.com", username)

	// 1.1 在 Cognito 注册用户
	fmt.Printf("1. 在 Cognito 注册用户: %s...\n", username)
	registerResp, err := client.Register(ctx, username, password, email)
	if err != nil {
		log.Printf("注册失败: %v", err)
	} else {
		fmt.Printf("✓ Cognito 注册成功! UserSub: %s\n", registerResp.UserSub)

		// 1.2 在 DynamoDB 创建用户 Profile
		fmt.Println("2. 在 DynamoDB 创建用户 Profile...")
		profile := UserProfile{
			UserID:      registerResp.UserSub,
			CognitoSub:  registerResp.UserSub,
			Email:       email,
			Username:    username,
			FullName:    "Test User",
			PhoneNumber: "+1234567890",
			Bio:         "This is a test user profile",
			AvatarURL:   "https://example.com/avatar.jpg",
			CreatedAt:   time.Now().Format(time.RFC3339),
			UpdatedAt:   time.Now().Format(time.RFC3339),
			Status:      "active",
		}

		err = dbClient.PutItem(ctx, "UserProfiles", profile)
		if err != nil {
			log.Printf("创建 Profile 失败: %v", err)
		} else {
			fmt.Printf("✓ Profile 创建成功!\n\n")
		}
	}

	// 场景2: 用户登录流程
	fmt.Println("场景2: 用户登录并更新活动记录")
	fmt.Println("----------------------------------------")

	// 2.1 Cognito 登录
	fmt.Println("1. 用户登录...")
	loginResp, err := client.Login(ctx, username, password)
	if err != nil {
		log.Printf("登录失败: %v", err)
	} else {
		fmt.Printf("✓ 登录成功! AccessToken: %s...\n", loginResp.AccessToken[:50])

		// 2.2 从 Token 获取用户信息
		userInfo, err := client.GetUserInfoFromToken(loginResp.AccessToken)
		if err != nil {
			log.Printf("获取用户信息失败: %v", err)
		} else {
			fmt.Printf("✓ 用户 Sub: %s\n", userInfo.Sub)

			// 2.3 更新 DynamoDB 中的最后登录时间
			fmt.Println("2. 更新用户 Profile（最后登录时间）...")
			key := map[string]types.AttributeValue{
				"user_id": &types.AttributeValueMemberS{Value: userInfo.Sub},
			}
			updateExpression := "SET last_login = :last_login, updated_at = :updated_at"
			expressionValues := map[string]types.AttributeValue{
				":last_login":  &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},
				":updated_at": &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},
			}

			err = dbClient.UpdateItem(ctx, "UserProfiles", key, updateExpression, expressionValues, nil)
			if err != nil {
				log.Printf("更新失败: %v", err)
			} else {
				fmt.Printf("✓ 最后登录时间已更新\n")
			}

			// 2.4 记录登录活动
			fmt.Println("3. 记录用户登录活动...")
			activity := Activity{
				ActivityID:  fmt.Sprintf("%s_%d", userInfo.Sub, time.Now().Unix()),
				UserID:      userInfo.Sub,
				Type:        "login",
				Description: "User logged in successfully",
				Timestamp:   time.Now().Format(time.RFC3339),
				IPAddress:   "127.0.0.1",
			}

			err = dbClient.PutItem(ctx, "UserActivities", activity)
			if err != nil {
				log.Printf("记录活动失败: %v", err)
			} else {
				fmt.Printf("✓ 登录活动已记录\n\n")
			}
		}
	}

	// 场景3: 获取用户完整信息
	fmt.Println("场景3: 获取用户完整信息")
	fmt.Println("----------------------------------------")

	if loginResp != nil && loginResp.AccessToken != "" {
		// 3.1 从 Token 获取基本信息
		userInfo, _ := client.GetUserInfoFromToken(loginResp.AccessToken)

		// 3.2 从 DynamoDB 获取详细 Profile
		fmt.Println("1. 从 DynamoDB 获取用户 Profile...")
		key := map[string]types.AttributeValue{
			"user_id": &types.AttributeValueMemberS{Value: userInfo.Sub},
		}

		var profile UserProfile
		err = dbClient.GetItem(ctx, "UserProfiles", key, &profile)
		if err != nil {
			log.Printf("获取 Profile 失败: %v", err)
		} else {
			fmt.Println("✓ 用户完整信息:")
			fmt.Printf("  - UserID: %s\n", profile.UserID)
			fmt.Printf("  - Username: %s\n", profile.Username)
			fmt.Printf("  - Email: %s\n", profile.Email)
			fmt.Printf("  - Full Name: %s\n", profile.FullName)
			fmt.Printf("  - Phone: %s\n", profile.PhoneNumber)
			fmt.Printf("  - Bio: %s\n", profile.Bio)
			fmt.Printf("  - Last Login: %s\n", profile.LastLogin)
			fmt.Printf("  - Status: %s\n\n", profile.Status)
		}
	}

	// 场景4: 查询用户活动历史
	fmt.Println("场景4: 查询用户活动历史")
	fmt.Println("----------------------------------------")

	if loginResp != nil {
		userInfo, _ := client.GetUserInfoFromToken(loginResp.AccessToken)
		
		fmt.Println("1. 查询用户所有活动...")
		var activities []Activity
		err = dbClient.Scan(ctx, "UserActivities", &activities, nil, nil)
		if err != nil {
			log.Printf("查询活动失败: %v", err)
		} else {
			fmt.Printf("✓ 找到 %d 条活动记录:\n", len(activities))
			for _, act := range activities {
				if act.UserID == userInfo.Sub {
					fmt.Printf("  - [%s] %s: %s\n", act.Timestamp, act.Type, act.Description)
				}
			}
			fmt.Println()
		}
	}

	// 场景5: 更新用户 Profile
	fmt.Println("场景5: 更新用户 Profile")
	fmt.Println("----------------------------------------")

	if loginResp != nil {
		userInfo, _ := client.GetUserInfoFromToken(loginResp.AccessToken)
		
		fmt.Println("1. 更新用户简介和头像...")
		key := map[string]types.AttributeValue{
			"user_id": &types.AttributeValueMemberS{Value: userInfo.Sub},
		}
		updateExpression := "SET bio = :bio, avatar_url = :avatar, updated_at = :updated"
		expressionValues := map[string]types.AttributeValue{
			":bio":     &types.AttributeValueMemberS{Value: "Updated bio - I love coding!"},
			":avatar":  &types.AttributeValueMemberS{Value: "https://example.com/new-avatar.jpg"},
			":updated": &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},
		}

		err = dbClient.UpdateItem(ctx, "UserProfiles", key, updateExpression, expressionValues, nil)
		if err != nil {
			log.Printf("更新失败: %v", err)
		} else {
			fmt.Println("✓ Profile 更新成功\n")

			// 记录更新活动
			activity := Activity{
				ActivityID:  fmt.Sprintf("%s_%d", userInfo.Sub, time.Now().Unix()),
				UserID:      userInfo.Sub,
				Type:        "profile_update",
				Description: "User updated profile information",
				Timestamp:   time.Now().Format(time.RFC3339),
				IPAddress:   "127.0.0.1",
			}
			dbClient.PutItem(ctx, "UserActivities", activity)
		}
	}

	// 场景6: 管理员操作 - 列出所有用户
	fmt.Println("场景6: 管理员查看所有用户")
	fmt.Println("----------------------------------------")

	fmt.Println("1. 获取所有用户 Profile...")
	var allProfiles []UserProfile
	err = dbClient.Scan(ctx, "UserProfiles", &allProfiles, nil, nil)
	if err != nil {
		log.Printf("扫描失败: %v", err)
	} else {
		fmt.Printf("✓ 系统中共有 %d 个用户:\n", len(allProfiles))
		for i, p := range allProfiles {
			fmt.Printf("  %d. %s (%s) - 状态: %s\n", i+1, p.Username, p.Email, p.Status)
		}
		fmt.Println()
	}

	fmt.Println("=== 示例完成 ===")
	fmt.Println("\n提示:")
	fmt.Println("- 本示例演示了 Cognito 认证与 DynamoDB 数据存储的完整集成")
	fmt.Println("- 实际应用中，建议为 DynamoDB 表添加适当的索引和分区键")
	fmt.Println("- 考虑使用 DynamoDB Streams 监听数据变化")
	fmt.Println("- 实施适当的错误处理和数据验证")
}
