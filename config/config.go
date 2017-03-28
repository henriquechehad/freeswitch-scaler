package config

import (
	"flag"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type config struct {
	ConsulServer string `yaml:"consul_servers"`
	ESLTimeout   int    `yaml:"esl_default_timeout"`
	ESLPassword  string `yaml:"esl_default_password"`
	ESLPort      int    `yaml:"esl_default_port"`
	ESLMaxTries  int    `yaml:"esl_max_tries"`
}

var Config *config

func Init() {
	cfgFile := flag.String("conf", "config.yml", "Configuration YAML file")
	src, err := ioutil.ReadFile(*cfgFile)
	if err != nil {
		log.Fatalln("Error loading config:", err)
	}

	err = yaml.Unmarshal(src, &Config)
	if err != nil {
		log.Fatalln("Error unmarshal config:", err)
	}

}
