package cognito

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

// UserInfo 用户信息
type UserInfo struct {
	Username      string            `json:"username"`
	Email         string            `json:"email"`
	EmailVerified bool              `json:"email_verified"`
	Enabled       bool              `json:"enabled"`
	UserStatus    string            `json:"user_status"`
	CreatedAt     string            `json:"created_at,omitempty"`
	UpdatedAt     string            `json:"updated_at,omitempty"`
	Attributes    map[string]string `json:"attributes,omitempty"`
}

// GetUser 获取用户信息
func (c *Client) GetUser(username string) (*UserInfo, error) {
	input := &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String(c.config.UserPoolID),
		Username:   aws.String(username),
	}

	resp, err := c.cognitoClient.AdminGetUser(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	userInfo := &UserInfo{
		Username:   *resp.Username,
		Enabled:    resp.Enabled,
		UserStatus: string(resp.UserStatus),
		Attributes: make(map[string]string),
	}

	if resp.UserCreateDate != nil {
		userInfo.CreatedAt = resp.UserCreateDate.String()
	}
	if resp.UserLastModifiedDate != nil {
		userInfo.UpdatedAt = resp.UserLastModifiedDate.String()
	}

	// 提取用户属性
	for _, attr := range resp.UserAttributes {
		key := *attr.Name
		value := *attr.Value
		switch key {
		case "email":
			userInfo.Email = value
		case "email_verified":
			userInfo.EmailVerified = value == "true"
		default:
			userInfo.Attributes[key] = value
		}
	}

	return userInfo, nil
}

// UpdateUserAttributes 更新用户属性
func (c *Client) UpdateUserAttributes(username string, attributes map[string]string) error {
	if len(attributes) == 0 {
		return fmt.Errorf("属性不能为空")
	}

	userAttrs := make([]types.AttributeType, 0, len(attributes))
	for key, value := range attributes {
		userAttrs = append(userAttrs, types.AttributeType{
			Name:  aws.String(key),
			Value: aws.String(value),
		})
	}

	input := &cognitoidentityprovider.AdminUpdateUserAttributesInput{
		UserPoolId:     aws.String(c.config.UserPoolID),
		Username:       aws.String(username),
		UserAttributes: userAttrs,
	}

	_, err := c.cognitoClient.AdminUpdateUserAttributes(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("更新用户属性失败: %w", err)
	}

	return nil
}

// EnableUser 启用用户
func (c *Client) EnableUser(username string) error {
	input := &cognitoidentityprovider.AdminEnableUserInput{
		UserPoolId: aws.String(c.config.UserPoolID),
		Username:   aws.String(username),
	}

	_, err := c.cognitoClient.AdminEnableUser(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("启用用户失败: %w", err)
	}

	return nil
}

// DisableUser 禁用用户
func (c *Client) DisableUser(username string) error {
	input := &cognitoidentityprovider.AdminDisableUserInput{
		UserPoolId: aws.String(c.config.UserPoolID),
		Username:   aws.String(username),
	}

	_, err := c.cognitoClient.AdminDisableUser(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("禁用用户失败: %w", err)
	}

	return nil
}

// SetPermanentPassword 为用户设置永久密码
func (c *Client) SetPermanentPassword(username, password string) error {
	input := &cognitoidentityprovider.AdminSetUserPasswordInput{
		UserPoolId: aws.String(c.config.UserPoolID),
		Username:   aws.String(username),
		Password:   aws.String(password),
		Permanent:  true,
	}

	_, err := c.cognitoClient.AdminSetUserPassword(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("设置永久密码失败: %w", err)
	}

	return nil
}

// ConfirmUser 管理员确认用户（用于SignUp创建的用户）
func (c *Client) ConfirmUser(username string) error {
	input := &cognitoidentityprovider.AdminConfirmSignUpInput{
		UserPoolId: aws.String(c.config.UserPoolID),
		Username:   aws.String(username),
	}

	_, err := c.cognitoClient.AdminConfirmSignUp(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("确认用户失败: %w", err)
	}

	return nil
}

// ListUsersFilter 列表筛选条件
type ListUsersFilter struct {
	Limit      int32
	NextToken  *string
	Filter     string   // 例如: "email ^= \"user@\""
	Attributes []string
}

// ListUsersResponse 用户列表响应
type ListUsersResponse struct {
	Users     []*UserInfo
	NextToken *string
}

// ListUsers 获取用户列表
func (c *Client) ListUsers(filter ListUsersFilter) (*ListUsersResponse, error) {
	input := &cognitoidentityprovider.ListUsersInput{
		UserPoolId: aws.String(c.config.UserPoolID),
	}

	if filter.Limit > 0 {
		input.Limit = aws.Int32(filter.Limit)
	} else {
		input.Limit = aws.Int32(10) // 默认10个
	}

	if filter.NextToken != nil {
		input.PaginationToken = filter.NextToken
	}

	if filter.Filter != "" {
		input.Filter = aws.String(filter.Filter)
	}

	if len(filter.Attributes) > 0 {
		input.AttributesToGet = filter.Attributes
	}

	resp, err := c.cognitoClient.ListUsers(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("获取用户列表失败: %w", err)
	}

	users := make([]*UserInfo, 0, len(resp.Users))
	for _, user := range resp.Users {
		userInfo := &UserInfo{
			Username:   *user.Username,
			Enabled:    user.Enabled,
			UserStatus: string(user.UserStatus),
			Attributes: make(map[string]string),
		}

		if user.UserCreateDate != nil {
			userInfo.CreatedAt = user.UserCreateDate.String()
		}
		if user.UserLastModifiedDate != nil {
			userInfo.UpdatedAt = user.UserLastModifiedDate.String()
		}

		// 提取用户属性
		for _, attr := range user.Attributes {
			key := *attr.Name
			value := *attr.Value
			switch key {
			case "email":
				userInfo.Email = value
			case "email_verified":
				userInfo.EmailVerified = value == "true"
			default:
				userInfo.Attributes[key] = value
			}
		}

		users = append(users, userInfo)
	}

	return &ListUsersResponse{
		Users:     users,
		NextToken: resp.PaginationToken,
	}, nil
}
