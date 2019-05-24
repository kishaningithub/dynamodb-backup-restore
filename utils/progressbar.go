package utils

import (
	"fmt"
	"github.com/gosuri/uiprogress"
)

func GetProgressBar(tableName string, itemCount int) *uiprogress.Bar {
	return uiprogress.AddBar(int(itemCount)).
		AppendCompleted().AppendElapsed().
		PrependFunc(func(b *uiprogress.Bar) string {
			return fmt.Sprintf("%s %d/%d", tableName, b.Current(), itemCount)
		})
}
