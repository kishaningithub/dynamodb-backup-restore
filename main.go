package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/jessevdk/go-flags"
	"github.com/kishaningithub/dynamodb-backup-restore/models"
	"github.com/kishaningithub/dynamodb-backup-restore/services"
	"github.com/kishaningithub/dynamodb-backup-restore/utils"
	"strings"
)

func main() {
	opts := getApplicationOptions()
	dynamoDB := getDynamoDbInstance(opts)
	if strings.EqualFold(opts.Mode, "backup") {
		backupService := services.NewBackup(dynamoDB, opts)
		backupService.Backup()
	} else {
		restoreService := services.NewRestore(dynamoDB, opts)
		restoreService.Restore()
	}
}

func getApplicationOptions() models.Options {
	opts := models.Options{}
	_, err := flags.NewParser(&opts, flags.HelpFlag|flags.PassDoubleDash).Parse()
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok {
			if flagsErr.Type == flags.ErrHelp {
				utils.CheckError("", flagsErr)
			}
		}
		utils.CheckError("invalid options", err)
	}
	return opts
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
