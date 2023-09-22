package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/hablof/order-viewer/config"
	"github.com/hablof/order-viewer/internal/cache/inmem"
	"github.com/hablof/order-viewer/internal/consumer"
	"github.com/hablof/order-viewer/internal/database"
	"github.com/hablof/order-viewer/internal/httpcontroller"
	"github.com/hablof/order-viewer/internal/repository"
	"github.com/hablof/order-viewer/internal/service"
	"github.com/hablof/order-viewer/internal/templates"
)

const (
	postgresURL string = "postgres://order_viewer:%s@127.0.0.1:5432/orders?sslmode=disable"
)

func main() {

	cfg, err := config.NewConfig()
	if err != nil {
		log.Error().Err(err).Msg("failed to get config")
		return
	}

	if cfg.Debug {
		log.Logger.Level(zerolog.DebugLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"})
	} else {
		log.Logger.Level(zerolog.InfoLevel)
	}

	log.Info().Str("log level", log.Logger.GetLevel().String()).Send()

	log := log.Logger.With().Str("func", "main").Caller().Logger()

	mainCtx, cf := context.WithCancel(context.Background())
	defer cf()

	t, err := templates.GetTemplates()
	if err != nil {
		log.Error().Err(err).Send()
		return
	}

	pg, err := database.NewPostgres(mainCtx, cfg.GetPgURL())
	if err != nil {
		log.Error().Err(err).Msg("failed to setup database connection")
		return
	}

	log.Info().Msg("connected to database")

	r, err := repository.NewRepository(mainCtx, pg)
	if err != nil {
		log.Error().Err(err).Msg("failed to setup repository")
		return
	}

	log.Info().Msg("set up repository")

	c := inmem.NewInMemCache()

	s, err := service.NewService(c, r)
	if err != nil {
		log.Error().Err(err).Msg("failed to setup service")
		return
	}

	log.Info().Msg("set up service")

	mux := httpcontroller.GetRouter(s, t)

	server := http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      mux,
		ReadTimeout:  cfg.HTTPTimeout,
		WriteTimeout: cfg.HTTPTimeout,
	}

	sc, err := consumer.RegisterStanClient(s, cfg)
	if err != nil {
		log.Error().Err(err).Msg("failed to setup STAN client")
		return
	}

	log.Info().Msg("set up STAN connection")

	go server.ListenAndServe()
	go sc.RunNconsumers(mainCtx, cfg)

	terminationChannel := make(chan os.Signal, 1)
	signal.Notify(terminationChannel, os.Interrupt, syscall.SIGTERM)

	<-terminationChannel
	cf()

	log.Info().Msg("terminating")

}
