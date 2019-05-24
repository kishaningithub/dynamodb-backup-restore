package services

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gosuri/uiprogress"
	"github.com/kishaningithub/dynamodb-backup-restore/models"
	"github.com/kishaningithub/dynamodb-backup-restore/utils"
	"os"
	"regexp"
	"sync"
)

type Backup interface {
	Backup()
}

type backup struct {
	dynamoDB         *dynamodb.DynamoDB
	tableNamePattern string
	backupFilePath   string
}

func NewBackup(dynamoDB *dynamodb.DynamoDB, tableNamePattern string, backupFilePath string) Backup {
	return &backup{
		dynamoDB:         dynamoDB,
		tableNamePattern: tableNamePattern,
		backupFilePath:   backupFilePath,
	}
}

func (backup *backup) Backup() {
	tables := backup.findTables()
	backupItems := make(chan models.BackupFormat, len(tables))
	var wg sync.WaitGroup
	wg.Add(len(tables))
	uiprogress.Start()
	for _, table := range tables {
		go func(table string) {
			defer wg.Done()
			items := backup.fetchItems(table)
			backupItem := models.BackupFormat{
				TableName: table,
				Items:     items,
			}
			backupItems <- backupItem
		}(table)
	}
	wg.Wait()
	close(backupItems)
	var items []models.BackupFormat
	for backupItem := range backupItems {
		items = append(items, backupItem)
	}
	backup.writeBackupJSON(backup.backupFilePath, items)
	fmt.Printf("Backup written to file %s successfully!! 🎉 🎉", backup.backupFilePath)
}

func (backup *backup) findTables() []string {
	input := &dynamodb.ListTablesInput{}
	result, err := backup.dynamoDB.ListTables(input)
	utils.CheckError("Unable to fetch list of tables", err)
	var tableNames []string
	r, err := regexp.Compile(backup.tableNamePattern)
	utils.CheckError("Invalid table name pattern", err)
	for _, tableName := range result.TableNames {
		if r.MatchString(*tableName) {
			tableNames = append(tableNames, *tableName)
		}
	}
	return tableNames
}

func (backup *backup) fetchItems(tableName string) []map[string]*dynamodb.AttributeValue {
	output, err := backup.dynamoDB.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	utils.CheckError("unable to fetch total count of records", err)
	bar := backup.getProgressBar(tableName, int(*output.Table.ItemCount))
	utils.CheckError("unable to set max count of progress bar", err)
	var items []map[string]*dynamodb.AttributeValue
	backup.scan(tableName, func(value map[string]*dynamodb.AttributeValue) {
		items = append(items, value)
		bar.Incr()
	})
	return items
}

func (backup *backup) getProgressBar(tableName string, itemCount int) *uiprogress.Bar {
	return uiprogress.AddBar(int(itemCount)).
		AppendCompleted().AppendElapsed().
		PrependFunc(func(b *uiprogress.Bar) string {
		return fmt.Sprintf("%s %d/%d",tableName, b.Current(), itemCount)
	})
}

func (backup *backup) writeBackupJSON(outputFilePath string, items []models.BackupFormat) {
	f, err := os.Create(outputFilePath)
	utils.CheckError("unable write fetchItems file", err)
	encoder := json.NewEncoder(f)
	err = encoder.Encode(items)
	utils.CheckError("Writing fetchItems JSON failed", err)
}

func (backup *backup) scan(tableName string, itemConsumer func(map[string]*dynamodb.AttributeValue)) {
	params := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}
	result, err := backup.dynamoDB.Scan(params)
	utils.CheckError("Query API call failed", err)
	for _, item := range result.Items {
		itemConsumer(item)
	}
	lastEvaluatedKey := result.LastEvaluatedKey
	for lastEvaluatedKey != nil {
		params := &dynamodb.ScanInput{
			TableName:         aws.String(tableName),
			ExclusiveStartKey: lastEvaluatedKey,
		}
		result, err := backup.dynamoDB.Scan(params)
		utils.CheckError("Query API call failed", err)
		for _, item := range result.Items {
			itemConsumer(item)
		}
		lastEvaluatedKey = result.LastEvaluatedKey
	}
}
