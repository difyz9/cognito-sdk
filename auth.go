package cognito

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

// RegisterResponse 注册响应
type RegisterResponse struct {
	UserSub string
	Success bool
}

// SignUpResponse 用户自注册响应
type SignUpResponse struct {
	UserSub         string
	UserConfirmed   bool
	CodeDeliveryDetails *CodeDeliveryDetails
}

// CodeDeliveryDetails 验证码发送详情
type CodeDeliveryDetails struct {
	Destination     string // 邮箱地址（部分隐藏）
	DeliveryMedium  string // EMAIL 或 SMS
	AttributeName   string // email 或 phone_number
}

// LoginResponse 登录响应
type LoginResponse struct {
	IDToken      string
	AccessToken  string
	RefreshToken string
	ExpiresIn    int32
}

// RefreshTokenResponse 刷新Token响应
type RefreshTokenResponse struct {
	IDToken     string
	AccessToken string
	ExpiresIn   int32
}

// Register 注册新用户（使用AdminCreateUser，自动确认用户）
func (c *Client) Register(email, password, username string) (*RegisterResponse, error) {
	if username == "" {
		username = email
	}

	// 使用AdminCreateUser创建用户（管理员方式，自动确认）
	input := &cognitoidentityprovider.AdminCreateUserInput{
		UserPoolId:        aws.String(c.config.UserPoolID),
		Username:          aws.String(email), // 使用邮箱作为用户名
		TemporaryPassword: aws.String(password),
		MessageAction:     types.MessageActionTypeSuppress, // 关闭默认邮件通知
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(email),
			},
			{
				Name:  aws.String("email_verified"),
				Value: aws.String("true"), // 直接标记邮箱已验证
			},
		},
	}

	resp, err := c.cognitoClient.AdminCreateUser(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("注册失败: %w", err)
	}

	// 设置永久密码
	err = c.SetPermanentPassword(email, password)
	if err != nil {
		return nil, fmt.Errorf("设置永久密码失败: %w", err)
	}

	var userSub string
	if resp.User != nil {
		for _, attr := range resp.User.Attributes {
			if *attr.Name == "sub" {
				userSub = *attr.Value
				break
			}
		}
	}

	return &RegisterResponse{
		UserSub: userSub,
		Success: true,
	}, nil
}

// SignUpUser 用户自注册（需要邮箱验证）
// 适用场景：用户自助注册，需要通过邮箱验证码确认
// 注册后需要调用 ConfirmSignUp 完成验证
func (c *Client) SignUpUser(email, password, username string) (*SignUpResponse, error) {
	if username == "" {
		username = email
	}

	input := &cognitoidentityprovider.SignUpInput{
		ClientId: aws.String(c.config.ClientID),
		Username: aws.String(email),
		Password: aws.String(password),
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(email),
			},
		},
	}

	resp, err := c.cognitoClient.SignUp(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("用户注册失败: %w", err)
	}

	response := &SignUpResponse{
		UserSub:       *resp.UserSub,
		UserConfirmed: resp.UserConfirmed,
	}

	if resp.CodeDeliveryDetails != nil {
		response.CodeDeliveryDetails = &CodeDeliveryDetails{
			Destination:    aws.ToString(resp.CodeDeliveryDetails.Destination),
			DeliveryMedium: string(resp.CodeDeliveryDetails.DeliveryMedium),
			AttributeName:  aws.ToString(resp.CodeDeliveryDetails.AttributeName),
		}
	}

	return response, nil
}

// ConfirmSignUp 确认用户注册（使用邮箱验证码）
// username: 用户名（通常是邮箱）
// confirmationCode: 用户收到的验证码
func (c *Client) ConfirmSignUp(username, confirmationCode string) error {
	input := &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         aws.String(c.config.ClientID),
		Username:         aws.String(username),
		ConfirmationCode: aws.String(confirmationCode),
	}

	_, err := c.cognitoClient.ConfirmSignUp(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("确认注册失败: %w", err)
	}

	return nil
}

// ResendConfirmationCode 重新发送验证码
func (c *Client) ResendConfirmationCode(username string) (*CodeDeliveryDetails, error) {
	input := &cognitoidentityprovider.ResendConfirmationCodeInput{
		ClientId: aws.String(c.config.ClientID),
		Username: aws.String(username),
	}

	resp, err := c.cognitoClient.ResendConfirmationCode(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("重发验证码失败: %w", err)
	}

	if resp.CodeDeliveryDetails == nil {
		return nil, nil
	}

	return &CodeDeliveryDetails{
		Destination:    aws.ToString(resp.CodeDeliveryDetails.Destination),
		DeliveryMedium: string(resp.CodeDeliveryDetails.DeliveryMedium),
		AttributeName:  aws.ToString(resp.CodeDeliveryDetails.AttributeName),
	}, nil
}

