package timeTools

import "time"

// String2TimeYMD t 应符合 "2006-01-02" 格式
func String2TimeYMD(t string) (ret time.Time, err error) {
	ret, err = time.Parse("2006-01-02", t)
	return
}

// String2TimeYMDhms t 应符合 "2006-01-02 15:04:05" 格式
func String2TimeYMDhms(t string) (ret time.Time, err error) {
	ret, err = time.Parse("2006-01-02 15:04:05", t)
	return
}
