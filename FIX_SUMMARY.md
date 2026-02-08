# Cognito SDK 修复说明

## 问题诊断

### 原因
cognito_golang 项目功能正常，但 cognito_sdk 封装的示例运行报错 `UserNotConfirmedException: User is not confirmed`

### 根本原因
两个项目使用了不同的用户注册方式：

**cognito_golang (正常工作):**
```go
// 使用 AdminCreateUser - 管理员创建用户
cognitoidentityprovider.AdminCreateUser()
- 自动设置 email_verified = true
- 用户状态为 CONFIRMED (已确认)
- 可以直接登录
```

**cognito_sdk (原始版本-有问题):**
```go
// 使用 SignUp - 用户自注册
cognitoidentityprovider.SignUp()
- email_verified = false (需要验证)
- 用户状态为 UNCONFIRMED (未确认)
- 无法登录，会报错: User is not confirmed
```

## 修复内容

### 1. 修改 Register 方法 (auth.go)
将 `SignUp` 改为 `AdminCreateUser`:
```go
// 旧代码 - 有问题
func (c *Client) Register(...) {
    input := &cognitoidentityprovider.SignUpInput{...}
    c.cognitoClient.SignUp(...)
}

// 新代码 - 已修复
func (c *Client) Register(...) {
    // 使用 AdminCreateUser (管理员创建，自动确认)
    input := &cognitoidentityprovider.AdminCreateUserInput{
        MessageAction: types.MessageActionTypeSuppress,
        UserAttributes: []types.AttributeType{
            {Name: "email_verified", Value: "true"},
        },
    }
    c.cognitoClient.AdminCreateUser(...)
    
    // 设置永久密码
    c.SetPermanentPassword(...)
}
```

### 2. 新增辅助方法 (user.go)
```go
// SetPermanentPassword - 设置永久密码
func (c *Client) SetPermanentPassword(username, password string) error

// ConfirmUser - 管理员确认用户（用于SignUp创建的用户）
func (c *Client) ConfirmUser(username string) error
```

## 测试结果

### 修复前
```
=== 登录示例 ===
登录失败: UserNotConfirmedException: User is not confirmed.
```

### 修复后
```
=== 登录示例 ===
登录成功!
ID Token: eyJraWQi...
Access Token: eyJraWQi...
Refresh Token: eyJjdHki...
过期时间: 3600 秒

=== Token验证示例 ===
Token验证成功! Claims: map[...]
```

## 使用方式对比

### cognito_golang (原始实现)
```go
// internal/auth/auth.go
RegisterUser(email, password, username)
  -> AdminCreateUser + SetPermanentPassword
```

### cognito_sdk (修复后)
```go
// client.Register
client.Register(email, password, username)
  -> AdminCreateUser + SetPermanentPassword
```

两者现在使用相同的底层实现，保证了一致性。

## 关键差异

| 方法 | 用户状态 | email_verified | 是否需要确认 | 适用场景 |
|------|---------|---------------|------------|---------|
| SignUp | UNCONFIRMED | false | 是 | 用户自注册，需要邮箱验证 |
| AdminCreateUser | CONFIRMED | true | 否 | 管理员创建，无需验证 |

## 最佳实践

1. **后端管理用户** - 使用 `AdminCreateUser` (当前实现)
   - 适合后台管理系统
   - 不需要邮件验证流程
   - 用户可以立即登录

2. **用户自注册** - 使用 `SignUp` + 验证流程
   - 需要实现邮箱验证
   - 调用 `ConfirmSignUp` 或 `AdminConfirmSignUp`
   - 更安全，但流程复杂

## 测试验证
运行测试脚本验证修复：
```bash
cd cognito_sdk/examples/basic
go run test_clean.go
```

输出：
```
删除旧的测试用户...
删除成功!

创建新用户...
注册成功! UserSub: 794a559c-...

测试登录...
登录成功!
```
