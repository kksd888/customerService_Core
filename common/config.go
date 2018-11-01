package common

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// 配置文件
type GinConfig struct {
	GoModel    string `yaml:"GoModel"`
	Port       string `yaml:"Port"`
	Mongodb    string `yaml:"Mongodb"`
	Multimedia string `yaml:"Multimedia"`
	AiSemantic string `yaml:"AiSemantic"`
}

var AppConfig GinConfig

// 读取配置
func NewGinConfig() {
	yamlFile, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		log.Fatal("未找到配置文件conf.yaml", err)
	}
	if err := yaml.Unmarshal(yamlFile, &AppConfig); err != nil {
		log.Fatal("配置文件格式错误", err)
	}

	// 打印配置信息
	log.Info("模式：", AppConfig.GoModel)
	log.Info("开放端口：", AppConfig.Port)
	log.Info("数据库地址：", AppConfig.Mongodb)
	log.Info("多媒体资源地址：", AppConfig.Multimedia)
	log.Info("AI接口地址：", AppConfig.AiSemantic)
}
