package cognito

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

// GroupInfo 用户组信息
type GroupInfo struct {
	GroupName   string
	Description string
}

// AddUserToGroup 将用户添加到组
func (c *Client) AddUserToGroup(username, groupName string) error {
	input := &cognitoidentityprovider.AdminAddUserToGroupInput{
		UserPoolId: aws.String(c.config.UserPoolID),
		Username:   aws.String(username),
		GroupName:  aws.String(groupName),
	}

	_, err := c.cognitoClient.AdminAddUserToGroup(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("添加用户到组失败: %w", err)
	}

	return nil
}

// RemoveUserFromGroup 从组中移除用户
func (c *Client) RemoveUserFromGroup(username, groupName string) error {
	input := &cognitoidentityprovider.AdminRemoveUserFromGroupInput{
		UserPoolId: aws.String(c.config.UserPoolID),
		Username:   aws.String(username),
		GroupName:  aws.String(groupName),
	}

	_, err := c.cognitoClient.AdminRemoveUserFromGroup(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("从组中移除用户失败: %w", err)
	}

	return nil
}

// ListUserGroups 列出用户所属的组
func (c *Client) ListUserGroups(username string) ([]*GroupInfo, error) {
	input := &cognitoidentityprovider.AdminListGroupsForUserInput{
		UserPoolId: aws.String(c.config.UserPoolID),
		Username:   aws.String(username),
	}

	resp, err := c.cognitoClient.AdminListGroupsForUser(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("获取用户组列表失败: %w", err)
	}

	groups := make([]*GroupInfo, 0, len(resp.Groups))
	for _, group := range resp.Groups {
		groupInfo := &GroupInfo{
			GroupName: *group.GroupName,
		}
		if group.Description != nil {
			groupInfo.Description = *group.Description
		}
		groups = append(groups, groupInfo)
	}

	return groups, nil
}
