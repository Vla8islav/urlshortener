package configuration

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

type Options struct {
	ServerAddress      string `env:"SERVER_ADDRESS"`
	ShortenerBaseURL   string `env:"BASE_URL"`
	FileStoragePath    string `env:"FILE_STORAGE_PATH"`
	DBConnectionString string `env:"DATABASE_DSN"`
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

	if mergeInto.FileStoragePath == "" && newValues.FileStoragePath != "" {
		mergeInto.FileStoragePath = newValues.FileStoragePath
	}

	if mergeInto.DBConnectionString == "" && newValues.DBConnectionString != "" {
		mergeInto.DBConnectionString = newValues.DBConnectionString
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
	flag.StringVar(&opt.FileStoragePath, "f", "/tmp/short-url-db.json", "local file storage path")
	flag.StringVar(&opt.DBConnectionString, "d",
		"postgresql://postgres:pass123@localhost:5432/urlshortener",
		"Postgres database connection string")
	flag.Parse()
	return opt
}
