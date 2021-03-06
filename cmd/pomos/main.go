package main

import (
	"flag"
	"log"
	"os"

	"github.com/gmidorii/pomoslack"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	SQLiteFile string `yaml:"sqlite_file"`
	Slack      Slack
}

type Slack struct {
	Token   string `yaml:"token"`
	Channel string `yaml:"channel"`
}

func parse(f string) (pomoslack.Config, error) {
	yf, err := os.Open(f)
	if err != nil {
		return pomoslack.Config{}, err
	}
	defer yf.Close()

	decoder := yaml.NewDecoder(yf)

	var config Config
	if err := decoder.Decode(&config); err != nil {
		return pomoslack.Config{}, err
	}

	return pomoslack.Config{
		SQLiteFile: config.SQLiteFile,
		Slack: pomoslack.Slack{
			Token:   config.Slack.Token,
			Channel: config.Slack.Channel,
		},
	}, nil
}

func main() {
	config := flag.String("c", "~/.config/pomoslack/config.yml", "config file")
	flag.Parse()

	c, err := parse(*config)
	if err != nil {
		log.Fatalln(err)
	}
	if err := pomoslack.Run(c); err != nil {
		log.Fatalln(err)
	}
}
