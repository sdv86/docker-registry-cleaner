package main

import (
	"log"

	"github.com/BurntSushi/toml"
)

type Config struct {
	DrHost string
	DrPort string
	DrUser string
	DrPass string
	DrRegs []string
	DrImgCount int
}

func ReadConfig() *Config {
	var c Config
	if _, err := toml.DecodeFile("config.toml", &c); err != nil {
		log.Println(err)
	}
	return &c
}
