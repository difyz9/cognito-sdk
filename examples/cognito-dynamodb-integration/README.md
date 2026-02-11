# Cognito + DynamoDB 集成示例

这个示例展示了如何将 AWS Cognito 用户认证与 DynamoDB 数据存储结合使用，构建一个完整的用户管理系统。

## 应用场景

这是一个真实世界的用户管理系统，包含：

1. **用户注册** - Cognito 认证 + DynamoDB Profile 创建
2. **用户登录** - Cognito 验证 + DynamoDB 活动记录
3. **用户信息管理** - Cognito 基础信息 + DynamoDB 扩展信息
4. **活动追踪** - 记录用户行为日志
5. **Profile 更新** - 同步更新用户数据
6. **管理员功能** - 查看和管理所有用户

## 数据架构

### Cognito 存储
- 用户名和密码
- 邮箱验证状态
- 用户基本属性
- 认证 Token

### DynamoDB 存储

#### UserProfiles 表
```
user_id (PK)    - Cognito UserSub
cognito_sub     - Cognito 用户标识
email           - 邮箱
username        - 用户名
full_name       - 全名
phone_number    - 电话号码
bio             - 个人简介
avatar_url      - 头像URL
created_at      - 创建时间
updated_at      - 更新时间
last_login      - 最后登录时间
status          - 状态 (active/inactive/banned)
```

#### UserActivities 表
```
activity_id (PK) - 活动ID
user_id         - 用户ID
type            - 活动类型 (login/logout/profile_update等)
description     - 描述
timestamp       - 时间戳
ip_address      - IP地址
```

## 前置要求

### 1. 创建 DynamoDB 表

#### UserProfiles 表
```bash
aws dynamodb create-table \
    --table-name UserProfiles \
    --attribute-definitions \
        AttributeName=user_id,AttributeType=S \
    --key-schema \
        AttributeName=user_id,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST \
    --region us-east-1
```

#### UserActivities 表
```bash
aws dynamodb create-table \
    --table-name UserActivities \
    --attribute-definitions \
        AttributeName=activity_id,AttributeType=S \
        AttributeName=user_id,AttributeType=S \
        AttributeName=timestamp,AttributeType=S \
    --key-schema \
        AttributeName=activity_id,KeyType=HASH \
    --global-secondary-indexes \
        '[{
            "IndexName": "UserIdIndex",
            "KeySchema": [
                {"AttributeName":"user_id","KeyType":"HASH"},
                {"AttributeName":"timestamp","KeyType":"RANGE"}
            ],
            "Projection": {"ProjectionType":"ALL"}
        }]' \
    --billing-mode PAY_PER_REQUEST \
    --region us-east-1
```

### 2. 配置环境变量

```bash
export AWS_REGION="us-east-1"
export COGNITO_USER_POOL_ID="us-east-1_xxxxx"
export COGNITO_CLIENT_ID="xxxxxxxxxxxxxxxxxx"
export AWS_ACCESS_KEY_ID="your-access-key"
export AWS_SECRET_ACCESS_KEY="your-secret-key"
```

## 运行示例

```bash
cd examples/cognito-dynamodb-integration
go run main.go
```

## 示例输出

```
=== Cognito + DynamoDB 集成示例 ===

场景1: 用户注册并创建 Profile
----------------------------------------
1. 在 Cognito 注册用户: testuser_20260211123456...
✓ Cognito 注册成功! UserSub: abc123-def456-ghi789
2. 在 DynamoDB 创建用户 Profile...
✓ Profile 创建成功!

场景2: 用户登录并更新活动记录
----------------------------------------
1. 用户登录...
✓ 登录成功! AccessToken: eyJraWQiOiJxxxxxx...
✓ 用户 Sub: abc123-def456-ghi789
2. 更新用户 Profile（最后登录时间）...
✓ 最后登录时间已更新
3. 记录用户登录活动...
✓ 登录活动已记录

场景3: 获取用户完整信息
----------------------------------------
1. 从 DynamoDB 获取用户 Profile...
✓ 用户完整信息:
  - UserID: abc123-def456-ghi789
  - Username: testuser_20260211123456
  - Email: testuser_20260211123456@example.com
  - Full Name: Test User
  - Phone: +1234567890
  - Bio: This is a test user profile
  - Last Login: 2026-02-11T12:35:00Z
  - Status: active

场景4: 查询用户活动历史
----------------------------------------
1. 查询用户所有活动...
✓ 找到 2 条活动记录:
  - [2026-02-11T12:35:00Z] login: User logged in successfully
  - [2026-02-11T12:35:05Z] profile_update: User updated profile information

场景5: 更新用户 Profile
----------------------------------------
1. 更新用户简介和头像...
✓ Profile 更新成功

场景6: 管理员查看所有用户
----------------------------------------
1. 获取所有用户 Profile...
✓ 系统中共有 3 个用户:
  1. testuser_20260211123456 (testuser_20260211123456@example.com) - 状态: active
  2. john_doe (john@example.com) - 状态: active
  3. jane_smith (jane@example.com) - 状态: active

=== 示例完成 ===
```

