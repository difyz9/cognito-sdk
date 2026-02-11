package cognito

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

// Config AWS Cognito SDK配置
type Config struct {
	UserPoolID string
	ClientID   string
	Region     string
}

// Client Cognito SDK客户端
type Client struct {
	config       Config
	awsConfig    aws.Config
	cognitoClient *cognitoidentityprovider.Client
}

// NewClient 创建新的Cognito客户端
func NewClient(cfg Config) (*Client, error) {
	if cfg.UserPoolID == "" || cfg.ClientID == "" || cfg.Region == "" {
		return nil, fmt.Errorf("配置不完整：UserPoolID、ClientID和Region都不能为空")
	}

	// 加载AWS配置
	awsConfig, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.Region),
	)
	if err != nil {
		return nil, fmt.Errorf("加载AWS配置失败: %w", err)
	}

	// 创建Cognito客户端
	cognitoClient := cognitoidentityprovider.NewFromConfig(awsConfig)

	return &Client{
		config:       cfg,
		awsConfig:    awsConfig,
		cognitoClient: cognitoClient,
	}, nil
}

// NewClientWithCredentials 使用指定的凭证创建客户端
func NewClientWithCredentials(cfg Config, accessKeyID, secretAccessKey string) (*Client, error) {
	if cfg.UserPoolID == "" || cfg.ClientID == "" || cfg.Region == "" {
		return nil, fmt.Errorf("配置不完整：UserPoolID、ClientID和Region都不能为空")
	}

	// 加载AWS配置
	awsConfig, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     accessKeyID,
				SecretAccessKey: secretAccessKey,
			}, nil
		})),
	)
	if err != nil {
		return nil, fmt.Errorf("加载AWS配置失败: %w", err)
	}

	// 创建Cognito客户端
	cognitoClient := cognitoidentityprovider.NewFromConfig(awsConfig)

	return &Client{
		config:       cfg,
		awsConfig:    awsConfig,
		cognitoClient: cognitoClient,
	}, nil
}

// GetConfig 获取配置
func (c *Client) GetConfig() Config {
	return c.config
}

// GetAWSConfig 获取AWS配置
func (c *Client) GetAWSConfig() aws.Config {
	return c.awsConfig
}
