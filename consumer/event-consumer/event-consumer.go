package event_consumer

import (
	"log"
	"time"

	"github.com/aibeksarsembayev/tbot-articles-no-extr-pkg/config"
	"github.com/aibeksarsembayev/tbot-articles-no-extr-pkg/events"
	apifetcher "github.com/aibeksarsembayev/tbot-articles-no-extr-pkg/tools/api-fetcher"
	digest "github.com/aibeksarsembayev/tbot-articles-no-extr-pkg/tools/digest"
	"go.uber.org/zap"
)

type Consumer struct {
	conf         config.Config
	lg           *zap.Logger
	fetcher      events.Fetcher
	processor    events.Processor
	apifetcher   apifetcher.APIFetcher
	digestsender digest.DigestSender
	batchSize    int
}

func New(conf config.Config, lg *zap.Logger, fetcher events.Fetcher, processor events.Processor, apif apifetcher.APIFetcher, ds digest.DigestSender, batchSize int) Consumer {
	return Consumer{
		conf:         conf,
		lg:           lg,
		fetcher:      fetcher,
		processor:    processor,
		apifetcher:   apif,
		digestsender: ds,
		batchSize:    batchSize,
	}
}

func (c Consumer) Start() error {
	// init api fetcher
	go c.apifetcher.Fetch()

	// init periodical article digest sender
	go c.digestsender.SendByWeekday()

	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())
			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err := c.handleEvents(gotEvents); err != nil {
			c.lg.Sugar().Error(err)
			continue
		}
	}
}

func (c *Consumer) handleEvents(events []events.Event) error {
	for _, event := range events {
		// log.Printf("got new event: %s", event.Text)
		c.lg.Sugar().Infof("got new event: %s", event.Text)
		if err := c.processor.Process(event); err != nil {
			// log.Printf("can't handle event: %s", err.Error())
			c.lg.Sugar().Errorf("can't handle event: %s", err.Error())
			continue
		}
	}
	return nil
}
