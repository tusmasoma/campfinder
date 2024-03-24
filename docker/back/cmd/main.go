package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/tusmasoma/campfinder/docker/back/driver"
)

func main() {
	// .envファイルから環境変数を読み込む
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	driver.Run()
}
