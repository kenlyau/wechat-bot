package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type StConfig struct {
	DllServer string
	Port      string
}

var Config StConfig
var templates map[string]string
var commands map[string]map[string]Command

type Command struct {
	Classify string `json:"classify"`
	Variate  string `json:"variate"`
}

func GetTemplates() map[string]string {
	return templates
}

func GetCommands() map[string]map[string]Command {
	return commands
}

func getFileList(path string) []string {
	fs, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	list := make([]string, 0)
	for _, file := range fs {
		list = append(list, file.Name())
	}
	return list
}

func SetUp() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	viper.Unmarshal(&Config)
	templates = make(map[string]string)
	commands = make(map[string]map[string]Command)

	commandFiles := getFileList("./data/commands")
	for _, fileName := range commandFiles {
		name := strings.Replace(fileName, ".json", "", -1)
		byteValue, e := ioutil.ReadFile("./data/commands/" + fileName)
		if e != nil {
			log.Fatal(e)
		}
		var nowCommand map[string]Command
		json.Unmarshal(byteValue, &nowCommand)

		commands[name] = nowCommand
	}
	templateFiles := getFileList("./data/templates")
	for _, fileName := range templateFiles {
		name := strings.Replace(fileName, ".txt", "", -1)
		byteValue, e := ioutil.ReadFile("./data/templates/" + fileName)
		if e != nil {
			log.Fatal(e)
		}
		templates[name] = string(byteValue)
	}
	log.Println(commands)
}
