package env

import (
	"flag"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	var path string
	flag.StringVar(&path, "config", "./env/.env", "path to env file")
	flag.Parse()
	err := godotenv.Load(path)
	if err != nil {
		panic("fail to read env file")
	}
}
