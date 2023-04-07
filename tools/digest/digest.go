package articledigest

import (
	"context"
	"errors"
	"fmt"
	"time"

	tgClient "github.com/aibeksarsembayev/tbot-articles-no-extr-pkg/clients/telegram"
	"github.com/aibeksarsembayev/tbot-articles-no-extr-pkg/storage"
	"github.com/aibeksarsembayev/tbot-articles-no-extr-pkg/storage/postgres"
	"go.uber.org/zap"
)

type Sender struct {
	lg      *zap.Logger
	chatID  int
	storage *postgres.Storage
	tg      *tgClient.Client
}

func New(logger *zap.Logger, chatID int, s *postgres.Storage, tg *tgClient.Client) *Sender {
	return &Sender{
		lg:      logger,
		chatID:  chatID,
		storage: s,
		tg:      tg,
	}
}

// SendByWeekday
func (s *Sender) SendByWeekday() {
	// ticker period
	minutePeriod := 1 * time.Minute
	t := time.NewTicker(minutePeriod)

	// send schedule - saturday 6AM 6UTC

	defer t.Stop()
	for {
		select {
		case <-t.C:

			weekday := time.Now().Weekday()
			hourNow := time.Now().UTC().Hour() + 6 // UTC0 hour + 6hrs
			minuteNow := time.Now().Minute()       // now in minutes

			// fmt.Println(time.Now().Weekday(), time.Now())

			if int(weekday) == 6 && hourNow == 6 && minuteNow == 0 { // Saturday 6H:00M AM
				err := s.pushDigest()
				if err != nil {
					s.lg.Sugar().Error(err)
				}
			}
		}
	}
}

// Send digest by each 7 days. TODO: delete once 2nd version will be tested
func (s *Sender) Send() {
	// initial digest
	err := s.pushDigest()
	if err != nil {
		s.lg.Sugar().Error(err)
	}
	weekPeriod := 168 * time.Hour   // 168hrs - 7d
	t := time.NewTicker(weekPeriod) // 7 days
	defer t.Stop()
	for {
		select {
		case <-t.C:
			err = s.pushDigest()
			if err != nil {
				s.lg.Sugar().Error(err)
			}
		}
	}
}

func (s *Sender) pushDigest() error {
	articles, err := s.storage.GetLatest(context.Background())
	if err != nil && !errors.Is(err, storage.ErrNoArticles) {
		return fmt.Errorf("weekly digest: can't pull articles from db: %w", err)
	}

	msg := ""

	if err == storage.ErrNoArticles {
		msg = `<b>Еженедельная сводка</b>
	
Приветствуем Вас, друзья!

За последние 7 дней в нашей «базе знаний» новых публикаций не было.
	
Список всех статей, инструкций и материалов в «Базе знаний» доступен по ссылке – <a href="https://sber-invest.kz/knowledgebase/list">«Перечень всех статей»</a>
			
С уважением, 
Сбережения и Инвестиции`
	} else {
		msg = `<b>Еженедельная сводка</b>	
	
Приветствуем Вас, друзья! 

За последние 7 дней в нашей «базе знаний» появились новые статьи на тему налогообложения, инвестиций и финансовой грамотности. Ниже весь перечень:
`

		msgEnd := `
Список всех статей, инструкций и материалов в «Базе знаний» доступен по ссылке – <a href="https://sber-invest.kz/knowledgebase/list">«Перечень всех статей»</a>
		
С уважением, 
Сбережения и Инвестиции
	`

		for _, a := range articles {

			temp := fmt.Sprintf(`
<b>%s</b>
<a href="%s">Прочитать статью полностью></a>
`, a.Title, a.URL)
			msg += temp
		}

		msg += msgEnd
	}

	err = s.tg.SendMessageURL(s.chatID, msg)
	if err != nil {
		return fmt.Errorf("weekly digest: can't push articles via tg client: %w", err)
	}

	if len(articles) != 0 {
		s.lg.Info("weekly digest was pushed to channel")
	} else {
		s.lg.Info("weekly digest was pushed to channel but no articles")
	}

	return nil
}
