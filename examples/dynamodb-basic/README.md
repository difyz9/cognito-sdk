# DynamoDB 基础操作示例

这个示例演示如何使用 Cognito SDK 集成的 DynamoDB 客户端进行基础操作。

## 功能演示

1. ✅ 检查表是否存在
2. ✅ 插入单个项目 (PutItem)
3. ✅ 获取项目 (GetItem)
4. ✅ 更新项目 (UpdateItem)
5. ✅ 批量插入 (BatchWriteItem)
6. ✅ 扫描表 (Scan)
7. ✅ 条件查询
8. ✅ 删除项目 (DeleteItem)
9. ✅ 列出所有表 (ListTables)

## 前置要求

### 1. 环境变量配置

```bash
export AWS_REGION="us-east-1"
export COGNITO_USER_POOL_ID="your-user-pool-id"
export COGNITO_CLIENT_ID="your-client-id"

# AWS 凭证（如果不使用默认配置）
export AWS_ACCESS_KEY_ID="your-access-key"
export AWS_SECRET_ACCESS_KEY="your-secret-key"
```

### 2. 创建 DynamoDB 表

使用 AWS CLI 创建测试表：

```bash
aws dynamodb create-table \
    --table-name Users \
    --attribute-definitions \
        AttributeName=user_id,AttributeType=S \
    --key-schema \
        AttributeName=user_id,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST \
    --region us-east-1
```

或使用 AWS Console 创建表：
- 表名: `Users`
- 分区键: `user_id` (String)
- 计费模式: 按需 (PAY_PER_REQUEST)

## 运行示例

```bash
# 进入示例目录
cd examples/dynamodb-basic

# 运行示例
go run main.go
```

## 预期输出

```
=== DynamoDB 基础操作示例 ===

1. 检查表是否存在...
表 'Users' 存在: true

2. 插入新用户...
✓ 成功插入用户: john_doe

3. 获取用户信息...
✓ 用户信息: ID=user-001, 用户名=john_doe, 邮箱=john@example.com

4. 更新用户邮箱...
✓ 成功更新用户邮箱

5. 批量插入用户...
✓ 成功批量插入 2 个用户

6. 扫描所有用户...
✓ 找到 3 个用户:
  - john_doe (john.doe@example.com)
  - jane_smith (jane@example.com)
  - bob_wilson (bob@example.com)

7. 条件查询（包含特定域名的邮箱）...
✓ 找到 3 个匹配用户

8. 删除用户...
✓ 成功删除用户 user-003

9. 列出所有 DynamoDB 表...
✓ 找到 1 个表:
  - Users

=== 示例完成 ===
```

## 核心代码说明

### 1. 初始化客户端

```go
// 创建 Cognito 客户端
client, err := cognito.NewClient(cognito.Config{
    UserPoolID: os.Getenv("COGNITO_USER_POOL_ID"),
    ClientID:   os.Getenv("COGNITO_CLIENT_ID"),
    Region:     os.Getenv("AWS_REGION"),
})

// 创建 DynamoDB 客户端（复用 AWS 配置）
dbClient := client.NewDynamoDBClient()
```

### 2. 定义数据结构

使用 `dynamodbav` 标签定义字段映射：

```go
type User struct {
    UserID    string `dynamodbav:"user_id"`
    Username  string `dynamodbav:"username"`
    Email     string `dynamodbav:"email"`
    CreatedAt string `dynamodbav:"created_at"`
}
```

### 3. 基础操作

#### 插入数据
```go
user := User{
    UserID:    "user-001",
    Username:  "john_doe",
    Email:     "john@example.com",
    CreatedAt: time.Now().Format(time.RFC3339),
}
err = dbClient.PutItem(ctx, "Users", user)
```

#### 获取数据
```go
key := map[string]types.AttributeValue{
    "user_id": &types.AttributeValueMemberS{Value: "user-001"},
}
var user User
err = dbClient.GetItem(ctx, "Users", key, &user)
```

#### 更新数据
```go
updateExpression := "SET email = :email"
expressionValues := map[string]types.AttributeValue{
    ":email": &types.AttributeValueMemberS{Value: "new@example.com"},
}
err = dbClient.UpdateItem(ctx, "Users", key, updateExpression, expressionValues, nil)
```

## 清理资源

运行后如需清理测试数据：

```bash
# 删除表
aws dynamodb delete-table --table-name Users --region us-east-1
```

## 常见问题

### 表不存在错误

确保已创建 `Users` 表，或修改代码中的 `tableName` 变量。

### 权限错误

确保 AWS 凭证具有以下权限：
- `dynamodb:PutItem`
- `dynamodb:GetItem`
- `dynamodb:UpdateItem`
- `dynamodb:DeleteItem`
- `dynamodb:Scan`
- `dynamodb:Query`
- `dynamodb:BatchWriteItem`
- `dynamodb:ListTables`
- `dynamodb:DescribeTable`

### 区域配置

确保 `AWS_REGION` 与创建的表所在区域一致。

## 进阶用法

查看其他示例：
- `../dynamodb-advanced/` - 高级查询和索引
- `../cognito-dynamodb-integration/` - Cognito + DynamoDB 集成场景