## 核心实现

### 用户注册流程

```go
// 1. Cognito 注册
registerResp, err := client.Register(ctx, username, password, email)

// 2. DynamoDB 创建 Profile
profile := UserProfile{
    UserID:     registerResp.UserSub,
    CognitoSub: registerResp.UserSub,
    Email:      email,
    Username:   username,
    CreatedAt:  time.Now().Format(time.RFC3339),
    Status:     "active",
}
err = dbClient.PutItem(ctx, "UserProfiles", profile)
```

### 用户登录流程

```go
// 1. Cognito 登录
loginResp, err := client.Login(ctx, username, password)

// 2. 获取用户信息
userInfo, err := client.GetUserInfoFromToken(loginResp.AccessToken)

// 3. 更新最后登录时间
updateExpression := "SET last_login = :last_login"
expressionValues := map[string]types.AttributeValue{
    ":last_login": &types.AttributeValueMemberS{
        Value: time.Now().Format(time.RFC3339),
    },
}
err = dbClient.UpdateItem(ctx, "UserProfiles", key, 
    updateExpression, expressionValues, nil)

// 4. 记录登录活动
activity := Activity{
    ActivityID: fmt.Sprintf("%s_%d", userInfo.Sub, time.Now().Unix()),
    UserID:     userInfo.Sub,
    Type:       "login",
    Timestamp:  time.Now().Format(time.RFC3339),
}
err = dbClient.PutItem(ctx, "UserActivities", activity)
```

## 最佳实践

### 1. 数据一致性

- 使用事务确保 Cognito 和 DynamoDB 操作的一致性
- 实施重试机制处理临时失败
- 考虑使用 DynamoDB Transactions 进行复杂操作

### 2. 性能优化

- 为频繁查询的字段创建 GSI（全局二级索引）
- 使用 BatchGetItem 批量获取数据
- 实施缓存策略减少数据库访问

### 3. 安全考虑

- 永远不要在 DynamoDB 存储密码
- 使用 Token 验证所有 API 调用
- 实施细粒度的 IAM 权限
- 记录所有敏感操作

### 4. 扩展建议

```go
// 添加用户偏好设置表
type UserPreferences struct {
    UserID         string `dynamodbav:"user_id"`
    Theme          string `dynamodbav:"theme"`
    Language       string `dynamodbav:"language"`
    Notifications  bool   `dynamodbav:"notifications"`
    EmailFrequency string `dynamodbav:"email_frequency"`
}

// 添加用户会话表
type UserSession struct {
    SessionID    string `dynamodbav:"session_id"`
    UserID       string `dynamodbav:"user_id"`
    AccessToken  string `dynamodbav:"access_token"`
    RefreshToken string `dynamodbav:"refresh_token"`
    CreatedAt    string `dynamodbav:"created_at"`
    ExpiresAt    string `dynamodbav:"expires_at"`
}
```

## 故障排除

### Cognito 用户已存在但 DynamoDB 中没有 Profile

这可能发生在注册流程中断时。解决方案：

```go
// 检查并创建缺失的 Profile
func ensureProfileExists(ctx context.Context, client *cognito.Client, 
    dbClient *cognito.DynamoDBClient, userSub string) error {
    
    key := map[string]types.AttributeValue{
        "user_id": &types.AttributeValueMemberS{Value: userSub},
    }
    
    var profile UserProfile
    err := dbClient.GetItem(ctx, "UserProfiles", key, &profile)
    if err != nil {
        // Profile 不存在，创建它
        // ... 创建逻辑
    }
    return nil
}
```

### DynamoDB 表权限错误

确保 IAM 角色/用户具有必要权限：

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "dynamodb:PutItem",
                "dynamodb:GetItem",
                "dynamodb:UpdateItem",
                "dynamodb:DeleteItem",
                "dynamodb:Query",
                "dynamodb:Scan"
            ],
            "Resource": [
                "arn:aws:dynamodb:us-east-1:*:table/UserProfiles",
                "arn:aws:dynamodb:us-east-1:*:table/UserActivities"
            ]
        }
    ]
}
```

## 清理资源

```bash
# 删除 DynamoDB 表
aws dynamodb delete-table --table-name UserProfiles --region us-east-1
aws dynamodb delete-table --table-name UserActivities --region us-east-1

# 删除 Cognito 用户（如果需要）
aws cognito-idp admin-delete-user \
    --user-pool-id us-east-1_xxxxx \
    --username testuser_20260211123456
```

## 相关文档

- [AWS Cognito 文档](https://docs.aws.amazon.com/cognito/)
- [DynamoDB 最佳实践](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/best-practices.html)
- [SDK API 参考](../../API_REFERENCE.md)
