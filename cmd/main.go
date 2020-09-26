package main

import (
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
}

func main() {
	var config config
	if err := envconfig.Process("", &config); err != nil {
		log.Fatal().Err(err).Msg("failed reading config")
	}

	http := web.New(":8080", log.Logger.With().Str("processor", "http").Logger())

	http.Run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	if err := http.Stop(); err != nil {
		log.Error().Err(err).Msg("failed to shutdown http server")
	}



}