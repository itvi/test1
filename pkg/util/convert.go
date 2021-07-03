package util

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// String2Int convert string to int
func String2Int(s string) int {
	out, err := strconv.Atoi(s)
	if err != nil {
		fmt.Printf("convert string '%s' to int error: %v", s, err)
	}
	return out
}

// IP2Int convert IP to int64
func IP2Int(ip string) int64 {
	if len(ip) == 0 {
		return 0
	}
	bits := strings.Split(ip, ".")
	if len(bits) < 4 {
		return 0
	}

	b0 := String2Int(bits[0])
	b1 := String2Int(bits[1])
	b2 := String2Int(bits[2])
	b3 := String2Int(bits[3])

	var sum int64
	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)
	return sum
}

// Int2IP convert int64 to IP string
func Int2IP(i int64) string {
	b0 := strconv.FormatInt((i>>24)&0xff, 10)
	b1 := strconv.FormatInt((i>>16)&0xff, 10)
	b2 := strconv.FormatInt((i>>8)&0xff, 10)
	b3 := strconv.FormatInt((i & 0xff), 10)
	return b0 + "." + b1 + "." + b2 + "." + b3
}

// Slice2String convert slice to string
func Slice2String(s []interface{}) string {
	//return strings.Join(s, ",")

	out, err := json.Marshal(s)
	if err != nil {
		return ""
	}
	return string(out)
}
