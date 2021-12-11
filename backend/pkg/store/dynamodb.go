package store

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDB struct {
	client *dynamodb.Client
}

func NewDynamoDB(cfg aws.Config) *DynamoDB {
	return &DynamoDB{
		client: dynamodb.NewFromConfig(cfg),
	}
}

func (storage *DynamoDB) Execute(ctx context.Context, stmt string, params ...interface{}) error {
	pp, err := attributevalue.MarshalList(params)
	if err != nil {
		return fmt.Errorf("failed to marshal params: %w", err)
	}

	if _, err := storage.client.ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
		Statement:  aws.String(stmt),
		Parameters: pp,
	}); err != nil {
		return fmt.Errorf("faield to execute statement: %w", err)
	}

	return nil
}

func (storage *DynamoDB) Query(ctx context.Context, dest interface{}, stmt string, params ...interface{}) error {
	items, err := storage.queryPage(ctx, stmt, nil, params...)
	if err != nil {
		return fmt.Errorf("failed to select: %w", err)
	}

	if err := attributevalue.UnmarshalListOfMaps(items, dest); err != nil {
		return fmt.Errorf("failed to unmarshal map: %w", err)
	}

	return nil
}

func (storage *DynamoDB) queryPage(ctx context.Context, stmt string, nextToken *string, params ...interface{}) ([]map[string]types.AttributeValue, error) {
	input := &dynamodb.ExecuteStatementInput{
		Statement: aws.String(stmt),
	}

	if len(params) > 0 {
		pp, err := attributevalue.MarshalList(params)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal params: %w", err)
		}
		input.Parameters = pp
	}

	if nextToken != nil {
		input.NextToken = nextToken
	}

	result, err := storage.client.ExecuteStatement(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("faield to execute statement: %w", err)
	}

	if result.NextToken == nil {
		return result.Items, nil
	}

	nextPage, err := storage.queryPage(ctx, stmt, result.NextToken, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to get next page: %w", err)
	}

	return append(result.Items, nextPage...), nil
}
