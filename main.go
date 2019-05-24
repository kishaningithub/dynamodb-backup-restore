package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/jessevdk/go-flags"
	"github.com/kishaningithub/dynamodb-backup-restore/models"
	"github.com/kishaningithub/dynamodb-backup-restore/services"
	"github.com/kishaningithub/dynamodb-backup-restore/utils"
	"os"
	"strings"
)

func main() {
	opts := models.Options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(0)
	}
	dynamoDB := getDynamoDbInstance(opts)
	if strings.EqualFold(opts.Mode,"backup") {
		backupService := services.NewBackup(dynamoDB, opts.TableNamePattern, opts.BackupOutputFilePath)
		backupService.Backup()
	} else {
		restoreService := services.NewRestore(dynamoDB, opts.RestoreInputFilePath)
		restoreService.Restore()
	}
}

func getDynamoDbInstance(opts models.Options) *dynamodb.DynamoDB {
	sess, err := session.NewSession()
	utils.CheckError("Unable to get AWS session", err)
	if len(opts.EndpointUrl) == 0 {
		return dynamodb.New(sess)
	} else {
		return dynamodb.New(sess, aws.NewConfig().WithEndpoint(opts.EndpointUrl))
	}
}


