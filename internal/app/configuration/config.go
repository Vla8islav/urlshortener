package configuration

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

type Options struct {
	ServerAddress    string `env:"SERVER_ADDRESS"`
	ShortenerBaseURL string `env:"BASE_URL"`
}

var instance *Options

func ReadFlags() *Options {
	if instance == nil {
		cmdOptions := getCmdOptions()
		envOptions := getEnvOptions()

		finalOptions := Options{}
		// env options are the priority
		mergeOptions(&finalOptions, envOptions)
		mergeOptions(&finalOptions, cmdOptions)
		instance = &finalOptions
	}
	return instance
}

func mergeOptions(mergeInto *Options, newValues Options) {
	if mergeInto.ServerAddress == "" && newValues.ServerAddress != "" {
		mergeInto.ServerAddress = newValues.ServerAddress
	}

	if mergeInto.ShortenerBaseURL == "" && newValues.ShortenerBaseURL != "" {
		mergeInto.ShortenerBaseURL = newValues.ShortenerBaseURL
	}
}

func getEnvOptions() Options {
	var opt Options
	err := env.Parse(&opt)
	if err != nil {
		log.Fatalln(err)
	}
	return opt
}

func getCmdOptions() Options {
	opt := Options{}
	flag.StringVar(&opt.ServerAddress, "a", "localhost:8080", "port on which the server should run")
	flag.StringVar(&opt.ShortenerBaseURL, "b", "http://localhost:8080", "base url for shortened links")
	flag.Parse()
	return opt
}
