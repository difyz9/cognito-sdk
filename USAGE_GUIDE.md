# SDK 使用指南

## 目录
- [在第三方项目中使用此SDK](#在第三方项目中使用此sdk)
- [两种注册方式说明](#两种注册方式说明)
- [完整示例项目](#完整示例项目)

## 两种注册方式说明

SDK 提供两种用户注册方式，适用于不同的业务场景：

### 方式1: Register - 管理员创建用户（推荐用于后台管理系统）

**特点：**
- ✅ 用户创建后**立即可用**，无需邮箱验证
- ✅ 用户状态自动为 `CONFIRMED`
- ✅ `email_verified` 自动设置为 `true`
- ✅ 适合后台管理系统创建用户
- ⚠️ 需要管理员权限（AWS凭证）

**使用示例：**
```go
// 管理员创建用户（立即可用）
resp, err := client.Register("user@example.com", "Password123!", "")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("用户创建成功! UserSub: %s\n", resp.UserSub)

// 用户可以立即登录
loginResp, err := client.Login("user@example.com", "Password123!")
```

**适用场景：**
- 后台管理系统创建用户
- 企业内部系统
- 不需要邮箱验证流程的场景

---

### 方式2: SignUpUser - 用户自注册（推荐用于公开注册）

**特点：**
- 📧 需要**邮箱验证**流程
- 📧 用户会收到包含验证码的邮件
- ⏳ 注册后状态为 `UNCONFIRMED`，需要调用 `ConfirmSignUp` 完成验证
- ✅ 更安全，适合公开注册场景
- ℹ️ 不需要管理员权限

**使用示例：**
```go
// 步骤1: 用户自注册
resp, err := client.SignUpUser("user@example.com", "Password123!", "")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("注册成功! UserSub: %s\n", resp.UserSub)
fmt.Printf("用户已确认: %v\n", resp.UserConfirmed) // false

if resp.CodeDeliveryDetails != nil {
    fmt.Printf("验证码已发送至: %s\n", resp.CodeDeliveryDetails.Destination)
    fmt.Printf("发送方式: %s\n", resp.CodeDeliveryDetails.DeliveryMedium)
}

// 步骤2: 用户输入收到的验证码进行确认
err = client.ConfirmSignUp("user@example.com", "123456") // 用户输入的验证码
if err != nil {
    log.Fatal(err)
}
fmt.Println("邮箱验证成功!")

// 步骤3: 现在用户可以登录了
loginResp, err := client.Login("user@example.com", "Password123!")
```

**重新发送验证码：**
```go
details, err := client.ResendConfirmationCode("user@example.com")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("验证码已重新发送至: %s\n", details.Destination)
```

**适用场景：**
- 用户自助注册（网站、APP）
- 需要验证邮箱真实性
- 公开注册系统
- SaaS 应用

---

### 两种方式对比

| 特性 | Register (管理员创建) | SignUpUser (用户自注册) |
|------|-------------------|---------------------|
| **权限要求** | 需要AWS管理员凭证 | 不需要（仅ClientID） |
| **邮箱验证** | ❌ 不需要 | ✅ 需要验证码 |
| **用户状态** | CONFIRMED（已确认） | UNCONFIRMED（未确认） |
| **email_verified** | true | false → true（验证后） |
| **立即登录** | ✅ 可以 | ❌ 需先验证 |
| **适用场景** | 后台管理、内部系统 | 公开注册、SaaS应用 |
| **安全性** | 依赖管理员权限 | 邮箱验证更安全 |

---

## 在第三方项目中使用此SDK

### 方式1: 通过Go Modules引用（推荐）

1. **发布SDK到GitHub**

```bash
cd /Users/apple/opt/difyz_2026/0205/cognito_sdk
git init
git add .
git commit -m "Initial commit"
git remote add origin https://github.com/yourusername/cognito-sdk.git
git push -u origin main
git tag v1.0.0
git push --tags
```

2. **在第三方项目中引用**

```bash
# 在你的项目中
go get github.com/yourusername/cognito-sdk@v1.0.0
```

3. **使用SDK**

```go
package main

import (
    "log"
    cognito "github.com/yourusername/cognito-sdk"
)

func main() {
    client, err := cognito.NewClient(cognito.Config{
        UserPoolID: "your-pool-id",
        ClientID:   "your-client-id",
        Region:     "ap-east-1",
    })
    if err != nil {
        log.Fatal(err)
    }

    // 使用SDK功能
    resp, err := client.Login("user@example.com", "password")
    // ...
}
```

### 方式2: 本地引用（开发测试）

1. **在你的项目的go.mod中添加replace指令**

```go
module myproject

go 1.21

require github.com/difyz9/cognito-sdk v0.0.0

replace github.com/difyz9/cognito-sdk => /Users/apple/opt/difyz_2026/0205/cognito_sdk
```

2. **直接导入使用**

```go
import cognito "github.com/difyz9/cognito-sdk"
```

### 方式3: 作为vendor依赖

```bash
# 在你的项目目录
mkdir -p vendor/github.com/yourname
cp -r /Users/apple/opt/difyz_2026/0205/cognito_sdk vendor/github.com/yourname/

# 在go.mod中启用vendor
go mod vendor
```

## 完整示例项目

### 创建新项目

```bash
mkdir my-cognito-app
cd my-cognito-app
go mod init my-cognito-app
```

### 安装依赖

```bash
# 方式1: 从GitHub安装
go get github.com/yourusername/cognito-sdk

# 方式2: 本地开发，在go.mod添加replace
```

### 创建main.go

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    cognito "github.com/yourusername/cognito-sdk"
)

func main() {
    // 从环境变量读取配置
    client, err := cognito.NewClient(cognito.Config{
        UserPoolID: os.Getenv("COGNITO_USER_POOL_ID"),
        ClientID:   os.Getenv("COGNITO_CLIENT_ID"),
        Region:     os.Getenv("AWS_REGION"),
    })
    if err != nil {
        log.Fatalf("初始化失败: %v", err)
    }

    // 注册用户
    registerResp, err := client.Register(
        "test@example.com",
        "Test@123456",
        "",
    )
    if err != nil {
        log.Printf("注册失败: %v", err)
    } else {
        fmt.Printf("注册成功! UserSub: %s\n", registerResp.UserSub)
    }

    // 登录
    loginResp, err := client.Login("test@example.com", "Test@123456")
    if err != nil {
        log.Printf("登录失败: %v", err)
        return
    }
    
    fmt.Printf("登录成功!\n")
    fmt.Printf("Token: %s...\n", loginResp.IDToken[:50])

    // 验证Token
    claims, err := client.VerifyToken(loginResp.IDToken)
    if err != nil {
        log.Printf("Token验证失败: %v", err)
        return
    }
    
    fmt.Printf("Token有效! 用户: %v\n", claims["cognito:username"])

    // 获取用户信息
    userInfo, err := client.GetUser("test@example.com")
    if err != nil {
        log.Printf("获取用户信息失败: %v", err)
        return
    }
    
    fmt.Printf("用户信息: %+v\n", userInfo)
}
```

### 运行

```bash
export COGNITO_USER_POOL_ID="ap-east-1_XXXXXXX"
export COGNITO_CLIENT_ID="xxxxxxxxxxxxxxxxxxxx"
export AWS_REGION="ap-east-1"
export AWS_ACCESS_KEY_ID="AKIAXXXXXXXXXXXXXXXX"
export AWS_SECRET_ACCESS_KEY="xxxxxxxxxxxxxxxx"

go run main.go
```

## Web应用集成示例

参考 `examples/http-server/main.go` 查看如何在Web应用中集成SDK。

```bash
cd examples/http-server
go run main.go
```

然后可以通过HTTP接口调用：

```bash
# 注册
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"Pass123!"}'

# 登录
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"Pass123!"}'

