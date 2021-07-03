package util

import (
	"strings"
	"time"
)

type JsonDate time.Time

// DateTime format the string like `20140319102030.000000+480` to yyyy-mm-dd hh:mm:ss
func DateTime(s string) string { // 20140319102030.000000+480
	str := strings.Split(s, ".")[0]
	date := str[0:4] + "-" + str[4:6] + "-" + str[6:8]
	time := str[8:10] + ":" + str[10:12] + ":" + str[12:14]
	return date + " " + time // 2014-03-19 10:20:30
}

func (j *JsonDate) UnmarshalJSON(data []byte) error {
	local, err := time.ParseInLocation("\"2006-01-02\"", string(data), time.Local)
	*j = JsonDate(local)
	return err
}

// func (j JsonDate) MarshalJSON() ([]byte, error) {
// 	var stamp = fmt.Sprintf("\"%s\"", time.Time(j).Format("2006-01-02"))
// 	return []byte(stamp), nil
// }
