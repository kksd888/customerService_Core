package common

// 配置文件
type GinConfig struct {
	GoModel    string `yaml:"GoModel"`
	Port       string `yaml:"Port"`
	Mongodb    string `yaml:"Mongodb"`
	AiSemantic string `yaml:"AiSemantic"`
}