# 获取用户信息（需要Token）
curl -X GET http://localhost:8080/api/profile \
  -H "Authorization: Bearer <your-id-token>"

# 更新用户信息
curl -X POST http://localhost:8080/api/update \
  -H "Authorization: Bearer <your-id-token>" \
  -H "Content-Type: application/json" \
  -d '{"attributes":{"name":"张三","phone_number":"+8613800138000"}}'
```

## 最佳实践

### 1. 配置管理

使用配置文件或环境变量管理敏感信息：

```go
type AppConfig struct {
    Cognito struct {
        UserPoolID string `env:"COGNITO_USER_POOL_ID"`
        ClientID   string `env:"COGNITO_CLIENT_ID"`
        Region     string `env:"AWS_REGION"`
    }
}

func initCognitoClient(cfg AppConfig) (*cognito.Client, error) {
    return cognito.NewClient(cognito.Config{
        UserPoolID: cfg.Cognito.UserPoolID,
        ClientID:   cfg.Cognito.ClientID,
        Region:     cfg.Cognito.Region,
    })
}
```

### 2. 错误处理

```go
resp, err := client.Login(username, password)
if err != nil {
    // 区分不同的错误类型
    if strings.Contains(err.Error(), "NotAuthorizedException") {
        return errors.New("用户名或密码错误")
    }
    if strings.Contains(err.Error(), "UserNotFoundException") {
        return errors.New("用户不存在")
    }
    return fmt.Errorf("登录失败: %w", err)
}
```

### 3. 单例模式

```go
var (
    cognitoClient *cognito.Client
    once          sync.Once
)

func GetCognitoClient() *cognito.Client {
    once.Do(func() {
        var err error
        cognitoClient, err = cognito.NewClient(cognito.Config{
            UserPoolID: os.Getenv("COGNITO_USER_POOL_ID"),
            ClientID:   os.Getenv("COGNITO_CLIENT_ID"),
            Region:     os.Getenv("AWS_REGION"),
        })
        if err != nil {
            log.Fatal(err)
        }
    })
    return cognitoClient
}
```

### 4. 中间件封装

```go
func AuthMiddleware(client *cognito.Client) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            token := extractToken(r)
            if token == "" {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }

            userInfo, err := client.GetUserInfoFromToken(token)
            if err != nil {
                http.Error(w, "Invalid token", http.StatusUnauthorized)
                return
            }

            ctx := context.WithValue(r.Context(), "user", userInfo)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

## 常见问题

### Q: 如何更新SDK到最新版本？

```bash
go get -u github.com/yourusername/cognito-sdk
go mod tidy
```

### Q: 如何在Docker中使用？

Dockerfile示例：

```dockerfile
FROM golang:1.21-alpine

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o app .

CMD ["./app"]
```

### Q: 如何进行单元测试？

使用接口和依赖注入便于测试：

```go
type CognitoService interface {
    Login(username, password string) (*cognito.LoginResponse, error)
    GetUser(username string) (*cognito.UserInfo, error)
}

// 在测试中使用mock
type MockCognitoService struct{}

func (m *MockCognitoService) Login(username, password string) (*cognito.LoginResponse, error) {
    return &cognito.LoginResponse{IDToken: "mock-token"}, nil
}
```

## 技术支持

如有问题，请查看：
- [GitHub Issues](https://github.com/yourusername/cognito-sdk/issues)
- [完整文档](README.md)
- [示例代码](examples/)