// Login 用户登录
func (c *Client) Login(username, password string) (*LoginResponse, error) {
	authParams := map[string]string{
		"USERNAME": username,
		"PASSWORD": password,
	}

	input := &cognitoidentityprovider.AdminInitiateAuthInput{
		UserPoolId: aws.String(c.config.UserPoolID),
		ClientId:   aws.String(c.config.ClientID),
		AuthFlow:   types.AuthFlowTypeAdminUserPasswordAuth,
		AuthParameters: authParams,
	}

	resp, err := c.cognitoClient.AdminInitiateAuth(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("登录失败: %w", err)
	}

	if resp.AuthenticationResult == nil {
		return nil, fmt.Errorf("登录失败: 未返回认证结果")
	}

	return &LoginResponse{
		IDToken:      *resp.AuthenticationResult.IdToken,
		AccessToken:  *resp.AuthenticationResult.AccessToken,
		RefreshToken: *resp.AuthenticationResult.RefreshToken,
		ExpiresIn:    resp.AuthenticationResult.ExpiresIn,
	}, nil
}

// RefreshToken 刷新访问令牌
func (c *Client) RefreshToken(refreshToken string) (*RefreshTokenResponse, error) {
	authParams := map[string]string{
		"REFRESH_TOKEN": refreshToken,
	}

	input := &cognitoidentityprovider.AdminInitiateAuthInput{
		UserPoolId:     aws.String(c.config.UserPoolID),
		ClientId:       aws.String(c.config.ClientID),
		AuthFlow:       types.AuthFlowTypeRefreshTokenAuth,
		AuthParameters: authParams,
	}

	resp, err := c.cognitoClient.AdminInitiateAuth(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("刷新Token失败: %w", err)
	}

	if resp.AuthenticationResult == nil {
		return nil, fmt.Errorf("刷新Token失败: 未返回认证结果")
	}

	return &RefreshTokenResponse{
		IDToken:     *resp.AuthenticationResult.IdToken,
		AccessToken: *resp.AuthenticationResult.AccessToken,
		ExpiresIn:   resp.AuthenticationResult.ExpiresIn,
	}, nil
}

// ChangePassword 修改密码（需要旧密码）
func (c *Client) ChangePassword(username, oldPassword, newPassword string) error {
	// 先登录获取AccessToken
	loginResp, err := c.Login(username, oldPassword)
	if err != nil {
		return fmt.Errorf("验证旧密码失败: %w", err)
	}

	input := &cognitoidentityprovider.ChangePasswordInput{
		PreviousPassword: aws.String(oldPassword),
		ProposedPassword: aws.String(newPassword),
		AccessToken:      aws.String(loginResp.AccessToken),
	}

	_, err = c.cognitoClient.ChangePassword(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("修改密码失败: %w", err)
	}

	return nil
}

// ResetPassword 重置密码（管理员权限，无需旧密码）
func (c *Client) ResetPassword(username, newPassword string) error {
	input := &cognitoidentityprovider.AdminSetUserPasswordInput{
		UserPoolId: aws.String(c.config.UserPoolID),
		Username:   aws.String(username),
		Password:   aws.String(newPassword),
		Permanent:  true,
	}

	_, err := c.cognitoClient.AdminSetUserPassword(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("重置密码失败: %w", err)
	}

	return nil
}

// SignOut 全局登出用户
func (c *Client) SignOut(username string) error {
	input := &cognitoidentityprovider.AdminUserGlobalSignOutInput{
		UserPoolId: aws.String(c.config.UserPoolID),
		Username:   aws.String(username),
	}

	_, err := c.cognitoClient.AdminUserGlobalSignOut(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("登出失败: %w", err)
	}

	return nil
}

// DeleteUser 删除用户
func (c *Client) DeleteUser(username string) error {
	input := &cognitoidentityprovider.AdminDeleteUserInput{
		UserPoolId: aws.String(c.config.UserPoolID),
		Username:   aws.String(username),
	}

	_, err := c.cognitoClient.AdminDeleteUser(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("删除用户失败: %w", err)
	}

	return nil
}
