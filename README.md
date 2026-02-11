# AWS Cognito Go SDK

一个简单易用的 AWS Cognito Go SDK，提供完整的用户认证和管理功能。

> **🔧 最新更新 (2026-02-08)**  
> - ✅ 修复了用户注册后无法登录的问题 → [修复说明](FIX_SUMMARY.md)
> - 🎉 **新功能**: 现在支持两种注册方式 → [使用说明](DUAL_REGISTRATION.md)
>   - **Register** - 管理员创建用户（立即可用）
>   - **SignUpUser** - 用户自注册（邮箱验证流程）

## 快速导航



- [功能特性](#功能特性)
- [安装](#安装)
- [快速开始](#快速开始)
- [两种注册方式](#2-用户注册) ⭐ 新增
- [完整文档](#文档)

## 功能特性

✨ **认证功能**
- 用户注册
- 用户登录
- Token 刷新
- Token 验证
- 密码修改
- 密码重置

✨ **用户管理**
- 获取用户信息
- 更新用户属性
- 启用/禁用用户
- 删除用户
- 用户列表查询
- 全局登出

✨ **组管理**
- 添加用户到组
- 从组中移除用户
- 查询用户所属组

✨ **DynamoDB 集成** 🆕
- 完整的 CRUD 操作
- 批量读写
- 条件查询和扫描
- 表管理（创建、删除、列表）
- 与 Cognito 无缝集成

## 安装

```bash
go get github.com/difyz9/cognito-sdk
```

## 快速开始

### 1. 创建客户端

```go
import cognito "github.com/difyz9/cognito-sdk"

// 方式1: 使用默认AWS凭证
client, err := cognito.NewClient(cognito.Config{
    UserPoolID: "ap-east-1_XXXXXXX",
    ClientID:   "xxxxxxxxxxxxxxxxxxxx",
    Region:     "ap-east-1",
})

// 方式2: 使用指定的凭证
client, err := cognito.NewClientWithCredentials(
    cognito.Config{
        UserPoolID: "ap-east-1_XXXXXXX",
        ClientID:   "xxxxxxxxxxxxxxxxxxxx",
        Region:     "ap-east-1",
    },
    "your-access-key-id",
    "your-secret-access-key",
)
```

### 2. 用户注册

SDK 提供两种注册方式，请根据您的业务场景选择：

#### 方式1: Register - 管理员创建用户（推荐用于后台管理）

用户创建后**立即可用**，无需邮箱验证：

```go
// 管理员创建用户（立即可用，无需验证）
resp, err := client.Register("user@example.com", "Password123!", "")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("注册成功! UserSub: %s\n", resp.UserSub)

// 用户可以立即登录
loginResp, err := client.Login("user@example.com", "Password123!")
```

**适用场景：**
- 后台管理系统创建用户
- 企业内部系统
- 不需要邮箱验证的场景

#### 方式2: SignUpUser - 用户自注册（推荐用于公开注册）

需要邮箱验证流程，更安全：

```go
// 步骤1: 用户自注册
resp, err := client.SignUpUser("user@example.com", "Password123!", "")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("注册成功! UserSub: %s\n", resp.UserSub)

if resp.CodeDeliveryDetails != nil {
    fmt.Printf("验证码已发送至: %s\n", resp.CodeDeliveryDetails.Destination)
}

// 步骤2: 用户输入收到的验证码
err = client.ConfirmSignUp("user@example.com", "123456")
if err != nil {
    log.Fatal(err)
}

// 步骤3: 验证后即可登录
loginResp, err := client.Login("user@example.com", "Password123!")
```

**适用场景：**
- 用户自助注册（网站、APP）
- 需要验证邮箱真实性
- 公开注册系统

**重新发送验证码：**
```go
details, err := client.ResendConfirmationCode("user@example.com")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("验证码已重新发送至: %s\n", details.Destination)
```

### 3. 用户登录

```go
resp, err := client.Login("user@example.com", "Password123!")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("ID Token: %s\n", resp.IDToken)
fmt.Printf("Access Token: %s\n", resp.AccessToken)
fmt.Printf("Refresh Token: %s\n", resp.RefreshToken)
```

### 4. 验证Token

```go
claims, err := client.VerifyToken(idToken)
if err != nil {
    log.Fatal(err)
}

// 或者直接提取用户信息
userInfo, err := client.GetUserInfoFromToken(idToken)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("用户ID: %s\n", userInfo["user_id"])
fmt.Printf("邮箱: %s\n", userInfo["email"])
```

### 5. 用户管理

```go
// 获取用户信息
userInfo, err := client.GetUser("user@example.com")
if err != nil {
    log.Fatal(err)
}

// 更新用户属性
err = client.UpdateUserAttributes("user@example.com", map[string]string{
    "name":         "张三",
    "phone_number": "+8613800138000",
    "address":      "北京市",
})

// 修改密码
err = client.ChangePassword("user@example.com", "OldPass123!", "NewPass123!")

// 禁用用户
err = client.DisableUser("user@example.com")

// 启用用户
err = client.EnableUser("user@example.com")

// 删除用户
err = client.DeleteUser("user@example.com")
```

### 6. 用户列表

```go
resp, err := client.ListUsers(cognito.ListUsersFilter{
    Limit:  10,
    Filter: "email ^= \"user@\"", // 可选的过滤条件
})

for _, user := range resp.Users {
    fmt.Printf("%s - %s (%s)\n", user.Username, user.Email, user.UserStatus)
}

// 分页查询
if resp.NextToken != nil {
    nextResp, _ := client.ListUsers(cognito.ListUsersFilter{
        Limit:     10,
        NextToken: resp.NextToken,
    })
}
```

### 7. 组管理

```go
// 添加用户到组
err := client.AddUserToGroup("user@example.com", "Admins")

// 列出用户所属的组
groups, err := client.ListUserGroups("user@example.com")
for _, group := range groups {
    fmt.Printf("%s: %s\n", group.GroupName, group.Description)
}

// 从组中移除用户
err = client.RemoveUserFromGroup("user@example.com", "Admins")
```

## API 文档

### Client 方法

#### 认证相关

- `Register(email, password, username string) (*RegisterResponse, error)` - 注册新用户
- `Login(username, password string) (*LoginResponse, error)` - 用户登录
- `RefreshToken(refreshToken string) (*RefreshTokenResponse, error)` - 刷新Token
- `ChangePassword(username, oldPassword, newPassword string) error` - 修改密码
- `ResetPassword(username, newPassword string) error` - 重置密码（管理员权限）
- `SignOut(username string) error` - 全局登出

#### Token 相关

- `VerifyToken(tokenString string) (jwt.MapClaims, error)` - 验证Token
- `GetUserInfoFromToken(tokenString string) (map[string]interface{}, error)` - 从Token提取用户信息
- `ValidateToken(tokenString string) bool` - 简单验证Token是否有效

#### 用户管理

- `GetUser(username string) (*UserInfo, error)` - 获取用户信息
- `UpdateUserAttributes(username string, attributes map[string]string) error` - 更新用户属性
- `EnableUser(username string) error` - 启用用户
- `DisableUser(username string) error` - 禁用用户
- `DeleteUser(username string) error` - 删除用户
- `ListUsers(filter ListUsersFilter) (*ListUsersResponse, error)` - 获取用户列表

#### 组管理

- `AddUserToGroup(username, groupName string) error` - 添加用户到组
- `RemoveUserFromGroup(username, groupName string) error` - 从组中移除用户
- `ListUserGroups(username string) ([]*GroupInfo, error)` - 列出用户所属的组

## 数据结构

### Config

```go
type Config struct {
    UserPoolID string  // Cognito用户池ID
    ClientID   string  // 应用客户端ID
    Region     string  // AWS区域
}
```

### UserInfo

```go
type UserInfo struct {
    Username      string            // 用户名
    Email         string            // 邮箱
    EmailVerified bool              // 邮箱是否验证
    Enabled       bool              // 是否启用
    UserStatus    string            // 用户状态
    CreatedAt     string            // 创建时间
    UpdatedAt     string            // 更新时间
    Attributes    map[string]string // 自定义属性
}
```

## 环境变量

建议使用环境变量管理敏感配置：

```bash
export COGNITO_USER_POOL_ID="ap-east-1_XXXXXXX"
export COGNITO_CLIENT_ID="xxxxxxxxxxxxxxxxxxxx"
export AWS_REGION="ap-east-1"
export AWS_ACCESS_KEY_ID="AKIAXXXXXXXXXXXXXXXX"
export AWS_SECRET_ACCESS_KEY="xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
```

## 完整示例

查看 [examples/basic/main.go](examples/basic/main.go) 获取更多使用示例。

## 错误处理

所有方法都会返回详细的错误信息：

```go
resp, err := client.Login("user@example.com", "wrong-password")
if err != nil {
    // 错误信息包含具体的失败原因
    log.Printf("登录失败: %v", err)
    return
}
```

## 注意事项

1. **AWS 凭证**: 确保配置了正确的 AWS 凭证（环境变量、~/.aws/credentials 或通过 IAM 角色）
2. **权限**: IAM 用户需要有 Cognito 和 DynamoDB 相关权限
3. **用户池配置**: 确保 Cognito 用户池已正确配置（认证流程、应用客户端等）
4. **密码策略**: 密码需要符合用户池的密码策略要求
5. **DynamoDB 表**: 使用 DynamoDB 功能前需要创建相应的表

## DynamoDB 集成 🆕

### 快速开始

```go
// 创建 Cognito 客户端
client, err := cognito.NewClient(cognito.Config{
    UserPoolID: "us-east-1_xxxxx",
    ClientID:   "xxxxxxxxxx",
    Region:     "us-east-1",
})

// 创建 DynamoDB 客户端（复用 AWS 配置）
dbClient := client.NewDynamoDBClient()

// 插入数据
type User struct {
    UserID   string `dynamodbav:"user_id"`
    Username string `dynamodbav:"username"`
    Email    string `dynamodbav:"email"`
}

user := User{
    UserID:   "user-001",
    Username: "john_doe",
    Email:    "john@example.com",
}

err = dbClient.PutItem(ctx, "Users", user)

// 获取数据
key := map[string]types.AttributeValue{
    "user_id": &types.AttributeValueMemberS{Value: "user-001"},
}

var retrievedUser User
err = dbClient.GetItem(ctx, "Users", key, &retrievedUser)
```

### DynamoDB 功能

- **基础操作**: PutItem, GetItem, UpdateItem, DeleteItem
- **批量操作**: BatchWriteItem
- **查询**: Query, Scan（支持条件过滤）
- **表管理**: CreateTable, DeleteTable, ListTables, DescribeTable, TableExists

### 示例项目

- [examples/dynamodb-basic/](examples/dynamodb-basic/) - DynamoDB 基础操作示例
- [examples/cognito-dynamodb-integration/](examples/cognito-dynamodb-integration/) - Cognito + DynamoDB 完整集成示例

查看 [API_REFERENCE.md](API_REFERENCE.md) 了解完整的 DynamoDB API 文档。

## 文档

📚 **完整文档**
- [DUAL_REGISTRATION.md](DUAL_REGISTRATION.md) - 🎉 **两种注册方式详解**（新）
- [API_REFERENCE.md](API_REFERENCE.md) - API 快速参考（包含 DynamoDB API）
- [USAGE_GUIDE.md](USAGE_GUIDE.md) - 完整使用指南
- [FIX_SUMMARY.md](FIX_SUMMARY.md) - 问题修复说明

📂 **示例代码**
- [examples/signup-flow/](examples/signup-flow/) - 注册方式对比示例
- [examples/basic/](examples/basic/) - 基础功能示例
- [examples/dynamodb-basic/](examples/dynamodb-basic/) - DynamoDB 基础操作 🆕
- [examples/cognito-dynamodb-integration/](examples/cognito-dynamodb-integration/) - Cognito + DynamoDB 集成 🆕

## 常见问题

### Q: 我应该使用哪种注册方式？
A: 
- **后台管理系统** → 使用 `Register`（管理员创建，立即可用）
- **公开网站/APP** → 使用 `SignUpUser`（用户自注册，邮箱验证）

查看 [DUAL_REGISTRATION.md](DUAL_REGISTRATION.md) 了解详细对比

### Q: SignUpUser 注册后无法登录？
A: 需要先调用 `ConfirmSignUp` 完成邮箱验证，或使用 `ConfirmUser`（管理员确认）

### Q: 登录时提示 "Auth flow not enabled"
A: 需要在 Cognito 控制台的应用客户端中启用 ADMIN_USER_PASSWORD_AUTH 认证流程

### Q: Token 验证失败
A: 检查 Region 和 UserPoolID 配置是否正确，Token 是否已过期

### Q: 权限不足
A: 确保 IAM 用户有以下权限：
- cognito-idp:AdminInitiateAuth
- cognito-idp:AdminGetUser
- cognito-idp:AdminUpdateUserAttributes
- cognito-idp:AdminCreateUser（Register 方法需要）
- 等其他需要的权限

## License

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！
