package timeTools

import "time"

// Time2TimeYMD000 将时间的年月日以外的内容都设置为零值，例如：2022-05-15 13:49:51.6522294 +0800 CST m=+0.002638901 会被转化成 2022-05-15 00:00:00 +0800 CST
func Time2TimeYMD000(t time.Time) (time.Time, error) {
	t, err := time.ParseInLocation("2006-01-02", t.Format("2006-01-02"), t.Location())
	return t, err
}
