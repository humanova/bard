// (c) 2020 Emir Erbasan (humanova)
// MIT License, see LICENSE for more details

package config

import "github.com/tkanos/gonfig"

type Config struct {
	PostDirectory    string
	ListenPrefixPath string
	Port             string
}

func GetConfig(configPath string, config *Config) error {
	err := gonfig.GetConf(configPath, config)
	if err != nil {
		return err
	}

	return nil
}
