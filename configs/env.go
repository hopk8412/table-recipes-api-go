package configs

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func EnvMongoURI() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	return os.Getenv("MONGO_URI")
}

func EnvUserInfoURI() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	return os.Getenv("KC_USERINFO_ENDPOINT")
}

func AllowedOrigins() []string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	return strings.Split(os.Getenv("CORS_ALLOWED_LIST"), ",")
}
