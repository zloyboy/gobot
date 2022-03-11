package config

import (
	"log"
	"os"
	"strconv"

	godotenv "github.com/joho/godotenv"
)

type Config struct {
	AnswerTout int
	ApiToken   string
	Version    string
	Notify     bool
}

func ReadConfig() *Config {
	if err := godotenv.Load("data/.env"); err != nil {
		log.Print("No .env file found")
		return nil
	}

	ANSWER_TIMEOUT := os.Getenv("ANSWER_TIMEOUT")
	tout, err := strconv.Atoi(ANSWER_TIMEOUT)
	if err != nil {
		tout = 30
	}
	API_TOKEN := os.Getenv("TELEGRAM_API_TOKEN")
	VERSION := os.Getenv("VERSION")

	return &Config{
		tout,
		API_TOKEN,
		VERSION,
		os.Getenv("NOTIFY") == "1"}
}
