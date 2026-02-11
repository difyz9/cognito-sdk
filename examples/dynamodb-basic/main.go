package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	cognito "github.com/difyz9/cognito-sdk"
)

// User 用户结构体示例
type User struct {
	UserID    string `dynamodbav:"user_id"`
	Username  string `dynamodbav:"username"`
	Email     string `dynamodbav:"email"`
	CreatedAt string `dynamodbav:"created_at"`
}

// Product 产品结构体示例
type Product struct {
	ProductID   string  `dynamodbav:"product_id"`
	Name        string  `dynamodbav:"name"`
	Description string  `dynamodbav:"description"`
	Price       float64 `dynamodbav:"price"`
	Stock       int     `dynamodbav:"stock"`
}

func main() {



// aws dynamodb create-table \
//     --table-name Users \
//     --attribute-definitions \
//         AttributeName=user_id,AttributeType=S \
//     --key-schema \
//         AttributeName=user_id,KeyType=HASH \
//     --billing-mode PAY_PER_REQUEST \
//     --region ap-southeast-1

// aws configure list

	// 初始化Cognito客户端
	client, err := cognito.NewClient(cognito.Config{
		UserPoolID: os.Getenv("COGNITO_USER_POOL_ID"),
		ClientID:   os.Getenv("COGNITO_CLIENT_ID"),
		Region:     os.Getenv("AWS_REGION"),
	})
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	// 创建DynamoDB客户端
	dbClient := client.NewDynamoDBClient()
	ctx := context.Background()

	// 表名
	tableName := "Users"

	fmt.Println("=== DynamoDB 基础操作示例 ===\n")

	// 1. 创建表 (可选)
	fmt.Println("1. 检查表是否存在...")
	exists, err := dbClient.TableExists(ctx, tableName)
	if err != nil {
		log.Printf("检查表失败: %v", err)
	}
	fmt.Printf("表 '%s' 存在: %v\n\n", tableName, exists)

	// 2. 插入数据
	fmt.Println("2. 插入新用户...")
	user := User{
		UserID:    "user-001",
		Username:  "john_doe",
		Email:     "john@example.com",
		CreatedAt: "2026-02-11T10:00:00Z",
	}

	err = dbClient.PutItem(ctx, tableName, user)
	if err != nil {
		log.Printf("插入用户失败: %v", err)
	} else {
		fmt.Printf("✓ 成功插入用户: %s\n\n", user.Username)
	}

	// 3. 获取数据
	fmt.Println("3. 获取用户信息...")
	key := map[string]types.AttributeValue{
		"user_id": &types.AttributeValueMemberS{Value: "user-001"},
	}

	var retrievedUser User
	err = dbClient.GetItem(ctx, tableName, key, &retrievedUser)
	if err != nil {
		log.Printf("获取用户失败: %v", err)
	} else {
		fmt.Printf("✓ 用户信息: ID=%s, 用户名=%s, 邮箱=%s\n\n", 
			retrievedUser.UserID, retrievedUser.Username, retrievedUser.Email)
	}

	// 4. 更新数据
	fmt.Println("4. 更新用户邮箱...")
	updateExpression := "SET email = :email"
	expressionAttributeValues := map[string]types.AttributeValue{
		":email": &types.AttributeValueMemberS{Value: "john.doe@example.com"},
	}

	err = dbClient.UpdateItem(ctx, tableName, key, updateExpression, 
		expressionAttributeValues, nil)
	if err != nil {
		log.Printf("更新用户失败: %v", err)
	} else {
		fmt.Println("✓ 成功更新用户邮箱\n")
	}

	// 5. 批量插入
	fmt.Println("5. 批量插入用户...")
	users := []interface{}{
		User{
			UserID:    "user-002",
			Username:  "jane_smith",
			Email:     "jane@example.com",
			CreatedAt: "2026-02-11T10:30:00Z",
		},
		User{
			UserID:    "user-003",
			Username:  "bob_wilson",
			Email:     "bob@example.com",
			CreatedAt: "2026-02-11T11:00:00Z",
		},
	}

	err = dbClient.BatchWriteItem(ctx, tableName, users)
	if err != nil {
		log.Printf("批量插入失败: %v", err)
	} else {
		fmt.Printf("✓ 成功批量插入 %d 个用户\n\n", len(users))
	}

	// 6. 扫描表
	fmt.Println("6. 扫描所有用户...")
	var allUsers []User
	err = dbClient.Scan(ctx, tableName, &allUsers, nil, nil)
	if err != nil {
		log.Printf("扫描表失败: %v", err)
	} else {
		fmt.Printf("✓ 找到 %d 个用户:\n", len(allUsers))
		for _, u := range allUsers {
			fmt.Printf("  - %s (%s)\n", u.Username, u.Email)
		}
		fmt.Println()
	}

	// 7. 条件扫描
	fmt.Println("7. 条件查询（包含特定域名的邮箱）...")
	filterExpression := aws.String("contains(email, :domain)")
	filterValues := map[string]types.AttributeValue{
		":domain": &types.AttributeValueMemberS{Value: "example.com"},
	}

	var filteredUsers []User
	err = dbClient.Scan(ctx, tableName, &filteredUsers, filterExpression, filterValues)
	if err != nil {
		log.Printf("条件查询失败: %v", err)
	} else {
		fmt.Printf("✓ 找到 %d 个匹配用户\n\n", len(filteredUsers))
	}

	// 8. 删除数据
	fmt.Println("8. 删除用户...")
	deleteKey := map[string]types.AttributeValue{
		"user_id": &types.AttributeValueMemberS{Value: "user-003"},
	}

	err = dbClient.DeleteItem(ctx, tableName, deleteKey)
	if err != nil {
		log.Printf("删除用户失败: %v", err)
	} else {
		fmt.Println("✓ 成功删除用户 user-003\n")
	}

	// 9. 列出所有表
	fmt.Println("9. 列出所有 DynamoDB 表...")
	tables, err := dbClient.ListTables(ctx)
	if err != nil {
		log.Printf("列出表失败: %v", err)
	} else {
		fmt.Printf("✓ 找到 %d 个表:\n", len(tables))
		for _, t := range tables {
			fmt.Printf("  - %s\n", t)
		}
	}

	fmt.Println("\n=== 示例完成 ===")
}
