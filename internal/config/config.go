package config

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

const configFilename = ".gatorconfig.json"

type Config struct {
	DbURL    string `json:"db_url"`
	Username string `json:"current_user_name"`
}

func getConfigFilepath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, configFilename), nil
}

func write(cfg Config) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	configFile, err := getConfigFilepath()
	if err != nil {
		return err
	}

	f, err := os.OpenFile(configFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}

func Read() (*Config, error) {
	configFile, err := getConfigFilepath()
	if err != nil {
		return &Config{}, err
	}

	f, err := os.Open(configFile)
	if err != nil {
		return &Config{}, err
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return &Config{}, err
	}

	err = f.Close()
	if err != nil {
		return &Config{}, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return &Config{}, err
	}

	return &config, nil
}

func (cfg *Config) SetUser(username string) error {
	cfg.Username = username
	return write(*cfg)
}
