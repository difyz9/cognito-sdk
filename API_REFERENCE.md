# API 快速参考 - 两种注册方式






















































































































































































































































































































































































































































- [AWS DynamoDB 官方文档](https://docs.aws.amazon.com/dynamodb/)- [Cognito + DynamoDB 集成示例](examples/cognito-dynamodb-integration/)- [DynamoDB 基础示例](examples/dynamodb-basic/)- [API 参考文档](API_REFERENCE.md)## 相关链接```}    ]        }            ]                "arn:aws:dynamodb:us-east-1:*:table/*"            "Resource": [            ],                "dynamodb:DeleteTable"                "dynamodb:CreateTable",                "dynamodb:ListTables",                "dynamodb:DescribeTable",                "dynamodb:BatchGetItem",                "dynamodb:BatchWriteItem",                "dynamodb:Scan",                "dynamodb:Query",                "dynamodb:DeleteItem",                "dynamodb:UpdateItem",                "dynamodb:GetItem",                "dynamodb:PutItem",            "Action": [            "Effect": "Allow",        {    "Statement": [    "Version": "2012-10-17",{```json使用 DynamoDB 功能需要以下 IAM 权限：## IAM 权限要求3. 对于复杂查询，考虑使用 ElasticSearch 或其他搜索服务2. 创建 GSI 支持不同的查询模式1. 使用 Query 配合 KeyConditionExpressionA: ### Q: 如何实现复杂查询？A: DynamoDB 限制每次最多 25 个项目，SDK 会自动分批处理超过 25 个项目的批量写入。### Q: 批量写入有大小限制吗？A: 可以，但需要直接使用 AWS SDK。本 SDK 提供的是基础的 CRUD 操作封装。### Q: 可以使用 DynamoDB Streams 吗？A: 使用分页扫描，或者考虑使用 Query 配合 GSI 来优化查询性能。### Q: 如何处理大量数据的扫描？A: 不会，`NewDynamoDBClient()` 复用 Cognito 客户端的 AWS 配置和凭证，无需额外配置。### Q: DynamoDB 客户端会创建新的 AWS 连接吗？## 常见问题- 记录所有数据访问日志- 实施细粒度的 IAM 权限- 使用加密存储敏感字段- 永远不要在 DynamoDB 存储敏感信息（如密码）### 5. 安全考虑- 使用 TTL 自动删除过期数据- 定期清理不需要的数据- 对于稳定流量，使用预配置容量更经济- 使用按需计费模式（PAY_PER_REQUEST）适合流量不可预测的场景### 4. 成本优化- 合理设计分区键避免热点- 实施数据缓存策略- 使用批量操作减少 API 调用次数- 为频繁查询的字段创建 GSI（全局二级索引）### 3. 性能优化对于需要原子性的操作，考虑使用 DynamoDB Transactions（需要直接使用 AWS SDK）。### 2. 使用事务```}    return err    // 实施重试逻辑或回滚操作    log.Printf("插入失败: %v", err)if err != nil {err := dbClient.PutItem(ctx, "Users", user)```go### 1. 错误处理## 最佳实践```}    EmailFrequency string `dynamodbav:"email_frequency"`    Notifications  bool   `dynamodbav:"notifications"`    Language       string `dynamodbav:"language"`    Theme          string `dynamodbav:"theme"`    UserID         string `dynamodbav:"user_id"`type UserPreferences struct {```go### UserPreferences 表```}    Metadata    string `dynamodbav:"metadata"`     // JSON 格式的额外信息    IPAddress   string `dynamodbav:"ip_address"`    Timestamp   string `dynamodbav:"timestamp"`    // GSI 排序键    Description string `dynamodbav:"description"`    Type        string `dynamodbav:"type"`         // login/logout/update等    UserID      string `dynamodbav:"user_id"`      // GSI 分区键    ActivityID  string `dynamodbav:"activity_id"`  // 主键: user_id + timestamptype Activity struct {```go### UserActivities 表```}    Status      string `dynamodbav:"status"`       // active/inactive/banned    LastLogin   string `dynamodbav:"last_login"`    UpdatedAt   string `dynamodbav:"updated_at"`    CreatedAt   string `dynamodbav:"created_at"`    AvatarURL   string `dynamodbav:"avatar_url"`    Bio         string `dynamodbav:"bio"`    PhoneNumber string `dynamodbav:"phone_number"`    FullName    string `dynamodbav:"full_name"`    Username    string `dynamodbav:"username"`    Email       string `dynamodbav:"email"`    CognitoSub  string `dynamodbav:"cognito_sub"`  // Cognito 用户标识    UserID      string `dynamodbav:"user_id"`      // 主键，使用 Cognito UserSubtype UserProfile struct {```go### UserProfiles 表## 数据模型设计建议   - 活动历史查询   - 用户信息管理   - 登录与活动记录   - 用户注册与 Profile 创建2. **Cognito + DynamoDB 集成**: [examples/cognito-dynamodb-integration/](examples/cognito-dynamodb-integration/)   - 表管理   - 批量读写   - CRUD 操作1. **DynamoDB 基础操作**: [examples/dynamodb-basic/](examples/dynamodb-basic/)查看完整的示例代码：## 完整示例```fmt.Printf("Profile: %+v\n", profile)fmt.Printf("Cognito Sub: %s\n", userInfo.Sub)// 合并信息err = dbClient.GetItem(ctx, "UserProfiles", key, &profile)var profile UserProfile}    "user_id": &types.AttributeValueMemberS{Value: userInfo.Sub},key := map[string]types.AttributeValue{// 从 DynamoDB 获取扩展信息userInfo, err := client.GetUserInfoFromToken(accessToken)// 从 Cognito Token 获取基本信息```go### 获取完整用户信息```err = dbClient.PutItem(ctx, "UserActivities", activity)}    Timestamp:  time.Now().Format(time.RFC3339),    Type:       "login",    UserID:     userInfo.Sub,    ActivityID: fmt.Sprintf("%s_%d", userInfo.Sub, time.Now().Unix()),activity := Activity{// 4. 记录登录活动    expressionValues, nil)err = dbClient.UpdateItem(ctx, "UserProfiles", key, updateExpression, }    ":time": &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},expressionValues := map[string]types.AttributeValue{updateExpression := "SET last_login = :time"}    "user_id": &types.AttributeValueMemberS{Value: userInfo.Sub},key := map[string]types.AttributeValue{// 3. 更新最后登录时间userInfo, err := client.GetUserInfoFromToken(loginResp.AccessToken)// 2. 获取用户信息loginResp, err := client.Login(ctx, username, password)// 1. Cognito 登录```go### 用户登录流程```err = dbClient.PutItem(ctx, "UserProfiles", profile)}    Status:     "active",    CreatedAt:  time.Now().Format(time.RFC3339),    Username:   username,    Email:      email,    CognitoSub: registerResp.UserSub,    UserID:     registerResp.UserSub,profile := UserProfile{// 2. DynamoDB 创建用户 ProfileregisterResp, err := client.Register(ctx, username, password, email)// 1. Cognito 注册```go### 用户注册流程## Cognito + DynamoDB 集成模式```fmt.Printf("项目数: %d\n", *tableDesc.ItemCount)fmt.Printf("状态: %s\n", tableDesc.TableStatus)fmt.Printf("表名: %s\n", *tableDesc.TableName)tableDesc, err := dbClient.DescribeTable(ctx, "Users")```go### 获取表信息```err := dbClient.DeleteTable(ctx, "Users")```go### 删除表```err := dbClient.CreateTable(ctx, input)}    BillingMode: types.BillingModePayPerRequest,    },        },            KeyType:       types.KeyTypeHash,            AttributeName: aws.String("user_id"),        {    KeySchema: []types.KeySchemaElement{    },        },            AttributeType: types.ScalarAttributeTypeS,            AttributeName: aws.String("user_id"),        {    AttributeDefinitions: []types.AttributeDefinition{    TableName: aws.String("Users"),input := &dynamodb.CreateTableInput{```go### 创建表```}    fmt.Println(table)for _, table := range tables {tables, err := dbClient.ListTables(ctx)```go### 列出所有表```}    fmt.Println("表不存在，需要创建")if !exists {exists, err := dbClient.TableExists(ctx, "Users")```go### 检查表是否存在## 表管理```    expressionAttributeValues, &users)err = dbClient.Query(ctx, "Users", keyConditionExpression, var users []User}    ":uid": &types.AttributeValueMemberS{Value: "user-001"},expressionAttributeValues := map[string]types.AttributeValue{keyConditionExpression := "user_id = :uid"```go### 查询```err = dbClient.Scan(ctx, "Users", &filteredUsers, filterExpression, filterValues)var filteredUsers []User}    ":domain": &types.AttributeValueMemberS{Value: "example.com"},filterValues := map[string]types.AttributeValue{filterExpression := aws.String("contains(email, :domain)")```go### 条件扫描```}    fmt.Printf("- %s (%s)\n", u.Username, u.Email)for _, u := range allUsers {fmt.Printf("找到 %d 个用户\n", len(allUsers))err = dbClient.Scan(ctx, "Users", &allUsers, nil, nil)var allUsers []User// 扫描所有记录```go### 扫描表```err = dbClient.BatchWriteItem(ctx, "Users", users)}    User{UserID: "user-003", Username: "bob", Email: "bob@example.com"},    User{UserID: "user-002", Username: "jane", Email: "jane@example.com"},    User{UserID: "user-001", Username: "john", Email: "john@example.com"},users := []interface{}{```go### 批量写入## 高级功能```err = dbClient.DeleteItem(ctx, "Users", key)}    "user_id": &types.AttributeValueMemberS{Value: "user-001"},key := map[string]types.AttributeValue{```go#### 删除数据```err = dbClient.UpdateItem(ctx, "Users", key, updateExpression, expressionValues, nil)}    ":email": &types.AttributeValueMemberS{Value: "newemail@example.com"},expressionValues := map[string]types.AttributeValue{updateExpression := "SET email = :email"}    "user_id": &types.AttributeValueMemberS{Value: "user-001"},key := map[string]types.AttributeValue{```go#### 更新数据```fmt.Printf("User: %s (%s)\n", retrievedUser.Username, retrievedUser.Email)err = dbClient.GetItem(ctx, "Users", key, &retrievedUser)var retrievedUser User}    "user_id": &types.AttributeValueMemberS{Value: "user-001"},key := map[string]types.AttributeValue{```go#### 获取数据```err = dbClient.PutItem(ctx, "Users", user)}    CreatedAt: time.Now().Format(time.RFC3339),    Email:     "john@example.com",    Username:  "john_doe",    UserID:    "user-001",user := User{```go#### 插入数据### 3. 基础操作```}    CreatedAt string `dynamodbav:"created_at"`    Email     string `dynamodbav:"email"`    Username  string `dynamodbav:"username"`    UserID    string `dynamodbav:"user_id"`type User struct {```go使用 `dynamodbav` 标签定义数据模型：### 2. 定义数据结构```ctx := context.Background()dbClient := client.NewDynamoDBClient()// 创建 DynamoDB 客户端（自动复用 AWS 配置）})    Region:     "us-east-1",    ClientID:   "xxxxxxxxxx",    UserPoolID: "us-east-1_xxxxx",client, err := cognito.NewClient(cognito.Config{// 创建 Cognito 客户端)    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"    cognito "github.com/difyz9/cognito-sdk"    "context"import (```go### 1. 初始化客户端## 快速开始本 SDK 已集成 DynamoDB 支持，可以轻松地将 Cognito 认证与 DynamoDB 数据存储结合使用。## 注册相关 API

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

## DynamoDB API 🆕

### 客户端初始化

```go
// 创建 DynamoDB 客户端（复用 Cognito 客户端的 AWS 配置）
func (c *Client) NewDynamoDBClient() *DynamoDBClient
```

**示例：**
```go
client, _ := cognito.NewClient(cognito.Config{...})
dbClient := client.NewDynamoDBClient()
```

---

### 基础操作

#### PutItem - 创建或更新项目

```go
func (db *DynamoDBClient) PutItem(ctx context.Context, tableName string, item interface{}) error
```

**参数：**
- `ctx`: 上下文
- `tableName`: 表名
- `item`: 要插入的项目（结构体）

**示例：**
```go
type User struct {
    UserID   string `dynamodbav:"user_id"`
    Username string `dynamodbav:"username"`
    Email    string `dynamodbav:"email"`
}

user := User{
    UserID:   "user-001",
    Username: "john",
    Email:    "john@example.com",
}

err := dbClient.PutItem(ctx, "Users", user)
```

---

#### GetItem - 获取项目

```go
func (db *DynamoDBClient) GetItem(ctx context.Context, tableName string, 
    key map[string]types.AttributeValue, result interface{}) error
```

**参数：**
- `ctx`: 上下文
- `tableName`: 表名
- `key`: 主键
- `result`: 结果对象的指针

**示例：**
```go
key := map[string]types.AttributeValue{
    "user_id": &types.AttributeValueMemberS{Value: "user-001"},
}

var user User
err := dbClient.GetItem(ctx, "Users", key, &user)
```

---

#### UpdateItem - 更新项目

```go
func (db *DynamoDBClient) UpdateItem(ctx context.Context, tableName string, 
    key map[string]types.AttributeValue, 
    updateExpression string, 
    expressionAttributeValues map[string]types.AttributeValue,
    expressionAttributeNames map[string]string) error
```

**示例：**
```go
key := map[string]types.AttributeValue{
    "user_id": &types.AttributeValueMemberS{Value: "user-001"},
}

updateExpr := "SET email = :email, updated_at = :time"
exprValues := map[string]types.AttributeValue{
    ":email": &types.AttributeValueMemberS{Value: "new@example.com"},
    ":time":  &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},
}

err := dbClient.UpdateItem(ctx, "Users", key, updateExpr, exprValues, nil)
```

---

#### DeleteItem - 删除项目

```go
func (db *DynamoDBClient) DeleteItem(ctx context.Context, tableName string, 
    key map[string]types.AttributeValue) error
```

**示例：**
```go
key := map[string]types.AttributeValue{
    "user_id": &types.AttributeValueMemberS{Value: "user-001"},
}

err := dbClient.DeleteItem(ctx, "Users", key)
```

---

### 查询操作

#### Query - 查询项目

```go
func (db *DynamoDBClient) Query(ctx context.Context, tableName string, 
    keyConditionExpression string,
    expressionAttributeValues map[string]types.AttributeValue, 
    results interface{}) error
```

**示例：**
```go
keyCondExpr := "user_id = :uid"
exprValues := map[string]types.AttributeValue{
    ":uid": &types.AttributeValueMemberS{Value: "user-001"},
}

var users []User
err := dbClient.Query(ctx, "Users", keyCondExpr, exprValues, &users)
```

---

#### Scan - 扫描表

```go
func (db *DynamoDBClient) Scan(ctx context.Context, tableName string, 
    results interface{}, 
    filterExpression *string, 
    expressionAttributeValues map[string]types.AttributeValue) error
```

**示例：**
```go
// 扫描所有
var allUsers []User
err := dbClient.Scan(ctx, "Users", &allUsers, nil, nil)

// 条件扫描
filterExpr := aws.String("contains(email, :domain)")
exprValues := map[string]types.AttributeValue{
    ":domain": &types.AttributeValueMemberS{Value: "example.com"},
}

var filteredUsers []User
err := dbClient.Scan(ctx, "Users", &filteredUsers, filterExpr, exprValues)
```

---

### 批量操作

#### BatchWriteItem - 批量写入

```go
func (db *DynamoDBClient) BatchWriteItem(ctx context.Context, tableName string, 
    items []interface{}) error
```

**示例：**
```go
users := []interface{}{
    User{UserID: "user-001", Username: "john", Email: "john@example.com"},
    User{UserID: "user-002", Username: "jane", Email: "jane@example.com"},
    User{UserID: "user-003", Username: "bob", Email: "bob@example.com"},
}

err := dbClient.BatchWriteItem(ctx, "Users", users)
```

**注意：** 每次最多写入 25 个项目，SDK 会自动分批处理。

---

### 表管理

#### CreateTable - 创建表

```go
func (db *DynamoDBClient) CreateTable(ctx context.Context, 
    input *dynamodb.CreateTableInput) error
```

**示例：**
```go
input := &dynamodb.CreateTableInput{
    TableName: aws.String("Users"),
    AttributeDefinitions: []types.AttributeDefinition{
        {
            AttributeName: aws.String("user_id"),
            AttributeType: types.ScalarAttributeTypeS,
        },
    },
    KeySchema: []types.KeySchemaElement{
        {
            AttributeName: aws.String("user_id"),
            KeyType:       types.KeyTypeHash,
        },
    },
    BillingMode: types.BillingModePayPerRequest,
}

err := dbClient.CreateTable(ctx, input)
```

---

#### DeleteTable - 删除表

```go
func (db *DynamoDBClient) DeleteTable(ctx context.Context, tableName string) error
```

---

#### ListTables - 列出所有表

```go
func (db *DynamoDBClient) ListTables(ctx context.Context) ([]string, error)
```

**示例：**
```go
tables, err := dbClient.ListTables(ctx)
for _, table := range tables {
    fmt.Println(table)
}
```

---

#### DescribeTable - 获取表信息

```go
func (db *DynamoDBClient) DescribeTable(ctx context.Context, tableName string) 
    (*types.TableDescription, error)
```

---

#### TableExists - 检查表是否存在

```go
func (db *DynamoDBClient) TableExists(ctx context.Context, tableName string) (bool, error)
```

**示例：**
```go
exists, err := dbClient.TableExists(ctx, "Users")
if !exists {
    // 创建表
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

// 获取AWS配置（DynamoDB复用）
func (c *Client) GetAWSConfig() aws.Config
```

---

## 更多信息

- [完整使用指南](USAGE_GUIDE.md)
- [两种注册方式](examples/signup-flow/)
- [DynamoDB 基础示例](examples/dynamodb-basic/)
- [Cognito + DynamoDB 集成](examples/cognito-dynamodb-integration/)
