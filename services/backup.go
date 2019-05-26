package services

import (
	"bufio"
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
	dynamoDB *dynamodb.DynamoDB
	options  models.Options
}

func NewBackup(dynamoDB *dynamodb.DynamoDB, options models.Options) Backup {
	return &backup{
		dynamoDB: dynamoDB,
		options:  options,
	}
}

func (backup *backup) Backup() {
	utils.PrintInfo("Starting backup...")
	f, err := os.Create(backup.options.GetBackupFilePath())
	utils.CheckError("unable create backup file", err)
	bufferedWriter := bufio.NewWriter(f)
	encoder := json.NewEncoder(bufferedWriter)
	backupHeader := backup.getBackupHeader(backup.findTables())
	err = encoder.Encode(backupHeader)
	utils.CheckError("Error occurred when encoding header", err)
	var wg sync.WaitGroup
	wg.Add(backupHeader.GetNoOfTables())
	uiprogress.Start()
	for _, table := range backupHeader.GetTables() {
		bar := utils.GetProgressBar(table, backupHeader.GetTableInfo(table).TotalRecords)
		go func(table string, bar *uiprogress.Bar) {
			defer wg.Done()
			backup.backupTable(table, bar, encoder)
		}(table, bar)
	}
	wg.Wait()
	utils.CheckError("unable write to backup file", err)
	err = bufferedWriter.Flush()
	utils.CheckError("unable write to backup file", err)
	err = f.Close()
	utils.CheckError("unable close backup file", err)
	uiprogress.Stop()
	utils.PrintInfo(fmt.Sprintf("Backup written to file %s successfully!! 🎉 🎉", backup.options.GetBackupFilePath()))
}

func (backup *backup) findTables() []string {
	if len(backup.options.TableName) > 0 {
		return []string{backup.options.TableName}
	}
	input := dynamodb.ListTablesInput{}
	result, err := backup.dynamoDB.ListTables(&input)
	utils.CheckError("Unable to fetch list of tables", err)
	var tableNames []string
	r, err := regexp.Compile(backup.options.TableNamePattern)
	utils.CheckError("Invalid table name pattern", err)
	for _, tableName := range result.TableNames {
		if r.MatchString(*tableName) {
			tableNames = append(tableNames, *tableName)
		}
	}
	return tableNames
}

func (backup *backup) getBackupHeader(tables []string) models.BackupHeader {
	tableInfo := make(map[string]models.TableMetaData)
	for _, table := range tables {
		totalRecords := backup.findNoOfRecords(table)
		if totalRecords == 0 {
			utils.PrintInfo(fmt.Sprintf("Skipping table %s as it is empty...", table))
			continue
		}
		tableInfo[table] = models.TableMetaData{
			TotalRecords: totalRecords,
		}
	}
	return models.BackupHeader{
		TableInfo: tableInfo,
	}
}

func (backup *backup) backupTable(tableName string, bar *uiprogress.Bar, encoder *json.Encoder) {
	backup.scan(tableName, func(value map[string]*dynamodb.AttributeValue) {
		backupData := models.BackupRecord{
			TableName: tableName,
			Item:      value,
		}
		err := encoder.Encode(backupData)
		utils.CheckError("unable to write to backup file", err)
		bar.Incr()
	})
}

func (backup *backup) findNoOfRecords(tableName string) int {
	output, err := backup.dynamoDB.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	utils.CheckError("unable to fetch total count of records", err)
	return int(*output.Table.ItemCount)
}

func (backup *backup) scan(tableName string, itemConsumer func(map[string]*dynamodb.AttributeValue)) {
	params := dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}
	result, err := backup.dynamoDB.Scan(&params)
	utils.CheckError("Query API call failed", err)
	for _, item := range result.Items {
		itemConsumer(item)
	}
	lastEvaluatedKey := result.LastEvaluatedKey
	for lastEvaluatedKey != nil {
		params := dynamodb.ScanInput{
			TableName:         aws.String(tableName),
			ExclusiveStartKey: lastEvaluatedKey,
		}
		result, err := backup.dynamoDB.Scan(&params)
		utils.CheckError("Query API call failed", err)
		for _, item := range result.Items {
			itemConsumer(item)
		}
		lastEvaluatedKey = result.LastEvaluatedKey
	}
}
