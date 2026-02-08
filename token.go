package cognito

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

// JWK JSON Web Key
type JWK struct {
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	Alg string `json:"alg"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// JWKS JSON Web Key Set
type JWKS struct {
	Keys []JWK `json:"keys"`
}

// VerifyToken 验证JWT Token
func (c *Client) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	// 获取JWKS
	jwksURL := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json",
		c.config.Region, c.config.UserPoolID)

	resp, err := http.Get(jwksURL)
	if err != nil {
		return nil, fmt.Errorf("获取JWKS失败: %w", err)
	}
	defer resp.Body.Close()

	var jwks JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, fmt.Errorf("解析JWKS失败: %w", err)
	}

	// 解析Token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("意外的签名方法: %v", token.Header["alg"])
		}

		// 获取kid
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("Token header中缺少kid")
		}

		// 查找对应的JWK
		for _, key := range jwks.Keys {
			if key.Kid == kid {
				return convertJWKToPublicKey(key)
			}
		}

		return nil, fmt.Errorf("未找到匹配的JWK")
	})

	if err != nil {
		return nil, fmt.Errorf("Token验证失败: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("Token无效")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("无法解析claims")
	}

	return claims, nil
}

// convertJWKToPublicKey 将JWK转换为RSA公钥
func convertJWKToPublicKey(jwk JWK) (*rsa.PublicKey, error) {
	// 解码N和E
	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, fmt.Errorf("解码N失败: %w", err)
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, fmt.Errorf("解码E失败: %w", err)
	}

	// 构造RSA公钥
	n := new(big.Int).SetBytes(nBytes)
	e := new(big.Int).SetBytes(eBytes)

	publicKey := &rsa.PublicKey{
		N: n,
		E: int(e.Int64()),
	}

	return publicKey, nil
}

// GetUserInfoFromToken 从Token中提取用户信息
func (c *Client) GetUserInfoFromToken(tokenString string) (map[string]interface{}, error) {
	claims, err := c.VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}

	userInfo := make(map[string]interface{})

	// 提取常用字段
	if username, ok := claims["cognito:username"]; ok {
		userInfo["username"] = username
	}
	if email, ok := claims["email"]; ok {
		userInfo["email"] = email
	}
	if sub, ok := claims["sub"]; ok {
		userInfo["user_id"] = sub
	}

	return userInfo, nil
}

// ValidateToken 简单验证Token是否有效（不解析内容）
func (c *Client) ValidateToken(tokenString string) bool {
	_, err := c.VerifyToken(tokenString)
	return err == nil
}
