package utilx

import "time"

const (
	dateTimeLayout = "2006-01-02 15:04:05"
	dateLayout     = "2006-01-02"
)

var (
	Loc0 = time.FixedZone("UTC", 0)
)

// 当前日期
func CurrentDate() string {
	return time.Now().Format(dateLayout)
}

// 当前时间
func CurrentTime() string {
	return time.Now().Format(dateTimeLayout)
}
