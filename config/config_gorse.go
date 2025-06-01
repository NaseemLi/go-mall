package config

type Gorse struct {
	Enable bool   `yaml:"enable"`
	Addr   string `yaml:"addr"`
	ApiKey string `yaml:"api_key"`
}
