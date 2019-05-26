package models

type Options struct {
	TableName            string `short:"t" long:"table-name" description:"Table name"`
	TableNamePattern     string `short:"p" long:"table-name-pattern" description:"Table name pattern"`
	Mode                 string `short:"m" long:"mode" description:"Mode of operation (backup,restore)" required:"true"`
	InputBackupFilePath  string `short:"i" long:"input-backup-file" description:"Input backup file path"`
	OutputBackupFilePath string `short:"o" long:"output-backup-file" description:"Output backup file path"`
	EndpointUrl          string `short:"e" long:"endpoint-url" description:"Endpoint url of destination dynamodb instance (Very useful for operating with local dynamodb instance)"`
}

func (options Options) GetBackupFilePath() string {
	if len(options.InputBackupFilePath) > 0 {
		return options.InputBackupFilePath
	}
	if len(options.OutputBackupFilePath) > 0 {
		return options.OutputBackupFilePath
	}
	return ""
}
