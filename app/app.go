package app

import (
	"time"

	tgClient "github.com/aibeksarsembayev/tbot-articles-no-extr-pkg/clients/telegram"
	"github.com/aibeksarsembayev/tbot-articles-no-extr-pkg/config"
	event_consumer "github.com/aibeksarsembayev/tbot-articles-no-extr-pkg/consumer/event-consumer"
	"github.com/aibeksarsembayev/tbot-articles-no-extr-pkg/events/telegram"
	"github.com/aibeksarsembayev/tbot-articles-no-extr-pkg/storage/postgres"
	apifetcher "github.com/aibeksarsembayev/tbot-articles-no-extr-pkg/tools/api-fetcher"
	digest "github.com/aibeksarsembayev/tbot-articles-no-extr-pkg/tools/digest"
	"github.com/aibeksarsembayev/tbot-articles-no-extr-pkg/tools/logger"
)

func Start() {
	// init logger
	lg := logger.InitLogger()

	// load configs
	conf, err := config.LoadConfig()
	if err != nil {
		// fmt.Println(err)
		lg.Sugar().Error(err)
	} else {
		// fmt.Println(conf)
		lg.Sugar().Info(conf)
	}

	// init postgres db by creation pool of connection for DB
	dbpool, err := postgres.InitPostgresDBConn(&conf)
	if err != nil {
		// log.Fatalf("database: %v", err)
		lg.Sugar().Fatalf("database: %v", err)
	}
	defer dbpool.Close()

	// article repository
	s := postgres.NewDBArticleRepo(dbpool)

	// API fetcher (periodical in minutes)
	apifetcher := apifetcher.New(
		lg,
		s,
		time.Duration(conf.TgBot.APIParsePeriod))

	// tgclient and event processor start ...
	tg := tgClient.New(
		conf.TgBot.Host,
		conf.TgBot.Token)

	eventsProcessor := telegram.New(
		lg,
		tg,
		s)
	// log.Print("service started")
	lg.Info("service started")

	// weekly article digest in tgchannel
	digestsender := digest.New(
		lg,
		conf.TgBot.DigestChatID,
		s,
		tg)

	// event consumer starts ...
	consumer := event_consumer.New(
		conf,
		lg,
		eventsProcessor,
		eventsProcessor,
		apifetcher,
		digestsender,
		conf.TgBot.BatchSize)
	if err := consumer.Start(); err != nil {
		// log.Fatal("service is stopped", err)
		lg.Sugar().Fatalf("service is stopped", err)
	}
}
