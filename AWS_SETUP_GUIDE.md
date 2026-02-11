# AWS 凭证配置指南

## 问题诊断

如果遇到 `UnrecognizedClientException: The security token included in the request is invalid` 错误，通常是以下原因：

1. AWS 凭证未配置
2. AWS 凭证已过期（临时凭证）
3. AWS 凭证配置错误
4. 使用的 AWS 账户没有相应权限

## 解决方案

### 方案 1: 重新配置 AWS CLI（推荐）

```bash
# 配置 AWS CLI
aws configure

# 按提示输入：
# AWS Access Key ID: 你的 Access Key
# AWS Secret Access Key: 你的 Secret Key
# Default region name: ap-southeast-1
# Default output format: json
```

### 方案 2: 使用环境变量

```bash
# 设置 AWS 凭证环境变量
export AWS_ACCESS_KEY_ID="your-access-key-id"
export AWS_SECRET_ACCESS_KEY="your-secret-access-key"
export AWS_REGION="ap-southeast-1"

# 如果使用临时凭证（需要 Session Token）
export AWS_SESSION_TOKEN="your-session-token"
```

### 方案 3: 使用 IAM 角色（EC2/ECS 环境）

如果在 AWS EC2 或 ECS 上运行，配置 IAM 角色即可，无需手动配置凭证。

### 方案 4: 在代码中显式指定凭证

```go
// 使用显式凭证创建客户端
client, err := cognito.NewClientWithCredentials(
    cognito.Config{
        UserPoolID: "ap-southeast-1_ceThZId98",
        ClientID:   "6bf1au63kb7geu32up3cgnc6q0",
        Region:     "ap-southeast-1",
    },
    "your-access-key-id",      // AWS Access Key ID
    "your-secret-access-key",   // AWS Secret Access Key
)
```

## 验证配置

### 1. 验证 AWS CLI 配置

```bash
# 查看配置
aws configure list

# 测试凭证是否有效
aws sts get-caller-identity

# 应该返回类似：
# {
#     "UserId": "AIDXXXXXXXXXX",
#     "Account": "123456789012",
#     "Arn": "arn:aws:iam::123456789012:user/your-username"
# }
```

### 2. 测试 DynamoDB 权限

```bash
# 列出 DynamoDB 表（测试权限）
aws dynamodb list-tables --region ap-southeast-1

# 应该返回：
# {
#     "TableNames": []
# }
```

## 创建 DynamoDB 表

确认凭证有效后，创建所需的表：

### 创建 Users 表

```bash
aws dynamodb create-table \
    --table-name Users \
    --attribute-definitions \
        AttributeName=user_id,AttributeType=S \
    --key-schema \
        AttributeName=user_id,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST \
    --region ap-southeast-1
```

### 创建 UserProfiles 表（用于集成示例）

```bash
aws dynamodb create-table \
    --table-name UserProfiles \
    --attribute-definitions \
        AttributeName=user_id,AttributeType=S \
    --key-schema \
        AttributeName=user_id,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST \
    --region ap-southeast-1
```

### 创建 UserActivities 表（用于集成示例）

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
    --region ap-southeast-1
```

## 所需 IAM 权限

确保您的 AWS 用户或角色具有以下权限：

### Cognito 权限

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "cognito-idp:AdminInitiateAuth",
                "cognito-idp:AdminCreateUser",
                "cognito-idp:AdminGetUser",
                "cognito-idp:AdminUpdateUserAttributes",
                "cognito-idp:AdminDeleteUser",
                "cognito-idp:AdminSetUserPassword",
                "cognito-idp:AdminResetUserPassword",
                "cognito-idp:AdminUserGlobalSignOut",
                "cognito-idp:AdminEnableUser",
                "cognito-idp:AdminDisableUser",
                "cognito-idp:ListUsers",
                "cognito-idp:SignUp",
                "cognito-idp:ConfirmSignUp"
            ],
            "Resource": "arn:aws:cognito-idp:ap-southeast-1:*:userpool/ap-southeast-1_ceThZId98"
        }
    ]
}
```

### DynamoDB 权限

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "dynamodb:CreateTable",
                "dynamodb:DeleteTable",
                "dynamodb:DescribeTable",
                "dynamodb:ListTables",
                "dynamodb:PutItem",
                "dynamodb:GetItem",
                "dynamodb:UpdateItem",
                "dynamodb:DeleteItem",
                "dynamodb:Query",
                "dynamodb:Scan",
                "dynamodb:BatchWriteItem",
                "dynamodb:BatchGetItem"
            ],
            "Resource": [
                "arn:aws:dynamodb:ap-southeast-1:*:table/Users",
                "arn:aws:dynamodb:ap-southeast-1:*:table/UserProfiles",
                "arn:aws:dynamodb:ap-southeast-1:*:table/UserActivities",
                "arn:aws:dynamodb:ap-southeast-1:*:table/UserActivities/index/*"
            ]
        }
    ]
}
```

## 常见问题

### Q: 如何获取 AWS Access Key？

1. 登录 AWS Console
2. 进入 IAM → Users → 选择你的用户
3. Security credentials 标签
4. Create access key
5. 下载并保存凭证（只显示一次）

### Q: 临时凭证多久过期？

- IAM 用户的长期凭证：永久有效（除非手动删除）
- STS 临时凭证：通常 1-12 小时
- SSO 凭证：取决于 SSO 配置

### Q: 如何使用 AWS SSO？

```bash
# 配置 SSO
aws configure sso

# 登录
aws sso login --profile your-profile

# 使用特定 profile
aws dynamodb list-tables --profile your-profile --region ap-southeast-1
```

### Q: 如何切换 AWS Profile？

```bash
# 查看所有 profiles
aws configure list-profiles

# 使用特定 profile
export AWS_PROFILE=your-profile-name

# 或在命令中指定
aws dynamodb list-tables --profile your-profile-name
```

## 快速测试脚本

创建这个脚本来测试您的 AWS 配置：

```bash
#!/bin/bash

echo "=== AWS 配置测试 ==="
echo ""

echo "1. 检查 AWS CLI 版本..."
aws --version
echo ""

echo "2. 检查当前配置..."
aws configure list
echo ""

echo "3. 验证凭证..."
aws sts get-caller-identity
echo ""

echo "4. 测试 DynamoDB 访问..."
aws dynamodb list-tables --region ap-southeast-1
echo ""

echo "5. 测试 Cognito 访问..."
aws cognito-idp list-user-pools --max-results 10 --region ap-southeast-1
echo ""

echo "=== 测试完成 ==="
```

保存为 `test_aws_config.sh` 并运行：

```bash
chmod +x test_aws_config.sh
./test_aws_config.sh
```

## 环境变量设置脚本

创建 `.env` 文件（不要提交到 git）：

```bash
# AWS 凭证
export AWS_ACCESS_KEY_ID="your-access-key-id"
export AWS_SECRET_ACCESS_KEY="your-secret-access-key"
export AWS_REGION="ap-southeast-1"

# Cognito 配置
export COGNITO_USER_POOL_ID="ap-southeast-1_ceThZId98"
export COGNITO_CLIENT_ID="6bf1au63kb7geu32up3cgnc6q0"
```

使用：

```bash
# 加载环境变量
source .env

# 运行示例
go run examples/dynamodb-basic/main.go
```

## 下一步

1. ✅ 配置 AWS 凭证
2. ✅ 验证凭证有效性
3. ✅ 创建 DynamoDB 表
4. ✅ 运行示例代码

如果还有问题，请检查：
- AWS 账户是否激活
- 是否有费用限制
- 区域配置是否正确
- IAM 权限是否足够
