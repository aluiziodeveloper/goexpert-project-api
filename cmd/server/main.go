package main

import "github.com/aluiziodeveloper/goexpert-project-api/configs"

func main() {
	config, _ := configs.LoadConfig(".")
	println(config.DBDriver)
}
