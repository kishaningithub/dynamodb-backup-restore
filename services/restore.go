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
	"sync"
)

type Restore interface {
	Restore()
}

type restore struct {
	dynamoDB       *dynamodb.DynamoDB
	backupFilePath string
}

func NewRestore(dynamoDB *dynamodb.DynamoDB, backupFilePath string) Restore {
	return &restore{
		dynamoDB:       dynamoDB,
		backupFilePath: backupFilePath,
	}
}

func (restore *restore) Restore() {
	itemsFromBackup := restore.getItemsFromBackup()
	var wg sync.WaitGroup
	wg.Add(len(itemsFromBackup))
	uiprogress.Start()
	for _, item := range itemsFromBackup {
		go func(item models.BackupFormat) {
			defer wg.Done()
			bar := utils.GetProgressBar(item.TableName, len(item.Items))
			restore.writeItems(item, func() {
				bar.Incr()
			})
		}(item)
	}
	wg.Wait()
}

func (restore *restore) getItemsFromBackup() []models.BackupFormat {
	var items []models.BackupFormat
	file, err := os.Open(restore.backupFilePath)
	utils.CheckError("Opening file failed", err)
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&items)
	utils.CheckError("Error while decoding backup file", err)
	return items
}

func (restore *restore) writeItems(backup models.BackupFormat, onComplete func()) {
	for i, item := range backup.Items {
		putItemInput := &dynamodb.PutItemInput{
			Item:      item,
			TableName: aws.String(backup.TableName),
		}
		_, err := restore.dynamoDB.PutItem(putItemInput)
		utils.CheckError(fmt.Sprintf("put item failed for item %v at %v position for table %v", item, i, backup.TableName), err)
		onComplete()
	}
}
