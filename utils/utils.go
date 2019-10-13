package utils

import "time"

func MakeTimestamp() int {
	return int(time.Now().UnixNano() / int64(time.Millisecond))
}