package common

// 配置文件
type GinConfig struct {
	RunModel   string `yaml:"RunModel"`
	Port       string `yaml:"Port"`
	Mongodb    string `yaml:"Mongodb"`
	AiSemantic string `yaml:"AiSemantic"`
}
