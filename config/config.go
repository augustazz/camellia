package config

import (
	"context"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Config interface {
	LoadConfig(ctx context.Context, profiles, path string)
	GetConfig() Config
}

var serverConf *ServerConfig

type ServerConfig struct {
	App AppConfig `yaml:"app"`
	Web WebConfig `yaml:"web"`
	Log LogConfig `yaml:"log"`
}

type AppConfig struct {
	Name string `yaml:"name"`
}

type WebConfig struct {
	Port int `yaml:"port"`
}

type LogConfig struct {
	Debug bool   `yaml:"debug"`
	Path  string `yaml:"path"`
}

//LoadConfig load yml config,if has profiles,load after default file,same config item will be rewrite
func (c *ServerConfig) LoadConfig(ctx context.Context, profiles, path string) {
	if !isYamlFileConfig(path) {
		log.Fatalln("not support file: ", path)
		return
	}
	//load default config file
	loadConfig0(path, c)
	if profiles != "" {
		//load config file with profile,eg:config-test.yaml
		w := strings.LastIndexByte(path, '.')
		name, suffix := path[:w], path[w+1:]
		for _, p := range strings.Split(profiles, ",") {
			fileWithProfile := fmt.Sprintf("%s-%s.%s", name, p, suffix)
			loadConfig0(fileWithProfile, c)
			log.Println("load config file success: ", fileWithProfile)
		}
	}
}

func (c *ServerConfig) GetConfig() Config {
	return serverConf
}

func loadConfig0(path string, conf Config) {
	open, err := os.Open(path)

	if err != nil {
		panic("start err:open config fail:" + err.Error())
	}
	defer open.Close()
	data, err := ioutil.ReadAll(open)

	if err != nil {
		panic("start err:read config fail" + err.Error())
	}

	err = yaml.Unmarshal(data, conf)

	if err != nil {
		panic("start err:bind config fail" + err.Error())
	}
}

func isYamlFileConfig(fileName string) bool {
	return strings.HasSuffix(fileName, ".yml") || strings.HasSuffix(fileName, ".yaml")
}
