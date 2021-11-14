package store

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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

func (storage *DynamoDB) Select(ctx context.Context, dest interface{}, stmt string, params ...interface{}) error {
	input := &dynamodb.ExecuteStatementInput{
		Statement: aws.String(stmt),
	}

	if len(params) > 0 {
		pp, err := attributevalue.MarshalList(params)
		if err != nil {
			return fmt.Errorf("failed to marshal params: %w", err)
		}
		input.Parameters = pp
	}

	result, err := storage.client.ExecuteStatement(ctx, input)
	if err != nil {
		return fmt.Errorf("faield to execute statement: %w", err)
	}

	if err := attributevalue.UnmarshalListOfMaps(result.Items, dest); err != nil {
		return fmt.Errorf("failed to unmarshal map: %w", err)
	}

	return nil
}
