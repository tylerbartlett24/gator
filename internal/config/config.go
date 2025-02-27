package config

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
)

const configFileName = "/.gatorconfig.json"

type Config struct {
	URL            string `json:"db_url"`
	Username       string `json:"current_user_name"`
}

func (c *Config) SetUser(username string) error {
	c.Username = username
	err := write(*c)
	if err != nil {
		return err
	}

	return nil
}

func write(cfg Config) error {
	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(filePath, os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.Encode(cfg)
	return nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	
	filePath := homeDir + configFileName
	return filePath, nil
}


func Read() Config {
	filePath, err := getConfigFilePath()
	if err != nil {
		log.Fatalf("Failed to find home directory: %v\n", err)
	}

	fileText, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v\n", err)
	}
	
	var c Config
	json.NewDecoder(bytes.NewBuffer(fileText)).Decode(&c)
	return c	
}




