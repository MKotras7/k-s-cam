package config

import (
	"encoding/json"
	"io"
	"os"

	"github.com/joho/godotenv"
)

type Server struct {
	Server_Url  string `json:"URL"`
	Server_Name string `json:"NAME"`
}

type CaptureConfig struct {
	CaptureDirectory string   `json:"CAPTURE_DIRECTORY"`
	CaptureHosts     []Server `json:"CAPTURE_HOSTS"`
	IntervalMS       int64    `json:"INTERVAL_MS"`
}

type DeleteConfig struct {
	DeleteMS   int64 `json:"DELETE_MS"` // Delete all files older than this many milliseconds
	IntervalMS int64 `json:"INTERVAL_MS"`
}

type Config struct {
	CaptureConfig CaptureConfig `json:"CAPTURE"`
	DeleteConfig  DeleteConfig  `json:"DELETE"`
}

func getEnvVar(key string) string {
	err := godotenv.Load("./.env")
	if err != nil {
		panic("godotenv.Load failed")
	}
	value := os.Getenv(key)
	return value
}

func LoadConfig() (*Config, error) {
	filename := getEnvVar("CONFIG_PATH")
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(byteValue, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
