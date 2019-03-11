package common

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"strings"
)

// 配置文件
type GinConfig struct {
	GoMode         string `yaml:"GoMode"`
	Port           string `yaml:"Port"`
	Mongodb        string `yaml:"Mongodb"`
	Multimedia     string `yaml:"Multimedia"`
	AiSemantic     string `yaml:"AiSemantic"`
	AutomaticReply bool   `yaml:"AutomaticReply"`
}

var AppConfig GinConfig

// 读取配置
func NewGinConfig() {
	var err error
	AppConfig, err = LoadConf("conf.yaml")
	if err != nil {
		logrus.Fatalln("未找到配置文件conf.yaml ", err)
	}

	// 打印配置信息
	logrus.Info("模式：", AppConfig.GoMode)
	logrus.Info("开放端口：", AppConfig.Port)
	logrus.Info("数据库地址：", AppConfig.Mongodb)
	logrus.Info("多媒体资源地址：", AppConfig.Multimedia)
	logrus.Info("AI接口地址：", AppConfig.AiSemantic)
	logrus.Info("AI自动回复是否启动：", AppConfig.AutomaticReply)
	if AppConfig.GoMode == "debug" {
		// /debug/pprof
		go func() { _ = http.ListenAndServe(":8899", nil) }()
		logrus.Info("pprof运行中：0.0.0.0:8899/debug/pprof")
	}
}

// 读取配置
func LoadConf(confPath string) (GinConfig, error) {
	var conf GinConfig

	// 环境变量
	viper.AutomaticEnv()
	viper.SetDefault("PRODUCTION", false)
	isProduction := viper.GetBool("PRODUCTION")

	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("gorush")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if confPath != "" {
		content, err := ioutil.ReadFile(confPath)

		if err != nil {
			return conf, err
		}

		if err := viper.ReadConfig(bytes.NewBuffer(content)); err != nil {
			return conf, err
		}
	} else {
		logrus.Fatalln("config is err")
	}

	if !isProduction {
		// 测试环境
		conf.GoMode = "debug"
	} else {
		// 生产环境
		conf.GoMode = "release"
	}

	// App
	conf.Port = viper.GetString(conf.GoMode + ".Port")
	conf.Mongodb = viper.GetString(conf.GoMode + ".Mongodb")
	conf.Multimedia = viper.GetString(conf.GoMode + ".Multimedia")
	conf.AiSemantic = viper.GetString(conf.GoMode + ".AiSemantic")
	conf.AutomaticReply = viper.GetBool(conf.GoMode + ".AutomaticReply")

	return conf, nil
}
