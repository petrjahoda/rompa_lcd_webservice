package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var DatabaseType string
var DatabaseIpAddress string
var DatabaseName string
var DatabasePort string
var DatabaseLogin string
var DatabasePassword string

type Config struct {
	DatabaseType string
	IpAddress    string
	DatabaseName string
	Port         string
	Login        string
	Password     string
}

func CreateConfigIfNotExists() {
	configDirectory := filepath.Join(".", "config")
	configFileName := "config.json"
	configFullPath := strings.Join([]string{configDirectory, configFileName}, "/")

	if _, checkPathError := os.Stat(configFullPath); checkPathError == nil {
		LogDebug("MAIN", "Config file already exists")
	} else if os.IsNotExist(checkPathError) {
		LogWarning("MAIN", "Config file does not exist, creating")
		mkdirError := os.MkdirAll(configDirectory, 0777)
		if mkdirError != nil {
			LogError("MAIN", "Unable to create directory for config file: "+mkdirError.Error())
		} else {
			LogInfo("MAIN", "Directory for config file created")
			data := Config{
				DatabaseType: "mysql",
				IpAddress:    "zapsidatabase",
				DatabaseName: "zapsi2",
				Port:         "3306",
				Login:        "zapsi_uzivatel",
				Password:     "zapsi",
			}
			file, _ := json.MarshalIndent(data, "", "  ")
			writingError := ioutil.WriteFile(configFullPath, file, 0666)
			LogInfo("MAIN", "Writing data to JSON file")
			if writingError != nil {
				LogError("MAIN", "Unable to write data to config file: "+writingError.Error())
			} else {
				LogInfo("MAIN", "Data written to config file")
			}
		}
	} else {
		LogError("MAIN", "Config file does not exist")
	}
}

func LoadSettingsFromConfigFile() {
	configDirectory := filepath.Join(".", "config")
	configFileName := "config.json"
	configFullPath := strings.Join([]string{configDirectory, configFileName}, "/")
	ConfigFile := Config{}
	for len(ConfigFile.DatabaseName) == 0 {
		readFile, err := ioutil.ReadFile(configFullPath)
		if err != nil {
			LogError("MAIN", "Problem reading config file")
			var err = os.Remove(configFullPath)
			if err != nil {
				LogError("MAIN", "Problem deleting file "+configFullPath+", "+err.Error())
				break
			}
			CreateConfigIfNotExists()
		}
		err = json.Unmarshal(readFile, &ConfigFile)
		if err != nil {
			LogError("MAIN", "Problem parsing config file, deleting config file")
			var err = os.Remove(configFullPath)
			if err != nil {
				LogError("MAIN", "Problem deleting file "+configFullPath+", "+err.Error())
				break
			}
			CreateConfigIfNotExists()
		}
	}
	DatabaseType = ConfigFile.DatabaseType
	DatabaseIpAddress = ConfigFile.IpAddress
	DatabaseName = ConfigFile.DatabaseName
	DatabasePort = ConfigFile.Port
	DatabaseLogin = ConfigFile.Login
	DatabasePassword = ConfigFile.Password
}
