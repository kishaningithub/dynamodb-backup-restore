package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/jessevdk/go-flags"
	"gopkg.in/cheggaaa/pb.v1"
	"log"
	"os"
)

type options struct {
	TableName string `short:"t" long:"table-name" description:"Name of the dynamo db table" required:"true"`
	Mode      string `short:"m" long:"mode" description:"Mode of operation (backup,restore)" required:"true"`
	BackupOutputFile string `short:"o" long:"output" description:"Output file for backup"`
	RestoreInputFile string `short:"i" long:"input" description:"Input file for restore"`
}

var items []map[string]*dynamodb.AttributeValue

func main() {
	opts := options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(0)
	}
	sess, err := session.NewSession()
	dynamoDB := getDynamoDbSession(sess)
	if opts.Mode == "backup" {
		backup(dynamoDB, opts.TableName, opts.BackupOutputFile)
	} else {
		restore(opts.RestoreInputFile, opts.TableName, dynamoDB)
	}
}

func backup(dynamoDB *dynamodb.DynamoDB, tableName string, outputFilePath string) {
	output, err := dynamoDB.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	checkError("unable to fetch total count of records", err)
	bar := pb.StartNew(int(*output.Table.ItemCount))
	result := doFirstScan(tableName, dynamoDB, bar)
	doSubSequentScan(result, tableName, dynamoDB, bar)
	writeBackupJSON(outputFilePath)
}

func writeBackupJSON(outputFilePath string) {
	f, err := os.Create(outputFilePath)
	checkError("unable write backup file", err)
	encoder := json.NewEncoder(f)
	err = encoder.Encode(items)
	checkError("Writing backup JSON failed", err)
}

func incrementProgressBarBy(count int, pb *pb.ProgressBar) {
	for i := 0; i < count; i++ {
		pb.Increment()
	}
}

func doFirstScan(tableName string, srcDynamoDB *dynamodb.DynamoDB, bar *pb.ProgressBar) *dynamodb.ScanOutput {
	params := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}
	result, err := srcDynamoDB.Scan(params)
	checkError("Query API call failed", err)
	items = append(items, result.Items...)
	incrementProgressBarBy(len(result.Items), bar)
	return result
}

func doSubSequentScan(firstScanResult *dynamodb.ScanOutput, tableName string, dynamoDB *dynamodb.DynamoDB, bar *pb.ProgressBar) {
	lastEvaluatedKey := firstScanResult.LastEvaluatedKey
	for lastEvaluatedKey != nil {
		params := &dynamodb.ScanInput{
			TableName:         aws.String(tableName),
			ExclusiveStartKey: lastEvaluatedKey,
		}
		result, err := dynamoDB.Scan(params)
		checkError("Query API call failed", err)
		items = append(items, result.Items...)
		incrementProgressBarBy(len(result.Items), bar)
		lastEvaluatedKey = result.LastEvaluatedKey
	}
}

func restore(backupFile string, tableName string, dynamoDB *dynamodb.DynamoDB) {
	items := getItemsFromBackup(backupFile)
	bar := pb.StartNew(len(items))
	writeItems(items, tableName,dynamoDB, bar)
}

func getItemsFromBackup(backupFile string) []map[string]*dynamodb.AttributeValue {
	var items []map[string]*dynamodb.AttributeValue
	file, err := os.Open(backupFile)
	checkError("Opening file failed", err)
	decoder := json.NewDecoder(file)
	err = decoder.Decode(items)
	checkError("Error while decoding backup file", err)
	return items
}

func writeItems(items []map[string]*dynamodb.AttributeValue, tableName string, dynamoDB *dynamodb.DynamoDB, bar *pb.ProgressBar) {
	for i, item := range items {
		putItemInput := &dynamodb.PutItemInput{
			Item:      item,
			TableName: aws.String(tableName),
		}
		_, err := dynamoDB.PutItem(putItemInput)
		bar.Increment()
		checkError(fmt.Sprintf("put item failed for item %v at %v position", item, i), err)
	}
}

func getDynamoDbSession(sess *session.Session) *dynamodb.DynamoDB {
	return dynamodb.New(sess)
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, ": ", err.Error())
	}
}
