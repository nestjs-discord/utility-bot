package env

import "github.com/joho/godotenv"

func init() {
	_ = godotenv.Load()
}
