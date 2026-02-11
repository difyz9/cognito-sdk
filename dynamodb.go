package cognito

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// DynamoDBClient DynamoDB客户端封装
type DynamoDBClient struct {
	client *dynamodb.Client
}

// NewDynamoDBClient 创建DynamoDB客户端
func (c *Client) NewDynamoDBClient() *DynamoDBClient {
	return &DynamoDBClient{
		client: dynamodb.NewFromConfig(c.GetAWSConfig()),
	}
}

// PutItem 创建或更新项目
func (db *DynamoDBClient) PutItem(ctx context.Context, tableName string, item interface{}) error {
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %w", err)
	}

	_, err = db.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	})
	if err != nil {
		return fmt.Errorf("failed to put item: %w", err)
	}

	return nil
}

// GetItem 获取单个项目
func (db *DynamoDBClient) GetItem(ctx context.Context, tableName string, key map[string]types.AttributeValue, result interface{}) error {
	output, err := db.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       key,
	})
	if err != nil {
		return fmt.Errorf("failed to get item: %w", err)
	}

	if output.Item == nil {
		return fmt.Errorf("item not found")
	}

	err = attributevalue.UnmarshalMap(output.Item, result)
	if err != nil {
		return fmt.Errorf("failed to unmarshal item: %w", err)
	}

	return nil
}

// UpdateItem 更新项目
func (db *DynamoDBClient) UpdateItem(ctx context.Context, tableName string, key map[string]types.AttributeValue, 
	updateExpression string, expressionAttributeValues map[string]types.AttributeValue, 
	expressionAttributeNames map[string]string) error {
	
	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(tableName),
		Key:                       key,
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeValues: expressionAttributeValues,
	}

	if len(expressionAttributeNames) > 0 {
		input.ExpressionAttributeNames = expressionAttributeNames
	}

	_, err := db.client.UpdateItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	return nil
}

// DeleteItem 删除项目
func (db *DynamoDBClient) DeleteItem(ctx context.Context, tableName string, key map[string]types.AttributeValue) error {
	_, err := db.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key:       key,
	})
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	return nil
}

// Query 查询项目
func (db *DynamoDBClient) Query(ctx context.Context, tableName string, keyConditionExpression string,
	expressionAttributeValues map[string]types.AttributeValue, results interface{}) error {
	
	output, err := db.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(tableName),
		KeyConditionExpression:    aws.String(keyConditionExpression),
		ExpressionAttributeValues: expressionAttributeValues,
	})
	if err != nil {
		return fmt.Errorf("failed to query items: %w", err)
	}

	err = attributevalue.UnmarshalListOfMaps(output.Items, results)
	if err != nil {
		return fmt.Errorf("failed to unmarshal items: %w", err)
	}

	return nil
}

// Scan 扫描表
func (db *DynamoDBClient) Scan(ctx context.Context, tableName string, results interface{}, 
	filterExpression *string, expressionAttributeValues map[string]types.AttributeValue) error {
	
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	if filterExpression != nil {
		input.FilterExpression = filterExpression
		input.ExpressionAttributeValues = expressionAttributeValues
	}

	output, err := db.client.Scan(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to scan table: %w", err)
	}

	err = attributevalue.UnmarshalListOfMaps(output.Items, results)
	if err != nil {
		return fmt.Errorf("failed to unmarshal items: %w", err)
	}

	return nil
}

// BatchWriteItem 批量写入项目
func (db *DynamoDBClient) BatchWriteItem(ctx context.Context, tableName string, items []interface{}) error {
	const batchSize = 25 // DynamoDB限制每次最多25个项目

	for i := 0; i < len(items); i += batchSize {
		end := i + batchSize
		if end > len(items) {
			end = len(items)
		}

		batch := items[i:end]
		writeRequests := make([]types.WriteRequest, len(batch))

		for j, item := range batch {
			av, err := attributevalue.MarshalMap(item)
			if err != nil {
				return fmt.Errorf("failed to marshal item: %w", err)
			}

			writeRequests[j] = types.WriteRequest{
				PutRequest: &types.PutRequest{
					Item: av,
				},
			}
		}

		_, err := db.client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				tableName: writeRequests,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to batch write items: %w", err)
		}
	}

	return nil
}

// CreateTable 创建DynamoDB表
func (db *DynamoDBClient) CreateTable(ctx context.Context, input *dynamodb.CreateTableInput) error {
	_, err := db.client.CreateTable(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	return nil
}

// DeleteTable 删除DynamoDB表
func (db *DynamoDBClient) DeleteTable(ctx context.Context, tableName string) error {
	_, err := db.client.DeleteTable(ctx, &dynamodb.DeleteTableInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		return fmt.Errorf("failed to delete table: %w", err)
	}
	return nil
}

// ListTables 列出所有表
func (db *DynamoDBClient) ListTables(ctx context.Context) ([]string, error) {
	output, err := db.client.ListTables(ctx, &dynamodb.ListTablesInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to list tables: %w", err)
	}
	return output.TableNames, nil
}

// DescribeTable 获取表信息
func (db *DynamoDBClient) DescribeTable(ctx context.Context, tableName string) (*types.TableDescription, error) {
	output, err := db.client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to describe table: %w", err)
	}
	return output.Table, nil
}

// TableExists 检查表是否存在
func (db *DynamoDBClient) TableExists(ctx context.Context, tableName string) (bool, error) {
	_, err := db.client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		// 检查是否是表不存在的错误
		if err.Error() != "" {
			// 简单的错误检查，表不存在时返回 false
			return false, nil
		}
		return false, fmt.Errorf("failed to check table existence: %w", err)
	}
	return true, nil
}
