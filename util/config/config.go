package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Database struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		Db       string `json:"db"`
		User     string `json:"user"`
		Password string `json:"password"`
	} `json:"database"`
	Host        string `json:"host"`
	MainPort    string `json:"mainPort"`
	CommentPort string `json:"commentPort"`
	UserPort    string `json:"userPort"`
	RelojPort   string `json:"relojPort"`
}

func LoadConfiguration(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}

func GetPort(port string) string {
	var buffer bytes.Buffer
	buffer.WriteString(":")
	buffer.WriteString(port)
	return buffer.String()
}

func GenerateStringPostgresDbConnection(configFile string) string {
	conf := LoadConfiguration(configFile)
	var buffer bytes.Buffer
	buffer.WriteString("postgres://")
	buffer.WriteString(conf.Database.User)
	buffer.WriteString(":")
	buffer.WriteString(conf.Database.Password)
	buffer.WriteString("@")
	buffer.WriteString(conf.Database.Host)
	buffer.WriteString(":")
	buffer.WriteString(conf.Database.Port)
	buffer.WriteString("/")
	buffer.WriteString(conf.Database.Db)
	buffer.WriteString("?sslmode=disable")
	return buffer.String()
}

func GenerateStringMysqlDbConnection(configFile string) string {
	conf := LoadConfiguration(configFile)
	var buffer bytes.Buffer
	buffer.WriteString(conf.Database.User)
	buffer.WriteString(":")
	buffer.WriteString(conf.Database.Password)
	buffer.WriteString("@tcp(")
	buffer.WriteString(conf.Database.Host)
	buffer.WriteString(":")
	buffer.WriteString(conf.Database.Port)
	buffer.WriteString(")/")
	buffer.WriteString(conf.Database.Db)
	return buffer.String()
}
