package cognito_test

import (
	"os"
	"testing"

	cognito "github.com/difyz9/cognito-sdk"
)

// 测试创建客户端
func TestNewClient(t *testing.T) {
	client, err := cognito.NewClient(cognito.Config{
		UserPoolID: "ap-east-1_test",
		ClientID:   "test-client-id",
		Region:     "ap-east-1",
	})

	if err != nil {
		t.Fatalf("创建客户端失败: %v", err)
	}

	if client == nil {
		t.Fatal("客户端为nil")
	}

	cfg := client.GetConfig()
	if cfg.UserPoolID != "ap-east-1_test" {
		t.Errorf("UserPoolID不匹配: got %s, want ap-east-1_test", cfg.UserPoolID)
	}
}

// 测试配置验证
func TestNewClientWithInvalidConfig(t *testing.T) {
	_, err := cognito.NewClient(cognito.Config{
		UserPoolID: "",
		ClientID:   "test",
		Region:     "ap-east-1",
	})

	if err == nil {
		t.Error("期望返回错误，但成功了")
	}
}

// 集成测试示例（需要真实的Cognito配置）
func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// 从环境变量获取配置
	userPoolID := os.Getenv("TEST_COGNITO_USER_POOL_ID")
	clientID := os.Getenv("TEST_COGNITO_CLIENT_ID")
	region := os.Getenv("TEST_AWS_REGION")

	if userPoolID == "" || clientID == "" || region == "" {
		t.Skip("缺少测试环境变量")
	}

	client, err := cognito.NewClient(cognito.Config{
		UserPoolID: userPoolID,
		ClientID:   clientID,
		Region:     region,
	})

	if err != nil {
		t.Fatalf("创建客户端失败: %v", err)
	}

	// 测试注册
	t.Run("Register", func(t *testing.T) {
		email := "test@example.com"
		password := "Test@123456"

		resp, err := client.Register(email, password, "")
		if err != nil {
			t.Logf("注册失败（可能用户已存在）: %v", err)
			return
		}

		if resp.UserSub == "" {
			t.Error("UserSub为空")
		}

		t.Logf("注册成功: UserSub=%s", resp.UserSub)
	})

	// 测试登录
	t.Run("Login", func(t *testing.T) {
		email := "test@example.com"
		password := "Test@123456"

		resp, err := client.Login(email, password)
		if err != nil {
			t.Fatalf("登录失败: %v", err)
		}

		if resp.IDToken == "" {
			t.Error("IDToken为空")
		}
		if resp.AccessToken == "" {
			t.Error("AccessToken为空")
		}
		if resp.RefreshToken == "" {
			t.Error("RefreshToken为空")
		}

		t.Logf("登录成功: IDToken前缀=%s...", resp.IDToken[:30])

		// 测试Token验证
		t.Run("VerifyToken", func(t *testing.T) {
			claims, err := client.VerifyToken(resp.IDToken)
			if err != nil {
				t.Fatalf("Token验证失败: %v", err)
			}

			if claims["cognito:username"] == nil {
				t.Error("claims中缺少username")
			}

			t.Logf("Token验证成功: username=%v", claims["cognito:username"])
		})

		// 测试获取用户信息
		t.Run("GetUser", func(t *testing.T) {
			userInfo, err := client.GetUser(email)
			if err != nil {
				t.Fatalf("获取用户信息失败: %v", err)
			}

			if userInfo.Email != email {
				t.Errorf("邮箱不匹配: got %s, want %s", userInfo.Email, email)
			}

			t.Logf("用户信息: %+v", userInfo)
		})
	})
}
