package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	cognito "github.com/difyz9/cognito-sdk"
)

var cognitoClient *cognito.Client

func main() {
	// 初始化Cognito客户端
	var err error
	cognitoClient, err = cognito.NewClient(cognito.Config{
		UserPoolID: os.Getenv("COGNITO_USER_POOL_ID"),
		ClientID:   os.Getenv("COGNITO_CLIENT_ID"),
		Region:     os.Getenv("AWS_REGION"),
	})
	if err != nil {
		log.Fatalf("初始化Cognito客户端失败: %v", err)
	}

	// 设置路由
	http.HandleFunc("/api/register", handleRegister)
	http.HandleFunc("/api/login", handleLogin)
	http.HandleFunc("/api/profile", authMiddleware(handleProfile))
	http.HandleFunc("/api/update", authMiddleware(handleUpdate))

	// 启动服务器
	port := ":8080"
	log.Printf("服务器启动在 http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// 注册接口
func handleRegister(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "请求参数错误"}`, http.StatusBadRequest)
		return
	}

	resp, err := cognitoClient.Register(req.Email, req.Password, "")
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"user_id": resp.UserSub,
	})
}

// 登录接口
func handleLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "请求参数错误"}`, http.StatusBadRequest)
		return
	}

	resp, err := cognitoClient.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":       true,
		"id_token":      resp.IDToken,
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
	})
}

// 获取用户信息
func handleProfile(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value("username").(string)

	userInfo, err := cognitoClient.GetUser(username)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"user":    userInfo,
	})
}

// 更新用户信息
func handleUpdate(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value("username").(string)

	var req struct {
		Attributes map[string]string `json:"attributes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "请求参数错误"}`, http.StatusBadRequest)
		return
	}

	err := cognitoClient.UpdateUserAttributes(username, req.Attributes)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}

// 认证中间件
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error": "缺少认证信息"}`, http.StatusUnauthorized)
			return
		}

		// 提取Token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, `{"error": "认证格式错误"}`, http.StatusUnauthorized)
			return
		}

		token := parts[1]

		// 验证Token并提取用户信息
		userInfo, err := cognitoClient.GetUserInfoFromToken(token)
		if err != nil {
			http.Error(w, `{"error": "Token验证失败"}`, http.StatusUnauthorized)
			return
		}

		// 将用户名存入Context
		ctx := r.Context()
		if username, ok := userInfo["username"].(string); ok {
			ctx = context.WithValue(ctx, "username", username)
		} else if email, ok := userInfo["email"].(string); ok {
			ctx = context.WithValue(ctx, "username", email)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
