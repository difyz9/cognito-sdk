# 两种注册方式对比演示

本示例演示了 Cognito SDK 提供的两种用户注册方式。

## 运行示例

```bash
cd cognito_sdk/examples/signup-flow
go run main.go
```

## 两种注册方式

### 方式1: Register - 管理员创建用户

**特点：**
- ✅ 用户创建后立即可用
- ✅ 无需邮箱验证
- ✅ 用户状态自动为 CONFIRMED
- ⚠️ 需要AWS管理员凭证

**使用场景：**
- 后台管理系统
- 企业内部系统
- 批量创建用户

**代码示例：**
```go
resp, err := client.Register("user@example.com", "Password123!", "")
if err != nil {
    log.Fatal(err)
}

// 用户可以立即登录
loginResp, err := client.Login("user@example.com", "Password123!")
```

---

### 方式2: SignUpUser - 用户自注册

**特点：**
- 📧 需要邮箱验证
- 📧 用户收到验证码邮件
- ⏳ 注册后状态为 UNCONFIRMED
- ✅ 更安全，适合公开注册

**使用场景：**
- 用户自助注册
- SaaS 应用
- 公开网站/APP注册

**代码示例：**
```go
// 1. 用户注册
resp, err := client.SignUpUser("user@example.com", "Password123!", "")
if err != nil {
    log.Fatal(err)
}

// 2. 确认注册（用户输入邮件中的验证码）
err = client.ConfirmSignUp("user@example.com", "123456")
if err != nil {
    log.Fatal(err)
}

// 3. 验证后可以登录
loginResp, err := client.Login("user@example.com", "Password123!")
```

## 对比表格

| 特性 | Register | SignUpUser |
|------|---------|-----------|
| **权限要求** | AWS管理员凭证 | 仅需ClientID |
| **邮箱验证** | ❌ 不需要 | ✅ 需要 |
| **用户状态** | CONFIRMED | UNCONFIRMED → CONFIRMED |
| **email_verified** | true | false → true |
| **立即登录** | ✅ 可以 | ❌ 需先验证 |
| **适用场景** | 后台管理 | 公开注册 |
| **实现方式** | AdminCreateUser | SignUp + ConfirmSignUp |

## 运行输出示例

```
===========================================
AWS Cognito 注册方式演示
===========================================

【方式1】管理员创建用户 - Register
--------------------------------------------------
✓ 用户创建成功! UserSub: xxx
用户状态: CONFIRMED
邮箱验证: true
✓ 登录成功!

【方式2】用户自注册 - SignUpUser
--------------------------------------------------
✓ 注册请求成功! UserSub: xxx
用户已确认: false
📧 验证码已发送至: s***@e***
用户状态: UNCONFIRMED ⚠️
✗ 登录失败: User is not confirmed
✓ 用户已确认
用户状态: CONFIRMED ✅
✓ 登录成功!
```

## 相关文档

- [完整使用指南](../../USAGE_GUIDE.md)
- [修复说明](../../FIX_SUMMARY.md)
- [项目README](../../README.md)
