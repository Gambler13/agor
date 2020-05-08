package conf

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Game struct {
	Port           int      `yaml:"port"`
	TickRate       int      `yaml:"tickRate"`
	AllowedOrigins []string `yaml:"allowedOrigins"`
	World          World    `yaml:"world"`
}

type World struct {
	Width     int `yaml:"width"`
	Height    int `yaml:"height"`
	MaxPlayer int `yaml:"maxPlayer"`
	Food      int `yaml:"food"`
}

func Load(file string) (conf Game, err error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	c := Game{}
	err = yaml.Unmarshal(data, &c)
	conf = c
	return
}
