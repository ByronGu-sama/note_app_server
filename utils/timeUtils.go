package utils

import "time"

// AtoT 字符串转秒数
func AtoT(str string) (time.Duration, error) {
	duration, err := time.ParseDuration(str)
	if err != nil {
		return -1, err
	}
	return duration, nil
}
