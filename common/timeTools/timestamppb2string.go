package timeTools

import "google.golang.org/protobuf/types/known/timestamppb"

func Timestamppb2TimeStringYMDhms(timestamp *timestamppb.Timestamp) string {
	return timestamp.AsTime().Local().Format("2006-01-02 15:04:05")
}

func Timestamppb2TimeStringYMD(timestamp *timestamppb.Timestamp) string {
	return timestamp.AsTime().Local().Format("2006-01-02")
}
