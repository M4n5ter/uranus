package timeTools

import "google.golang.org/protobuf/types/known/timestamppb"

// Timestamppb2TimeStringYMDhms 将 *timestamppb.Timestamp 转化为 对应 time.Time 的格式化字符串("2006-01-02 15:04:05")
func Timestamppb2TimeStringYMDhms(timestamp *timestamppb.Timestamp) string {
	return timestamp.AsTime().Local().Format("2006-01-02 15:04:05")
}

// Timestamppb2TimeStringYMD 将 *timestamppb.Timestamp 转化为 对应 time.Time 的格式化字符串("2006-01-02")
func Timestamppb2TimeStringYMD(timestamp *timestamppb.Timestamp) string {
	return timestamp.AsTime().Local().Format("2006-01-02")
}
