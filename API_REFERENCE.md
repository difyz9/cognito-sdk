# API 快速参考 - 两种注册方式

## 注册相关 API

### 1. Register - 管理员创建用户

```go
func (c *Client) Register(email, password, username string) (*RegisterResponse, error)
```

**参数：**
- `email`: 用户邮箱
- `password`: 用户密码
- `username`: 用户名（可为空，默认使用email）

**返回：**
```go
type RegisterResponse struct {
    UserSub string  // 用户唯一ID
    Success bool    // 是否成功
}
```

**特点：**
- ✅ 用户立即可用（状态：CONFIRMED）
- ✅ email_verified 自动为 true
- ⚠️ 需要 AWS 管理员凭证

**示例：**
```go
resp, err := client.Register("user@example.com", "Password123!", "")
// 用户立即可登录
```

---

### 2. SignUpUser - 用户自注册

```go
func (c *Client) SignUpUser(email, password, username string) (*SignUpResponse, error)
```

**参数：**
- `email`: 用户邮箱
- `password`: 用户密码
- `username`: 用户名（可为空，默认使用email）

**返回：**
```go
type SignUpResponse struct {
    UserSub         string
    UserConfirmed   bool
    CodeDeliveryDetails *CodeDeliveryDetails
}

type CodeDeliveryDetails struct {
    Destination    string  // 邮箱地址（部分隐藏，如 s***@e***）
    DeliveryMedium string  // "EMAIL" 或 "SMS"
    AttributeName  string  // "email" 或 "phone_number"
}
```

**特点：**
- 📧 发送验证码到用户邮箱
- ⏳ 用户状态：UNCONFIRMED
- ℹ️ 不需要管理员凭证

**示例：**
```go
resp, err := client.SignUpUser("user@example.com", "Password123!", "")
if resp.CodeDeliveryDetails != nil {
    fmt.Printf("验证码已发送至: %s\n", resp.CodeDeliveryDetails.Destination)
}
// 用户需要验证后才能登录
```

---

### 3. ConfirmSignUp - 确认用户注册

```go
func (c *Client) ConfirmSignUp(username, confirmationCode string) error
```

**参数：**
- `username`: 用户名（通常是邮箱）
- `confirmationCode`: 6位验证码

**说明：**
- 用于确认 SignUpUser 创建的用户
- 验证码来自用户邮箱
- 确认后用户状态变为 CONFIRMED

**示例：**
```go
err := client.ConfirmSignUp("user@example.com", "123456")
if err != nil {
    log.Fatal(err)
}
// 现在用户可以登录了
```

---

### 4. ResendConfirmationCode - 重发验证码

```go
func (c *Client) ResendConfirmationCode(username string) (*CodeDeliveryDetails, error)
```

**参数：**
- `username`: 用户名（通常是邮箱）

**返回：**
- 验证码发送详情

**示例：**
```go
details, err := client.ResendConfirmationCode("user@example.com")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("验证码已重新发送至: %s\n", details.Destination)
```

---

### 5. ConfirmUser - 管理员确认用户

```go
func (c *Client) ConfirmUser(username string) error
```

**参数：**
- `username`: 用户名

**说明：**
- 管理员权限直接确认用户
- 跳过邮箱验证流程
- 适用于测试或特殊场景

**示例：**
```go
err := client.ConfirmUser("user@example.com")
// 用户状态变为 CONFIRMED
```

---

## 完整注册流程对比

### 流程1: 管理员创建（Register）

```go
// 1. 创建用户（一步完成）
resp, _ := client.Register("user@example.com", "Password123!", "")

// 2. 立即登录
loginResp, _ := client.Login("user@example.com", "Password123!")
```

**步骤：** 2步  
**用户体验：** 简单，但需要后台创建  
**安全性：** 依赖管理员权限控制

---

### 流程2: 用户自注册（SignUpUser）

```go
// 1. 用户提交注册
resp, _ := client.SignUpUser("user@example.com", "Password123!", "")

// 2. 系统发送验证码到邮箱
fmt.Printf("验证码已发送至: %s\n", resp.CodeDeliveryDetails.Destination)

// 3. 用户输入验证码
code := getUserInput() // 从表单获取

// 4. 确认注册
err := client.ConfirmSignUp("user@example.com", code)

// 5. 登录
loginResp, _ := client.Login("user@example.com", "Password123!")
```

**步骤：** 5步  
**用户体验：** 需要验证，但更标准  
**安全性：** 邮箱验证，更安全

---

## 选择建议

### 使用 Register 的场景

```
✓ 后台管理系统
✓ 企业内部系统  
✓ 批量创建用户
✓ 测试环境快速创建
✓ 管理员代替用户注册
```

### 使用 SignUpUser 的场景

```
✓ 公开网站注册
✓ 移动应用注册
✓ SaaS 产品
✓ 需要验证邮箱真实性
✓ 符合标准注册流程
```

---

## 错误处理

### Register 常见错误

```go
resp, err := client.Register(email, password, "")
if err != nil {
    switch {
    case strings.Contains(err.Error(), "UsernameExistsException"):
        // 用户已存在
    case strings.Contains(err.Error(), "InvalidPasswordException"):
        // 密码不符合策略
    default:
        // 其他错误
    }
}
```

### SignUpUser 常见错误

```go
resp, err := client.SignUpUser(email, password, "")
if err != nil {
    switch {
    case strings.Contains(err.Error(), "UsernameExistsException"):
        // 用户已存在
    case strings.Contains(err.Error(), "InvalidPasswordException"):
        // 密码不符合策略
    case strings.Contains(err.Error(), "InvalidParameterException"):
        // 参数无效（如邮箱格式）
    }
}

// 确认错误
err = client.ConfirmSignUp(username, code)
if err != nil {
    switch {
    case strings.Contains(err.Error(), "CodeMismatchException"):
        // 验证码错误
    case strings.Contains(err.Error(), "ExpiredCodeException"):
        // 验证码已过期
    case strings.Contains(err.Error(), "NotAuthorizedException"):
        // 用户已确认
    }
}
```

---

## 相关方法

### 辅助方法

```go
// 设置永久密码（内部使用）
func (c *Client) SetPermanentPassword(username, password string) error

// 获取用户信息
func (c *Client) GetUser(username string) (*UserInfo, error)

// 删除用户
func (c *Client) DeleteUser(username string) error
```

---

## 更多信息

- [完整使用指南](USAGE_GUIDE.md)
- [两种方式演示](examples/signup-flow/)
- [项目 README](README.md)
