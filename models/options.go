package models

type Options struct {
	TableNamePattern     string `short:"t" long:"table-name" description:"Table name pattern"`
	Mode                 string `short:"m" long:"mode" description:"Mode of operation (backup,restore)" required:"true"`
	BackupOutputFilePath string `short:"o" long:"output" description:"Output file for backup"`
	RestoreInputFilePath string `short:"i" long:"input" description:"Input file for restore"`
	EndpointUrl          string `short:"e" long:"endpoint-url" description:"Endpoint url of destination dynamodb instance (Very useful for restoring data into local dynamodb instance)"`
}
