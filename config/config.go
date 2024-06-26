package config

import (
	"fmt"
	"io/ioutil"
	"log"

	"multiple-k8s-informer/controller"

	"github.com/go-yaml/yaml"
)

// TODO: 配置文件

var SysConfig *Config

type Config struct {
	MaxReQueueTime int                  `json:"maxRequeueTime" yaml:"maxRequeueTime"`
	Clusters       []controller.Cluster `json:"clusters" yaml:"clusters"`
}

func NewConfig() *Config {
	return &Config{}
}

func loadConfigFile(path string) []byte {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println(err)
		return nil
	}
	return b
}

func LoadConfig(path string) (*Config, error) {
	config := NewConfig()
	if b := loadConfigFile(path); b != nil {

		err := yaml.Unmarshal(b, config)
		if err != nil {
			return nil, err
		}

		fmt.Println(config)
		return config, err
	} else {
		return nil, fmt.Errorf("load config file error...")
	}

}
