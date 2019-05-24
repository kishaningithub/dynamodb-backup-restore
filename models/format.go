package models

import "github.com/aws/aws-sdk-go/service/dynamodb"

type BackupFormat struct {
	TableName string
	Items []map[string]*dynamodb.AttributeValue
}