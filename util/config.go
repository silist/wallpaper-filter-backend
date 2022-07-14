package util

import (
	"fmt"
	"log"
	"path/filepath"
	"sync"

	"github.com/BurntSushi/toml"
)

const configPath = "./config/conf.toml"

type TomlConfig struct {
	BaseDir     string `toml:"base_dir"`
	DownloadDir string `toml:"download_dir"`
	Addr        string `toml:"addr"`
}

var (
	cfg  TomlConfig
	once sync.Once
)

func Config() TomlConfig {
	once.Do(func() {
		filePath, err := filepath.Abs(configPath)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("parse toml file once. filePath: %s\n", filePath)
		if _, err := toml.DecodeFile(filePath, &cfg); err != nil {
			log.Fatal(err)
		}
	})
	return cfg
}
