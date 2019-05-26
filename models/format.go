package models

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type BackupHeader struct {
	TableInfo map[string]TableMetaData
}

func (backupHeader BackupHeader) GetNoOfTables() int {
	return len(backupHeader.TableInfo)
}

func (backupHeader BackupHeader) GetTables() []string {
	tables := make([]string, len(backupHeader.TableInfo))
	i := 0
	for key := range backupHeader.TableInfo {
		tables[i] = key
		i++
	}
	return tables
}

func (backupHeader BackupHeader) GetTableInfo(tableName string) TableMetaData {
	return backupHeader.TableInfo[tableName]
}

type TableMetaData struct {
	TotalRecords int
}

type BackupRecord struct {
	TableName string
	Item      map[string]*dynamodb.AttributeValue
}
