package util

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/BurntSushi/toml"
)

type tomlConfig struct {
	BaseDir string `toml:"base_dir"`
}

var (
	cfg  *tomlConfig
	once sync.Once
)

func Config() *tomlConfig {
	configPath := "./config/conf.toml"
	once.Do(func() {
		filePath, err := filepath.Abs(configPath)
		if err != nil {
			panic(err)
		}
		fmt.Printf("parse toml file once. filePath: %s\n", filePath)
		if _, err := toml.DecodeFile(filePath, cfg); err != nil {
			panic(err)
		}
	})
	return cfg
}
