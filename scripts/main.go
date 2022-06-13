package main

import (
	"fmt"
	"time"
)

func main() {
	today := time.Now()
	today, _ = time.Parse("2006-01-02", today.Format("2006-01-02"))
	todayString := today.Format("2006-01-02 15:04:05")
	lastDate := today.AddDate(0, 0, 2)
	lastDateString := lastDate.Format("2006-01-02 15:04:05")
	fmt.Println(todayString)
	fmt.Println(lastDateString)
}
