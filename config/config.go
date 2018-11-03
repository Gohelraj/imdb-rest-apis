package config

import (
	"github.com/BurntSushi/toml"
)

var path = "C:/Users/gohel/go/src/imdb_rest_api/config.toml"

type Service struct {
	Cockroach Cockroach
}

type Cockroach struct {
	Host        string
	Port     string
	User string
	DbName      string
	Password      string
	Dialect string
}

var Conf Service
var isLoaded = false

func Load() (Service, error) {
	if (isLoaded == true) {
		return Conf, nil
	}
	_, err := toml.DecodeFile(path, &Conf)
	isLoaded = true
	return Conf, err
}
