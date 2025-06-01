package sql

import (
	"fmt"
	"strings"
)

func OrderRevert(idList []string, column string) (sql string) {
	var list []string
	for _, v := range idList {
		list = append(list, fmt.Sprintf(" %s = %s desc", column, v))
	}
	sql = strings.Join(list, ", ")
	return
}
