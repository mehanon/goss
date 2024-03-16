package shrimp

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Profile struct {
	Tag  string   `yaml:"tag"`
	URLs []string `yaml:"url"`
}

type Config struct {
	Tokens    []string       `yaml:"tokens"`
	Source    string         `yaml:"source"`
	ChatId    int64          `yaml:"chat-id"`
	Dir       string         `yaml:"dir"`
	TagTopics map[string]int `yaml:"tag-topics"`
	Profiles  []Profile      `yaml:"profiles"`
	filename  string
}

func (cfg *Config) Update() error {
	buff, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(cfg.filename, buff, 0755)
}

func ConfigFromFile(filename string) (*Config, error) {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(buf, &cfg); err != nil {
		return nil, err
	}
	if cfg.TagTopics == nil {
		cfg.TagTopics = make(map[string]int)
	}
	cfg.filename = filename
	return &cfg, nil
}
