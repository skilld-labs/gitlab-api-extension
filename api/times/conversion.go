package times

import (
	"strconv"
	"strings"
)

func HumanTimeConversion(seconds int64, styleType string, limitKey string, separator string) (humanTime string) {
	timeMap := [][]interface{}{[]interface{}{"mo", "month", int64(18144000)}, []interface{}{"w", "week", int64(604800)}, []interface{}{"d", "day", int64(86400)}, []interface{}{"h", "hours", int64(3600)}, []interface{}{"m", "minutes", int64(60)}, []interface{}{"s", "second", int64(1)}}
	unitMax := map[string]int{"month": 0, "week": 1, "day": 2, "hour": 3, "minute": 4, "second": 5}
	unitKey := map[bool]int{true: 0, false: 1}[styleType == "short"]
	var t int64
	for n := unitMax[limitKey]; n < 6; n++ {
		limit := timeMap[n][2].(int64)
		if seconds >= limit {
			t = seconds / limit
			humanTime += separator + strconv.FormatInt(t, 10) + timeMap[n][unitKey].(string)
			seconds -= t * limit
		}
	}
	return strings.TrimLeft(humanTime, separator)
}
