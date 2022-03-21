package main

import "uranus/test/model"

func main() {
	model.GlobalDB.AutoMigrate(&model.Flight{})
}
