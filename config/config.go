package config

import (
	"github.com/augustazz/camellia/constants"

	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var serverConf *ServerConfig

func init() {
	configPath := os.Getenv("configPath")
	if configPath == "" {
		configPath = constants.DefaultConfigPath
	}
	loadConfig(configPath)
}

func GetSrvConfig() *ServerConfig {
	return serverConf
}

func GetServiceConfig() []*Service {
	return serverConf.Services
}

func GetConnConfig() ConnConfig {
	return serverConf.Conn
}

type ServerConfig struct {
	App      AppConfig  `yaml:"app"`
	Web      WebConfig  `yaml:"web"`
	Log      LogConfig  `yaml:"log"`
	Conn     ConnConfig `yaml:"conn"`
	Services []*Service `yaml:"services"`
}

type AppConfig struct {
	Name     string `yaml:"name"`
	Profiles string `yaml:"profiles"`
}

type WebConfig struct {
	Port int `yaml:"port"`
}

type ConnConfig struct {
	AuthTimeout  uint32 `yaml:"auth-timeout"`
	IdleTimeout  uint32 `yaml:"idle-timeout"`
	AuthFilePath string `yaml:"auth-file-path"`
}

type LogConfig struct {
	Debug bool   `yaml:"debug"`
	Path  string `yaml:"path"`
}

type Service struct {
	Name string   `yaml:"name"`
	Url  []string `yaml:"url"`
}

//LoadConfig load yml config,if has profiles,load after default file,same config item will be rewrite
func loadConfig(path string) *ServerConfig {
	if !isYamlFileConfig(path) {
		log.Fatalln("not support file: ", path)
		return nil
	}
	if serverConf == nil {
		serverConf = &ServerConfig{}
	}
	//load default config file
	loadConfig0(path, serverConf)
	if serverConf.App.Profiles != "" {
		//load config file with profile,eg:config-test.yaml
		w := strings.LastIndexByte(path, '.')
		name, suffix := path[:w], path[w+1:]
		for _, p := range strings.Split(serverConf.App.Profiles, ",") {
			fileWithProfile := fmt.Sprintf("%s-%s.%s", name, p, suffix)
			loadConfig0(fileWithProfile, serverConf)
			log.Println("load config file success: ", fileWithProfile)
		}
	}
	return serverConf
}

func loadConfig0(path string, conf *ServerConfig) {
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
