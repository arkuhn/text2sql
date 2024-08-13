package config

import (
        "encoding/json"
        "os"
        "path/filepath"
)

const configFileName = ".text2sql_config.json"

func getConfigPath() string {
        homeDir, err := os.UserHomeDir()
        if err != nil {
                return configFileName
        }
        return filepath.Join(homeDir, configFileName)
}

func GetConfig(key string) string {
        configPath := getConfigPath()
        data, err := os.ReadFile(configPath)
        if err != nil {
                return ""
        }

        var config map[string]string
        err = json.Unmarshal(data, &config)
        if err != nil {
                return ""
        }

        return config[key]
}

func SetConfig(key, value string) error {
        configPath := getConfigPath()
        var config map[string]string

        data, err := os.ReadFile(configPath)
        if err == nil {
                err = json.Unmarshal(data, &config)
                if err != nil {
                        config = make(map[string]string)
                }
        } else {
                config = make(map[string]string)
        }

        config[key] = value

        data, err = json.MarshalIndent(config, "", "  ")
        if err != nil {
                return err
        }

        return os.WriteFile(configPath, data, 0644)
}
