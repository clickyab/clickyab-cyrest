package accutils

import (
	"common/utils"
	"fmt"
	"modules/balance/libs"
	"strings"
)

// AccessToQuery create a query filter for using in lists to filter not accessed items
func AccessToQuery(a libs.Access, read bool, start int) (string, []interface{}) {
	if a.IsOwner() {
		return " TRUE ", nil
	}
	var ids []int64
	var all = false
	if read {
		all, ids = a.VisibleObjectIDs()
	} else {
		all, ids = a.WritableObjectIDs()
	}

	if all {
		return " TRUE ", nil
	}

	idIn := make([]interface{}, len(ids))
	for i := range ids {
		idIn[i] = ids[i]
	}
	s, in := utils.BuildPgPlaceHolder(start, idIn...)

	filter := fmt.Sprintf(" IN (%s) ", strings.Join(s, ","))
	return filter, in
}
