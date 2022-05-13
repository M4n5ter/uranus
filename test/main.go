package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println(time.Now().AddDate(0, 0, 1))
	fmt.Println(time.Now().AddDate(0, 0, -1))
}
