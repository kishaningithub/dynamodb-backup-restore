package services

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gammazero/workerpool"
	"github.com/gosuri/uiprogress"
	"github.com/kishaningithub/dynamodb-backup-restore/models"
	"github.com/kishaningithub/dynamodb-backup-restore/utils"
	"os"
)

type Restore interface {
	Restore()
}

type restore struct {
	dynamoDB     *dynamodb.DynamoDB
	options      models.Options
	progressBars map[string]*uiprogress.Bar
}

func NewRestore(dynamoDB *dynamodb.DynamoDB, options models.Options) Restore {
	return &restore{
		dynamoDB:     dynamoDB,
		options:      options,
		progressBars: make(map[string]*uiprogress.Bar),
	}
}

func (restore *restore) Restore() {
	utils.PrintInfo("Starting restore...")
	uiprogress.Start()
	file, err := os.Open(restore.options.GetBackupFilePath())
	utils.CheckError("Opening file failed", err)
	decoder := json.NewDecoder(bufio.NewReader(file))
	header := restore.getHeader(decoder)
	for table, tableMetaData := range header.TableInfo {
		restore.progressBars[table] = utils.GetProgressBar(table, tableMetaData.TotalRecords)
	}
	wp := workerpool.New(header.GetNoOfTables())
	restore.getItemsFromBackup(decoder, func(record models.BackupRecord) {
		wp.Submit(func() {
			restore.writeItem(record)
			restore.incrementProgressBar(record)
		})
	})
	wp.StopWait()
	err = file.Close()
	utils.CheckError("Error while closing backup file", err)
	uiprogress.Stop()
	utils.PrintInfo("Restore completed successfully!! 🎉 🎉")
}

func (restore *restore) incrementProgressBar(backupFormat models.BackupRecord) {
	bar := restore.progressBars[backupFormat.TableName]
	bar.Incr()
}

func (restore *restore) getHeader(decoder *json.Decoder) models.BackupHeader {
	var header models.BackupHeader
	err := decoder.Decode(&header)
	utils.CheckError("Error while decoding backup file header", err)
	return header
}

func (restore *restore) getItemsFromBackup(decoder *json.Decoder, itemConsumer func(models.BackupRecord)) {
	for decoder.More() {
		var item models.BackupRecord
		err := decoder.Decode(&item)
		utils.CheckError("Error while decoding backup file", err)
		itemConsumer(item)
	}
}

func (restore *restore) writeItem(backup models.BackupRecord) {
	putItemInput := dynamodb.PutItemInput{
		Item:     backup.Item,
		TableName: aws.String(backup.TableName),
	}
	_, err := restore.dynamoDB.PutItem(&putItemInput)
	utils.CheckError(fmt.Sprintf("put item failed for item %v for table %v", backup.Item, backup.TableName), err)
}
