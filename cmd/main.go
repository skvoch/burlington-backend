package main

import (
	"github.com/skvoch/burlington-backend/tree/master/internal/repository/couchbase"
	"github.com/skvoch/burlington-backend/tree/master/internal/service"
	"github.com/skvoch/burlington-backend/tree/master/internal/web"
	"github.com/rs/zerolog/log"
	"github.com/kelseyhightower/envconfig"
	"os"
	"os/signal"
)


type config struct {
	LogLevel       string `envconfig:"log_level" default:"debug"`
	Listen         string `envconfig:"listen" default:":8080"`
	ListenInternal string `envconfig:"listen_internal" default:":8000"`

	CouchbaseHost      string `envconfig:"COUCHBASE_HOST" default:"127.0.0.1"`
	CouchbaseUsername  string `envconfig:"COUCHBASE_USERNAME" default:"Administrator"`
	CouchbasePassword  string `envconfig:"COUCHBASE_PASSWORD" default:"password"`

}

func main() {
	var config config
	if err := envconfig.Process("", &config); err != nil {
		log.Fatal().Err(err).Msg("failed reading config")
	}

	repository, err := couchbase.New(config.CouchbaseUsername, config.CouchbasePassword, config.CouchbaseHost)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init couchbase repository")
	}
	svc := service.New(&service.Opts{
		Logger: log.Logger.With().Str("processor", "http").Logger(),
		Repo: repository,
	})
	http := web.New(config.Listen, log.Logger.With().Str("processor", "http", ).Logger(),svc)
	http.Run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	if err := http.Stop(); err != nil {
		log.Error().Err(err).Msg("failed to shutdown http server")
	}

}