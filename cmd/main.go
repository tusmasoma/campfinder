package main

import (
	"flag"
	"log"

	"github.com/joho/godotenv"
	"github.com/tusmasoma/campfinder/pkg/server"
)

func main() {
	var addr string
	// .envファイルから環境変数を読み込む
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	flag.StringVar(&addr, "addr", ":8083", "tcp host:port to connect")
	flag.Parse()

	server.Serve(addr)
}
