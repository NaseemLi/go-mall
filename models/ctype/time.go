package ctype

import "time"

const TimeFormat = "2006-01-02 15:04:05"

type LocalTime time.Time

func (t *LocalTime) UnmarshalJSON(data []byte) (err error) {
	// 空值不进行解析
	if len(data) == 2 {
		*t = LocalTime(time.Time{})
		return
	}

	// 指定解析的格式
	now, err := time.Parse(`"`+TimeFormat+`"`, string(data))
	*t = LocalTime(now.Local().Add(-8 * time.Hour))
	return
}
