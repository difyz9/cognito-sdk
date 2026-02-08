# 双注册方式支持 - 更新说明

## 📝 更新概述

SDK 现在同时支持两种用户注册方式，用户可以根据业务场景自由选择：

1. **Register** - 管理员创建用户（AdminCreateUser）
2. **SignUpUser** - 用户自注册（SignUp + 验证流程）

## 🎯 设计理念

保持 API 简洁的同时，提供足够的灵活性：
- `Register` 保持原有实现（向后兼容）
- 新增 `SignUpUser` 提供标准注册流程
- 提供完整的验证码流程支持

## ✨ 新增 API

### 1. SignUpUser - 用户自注册
```go
func (c *Client) SignUpUser(email, password, username string) (*SignUpResponse, error)
```

### 2. ConfirmSignUp - 确认注册
```go
func (c *Client) ConfirmSignUp(username, confirmationCode string) error
```

### 3. ResendConfirmationCode - 重发验证码
```go
func (c *Client) ResendConfirmationCode(username string) (*CodeDeliveryDetails, error)
```

### 4. 新增数据结构
```go
type SignUpResponse struct {
    UserSub             string
    UserConfirmed       bool
    CodeDeliveryDetails *CodeDeliveryDetails
}

type CodeDeliveryDetails struct {
    Destination    string  // 邮箱地址（部分隐藏）
    DeliveryMedium string  // EMAIL 或 SMS
    AttributeName  string  // email 或 phone_number
}
```

## 📊 对比表格

| 特性 | Register | SignUpUser |
|------|---------|-----------|
| 实现方式 | AdminCreateUser | SignUp |
| 权限要求 | AWS管理员凭证 ✓ | ClientID即可 |
| 邮箱验证 | 不需要 | 需要 ✓ |
| 用户状态 | CONFIRMED | UNCONFIRMED → CONFIRMED |
| email_verified | true | false → true |
| 立即可用 | ✓ | 需验证后 |
| 适用场景 | 后台管理 | 公开注册 |
| 验证流程 | 无 | 验证码邮件 |

## 💡 使用示例

### 场景1: 后台管理系统创建用户

```go
// 管理员创建，立即可用
resp, err := client.Register("user@example.com", "Password123!", "")
loginResp, err := client.Login("user@example.com", "Password123!")
// ✅ 立即登录成功
```

### 场景2: 网站用户自助注册

```go
// 用户注册
resp, err := client.SignUpUser("user@example.com", "Password123!", "")
// 📧 用户收到验证码邮件

// 用户输入验证码
err = client.ConfirmSignUp("user@example.com", "123456")

// 验证后登录
loginResp, err := client.Login("user@example.com", "Password123!")
// ✅ 登录成功
```

## 📂 文件变更

### 修改的文件
1. **auth.go**
   - 保留 `Register` 方法（AdminCreateUser）
   - 新增 `SignUpUser` 方法
   - 新增 `ConfirmSignUp` 方法
   - 新增 `ResendConfirmationCode` 方法
   - 新增 `SignUpResponse` 和 `CodeDeliveryDetails` 类型

2. **user.go**
   - 新增 `SetPermanentPassword` 方法
   - 新增 `ConfirmUser` 方法（管理员确认）

### 新增的文件
1. **examples/signup-flow/main.go** - 两种方式对比演示
2. **examples/signup-flow/README.md** - 示例说明
3. **API_REFERENCE.md** - API 快速参考
4. **DUAL_REGISTRATION.md** - 本文档

### 更新的文档
1. **README.md** - 添加两种注册方式说明
2. **USAGE_GUIDE.md** - 添加详细对比和使用指南

## 🧪 测试验证

运行对比演示：
```bash
cd cognito_sdk/examples/signup-flow
go run main.go
```

预期输出：
```
【方式1】管理员创建用户
✓ 用户创建成功! UserSub: xxx
用户状态: CONFIRMED
✓ 登录成功!

【方式2】用户自注册
✓ 注册请求成功!
📧 验证码已发送至: s***@e***
用户状态: UNCONFIRMED
✗ 登录失败: User is not confirmed
✓ 用户已确认
✓ 登录成功!
```

## 🔄 向后兼容性

✅ **完全向后兼容**
- 原有的 `Register` 方法保持不变
- 现有代码无需修改即可正常工作
- 只是新增了更多选择

## 📖 相关文档

- [API 快速参考](API_REFERENCE.md) - 详细 API 文档
- [使用指南](USAGE_GUIDE.md) - 完整使用说明
- [修复说明](FIX_SUMMARY.md) - 原始问题修复记录
- [示例代码](examples/signup-flow/) - 两种方式对比演示

## 🎓 最佳实践建议

### 后台管理系统
```go
// 使用 Register - 简单高效
client.Register(email, password, "")
```

### 公开注册网站
```go
// 使用 SignUpUser - 标准流程
client.SignUpUser(email, password, "")
client.ConfirmSignUp(email, code)
```

### 测试环境
```go
// Register 或 SignUpUser + ConfirmUser
client.SignUpUser(email, password, "")
client.ConfirmUser(email) // 管理员直接确认，跳过验证码
```

## 📞 常见问题

**Q: 我应该用哪种注册方式？**  
A: 
- 后台管理、内部系统 → `Register`
- 公开网站、APP → `SignUpUser`

**Q: SignUpUser 必须要验证码吗？**  
A: 是的，这是标准的安全流程。如果你想跳过，可以用 `ConfirmUser`（需要管理员权限）

**Q: 我现有代码需要改吗？**  
A: 不需要，`Register` 方法保持不变，完全向后兼容

**Q: 验证码过期了怎么办？**  
A: 使用 `ResendConfirmationCode` 重新发送

**Q: 可以同时使用两种方式吗？**  
A: 可以，它们是独立的API，可以在同一个项目中混用

## 🚀 下一步

1. 查看 [API_REFERENCE.md](API_REFERENCE.md) 了解详细API
2. 运行 [示例代码](examples/signup-flow/) 查看实际效果
3. 阅读 [USAGE_GUIDE.md](USAGE_GUIDE.md) 了解更多使用场景
