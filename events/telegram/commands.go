package telegram

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	tgclient "github.com/aibeksarsembayev/tbot-articles-no-extr-pkg/clients/telegram"
	"github.com/aibeksarsembayev/tbot-articles-no-extr-pkg/lib/e"
	"github.com/aibeksarsembayev/tbot-articles-no-extr-pkg/storage"
)

const (
	// Main menu
	HelpCmd  = "/help"
	StartCmd = "/start"
	// Article
	ArticleCmd            = "/articles"
	ArticleDigestCmd      = "bydigest"
	CategoryCmd           = "bycategory"
	AuthorCmd             = "byauthor"
	AllArticleCmd         = "allarticles"
	LatestArticleCmd      = "latestarticles"
	LatestMonthArticleCmd = "latestmontharticles"
	ByCategoryCmd         = "category-"
	ByAuthorCmd           = "author-"
	// Status
	StatusCreatorCmd    = "creator"
	StatusAdminCmd      = "admin"
	StatusMemberCmd     = "member"
	StatusRestrictedCMD = "resticted"
	StatusLeftCMD       = "left"
	StatusBannedCMD     = "kicked"
)

func (p *Processor) doCmd(text string, meta Meta) error {
	input := ""
	callback_query := ""
	// print meta
	fmt.Println(meta)
	if meta.ChannelPost != "" || meta.EditedChannelPost != "" {
		input = fmt.Sprintf("channel_post-%s", meta.ChannelPost)
		p.lg.Sugar().Infof("got new message '%s' from '%s' channel", meta.ChannelPost, meta.Username)
	} else if meta.PollID != "" {
		input = fmt.Sprintf("channel_post-%s", meta.ChannelPost)
		p.lg.Sugar().Infof("got new message poll_id '%s' from 'unknown' channel", meta.PollID)
	} else if text != "" {
		text = strings.TrimSpace(text)
		p.lg.Sugar().Infof("got new command '%s' from '%s'", text, meta.Username)
		// group commands style correction
		if strings.HasSuffix(text, "@sber_invest_bot") {
			text = strings.TrimSuffix(text, "@sber_invest_bot")
		}
		input = text
	} else if meta.CallbackQuery != "" {
		callback_query = strings.TrimSpace(meta.CallbackQuery)
		p.lg.Sugar().Infof("got new callbackquery '%s' from '%s'", callback_query, meta.Username)
		input = callback_query
	} else if meta.Status != "" {
		p.lg.Sugar().Infof("got new status change command '%s' from '%s'", meta.Status, meta.Username)
		input = "status-" + meta.Status
	}

	switch {
	// Main menu
	case input == HelpCmd:
		return p.sendHelp(meta.ChatID)
	case input == StartCmd:
		return p.sendHello(meta.ChatID)
		// Article
	case input == ArticleCmd:
		return p.articlesFilter(meta.ChatID)
	case input == ArticleDigestCmd:
		return p.digestFilter(meta.ChatID)
	case input == AuthorCmd:
		return p.sendAuthor(meta.ChatID)
	case input == CategoryCmd:
		return p.sendCategory(meta.ChatID)
	case input == AllArticleCmd:
		return p.sendAll(meta.ChatID)
	case input == LatestArticleCmd:
		return p.sendLatest(meta.ChatID)
	case input == LatestMonthArticleCmd:
		return p.sendLatestMonth(meta.ChatID)
	case strings.HasPrefix(input, ByAuthorCmd):
		author := strings.TrimPrefix(callback_query, "author-")
		return p.sendByAuthor(meta.ChatID, author)
	case strings.HasPrefix(input, ByCategoryCmd):
		category := strings.TrimPrefix(callback_query, "category-")
		return p.sendByCategory(meta.ChatID, category)
		// Status
	case strings.HasPrefix(input, "status-"):
		// return p.tg.SendMessage(chatID, msgStatusChanged)
		return nil // TODO: make proper handle for status changes
	case strings.HasPrefix(input, "channel_post-"):
		// text := strings.TrimPrefix(input, "channel_post-")
		// return p.tg.SendMessage(chatID, text)
		return nil // TODO: make proper handle for status changes
	default:
		return p.tg.SendMessage(meta.ChatID, msgUnknownCommand)
	}
}

func (p *Processor) articlesFilter(chatID int) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send articles by category", err) }()

	mainMenu := []string{
		"по рубрикам",
		"по авторам",
		"cводка за 7/30 дней",
		"все статьи",
	}

	msg := &tgclient.SendMessage{
		ChatID:    chatID,
		Text:      "<b>Выберите подходящий фильтр</b>",
		ParseMode: "HTML",
	}

	msg.ReplyMarkup.InlineKeyboard = [][]tgclient.InlineKeyboardButton{
		{
			{Text: mainMenu[0], CallbackData: CategoryCmd},
			{Text: mainMenu[1], CallbackData: AuthorCmd},
		},
		{
			{Text: mainMenu[2], CallbackData: ArticleDigestCmd},
			{Text: mainMenu[3], CallbackData: AllArticleCmd},
		},
	}

	if err := p.tg.SendMessagePost(msg); err != nil {
		return err
	}

	return nil
}

func (p *Processor) digestFilter(chatID int) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send articles by category", err) }()

	mainMenu := []string{
		"7 дней",
		"30 дней",
	}

	msg := &tgclient.SendMessage{
		ChatID:    chatID,
		Text:      "<b>Выберите интервал выборки</b>",
		ParseMode: "HTML",
	}

	msg.ReplyMarkup.InlineKeyboard = [][]tgclient.InlineKeyboardButton{
		{
			{Text: mainMenu[0], CallbackData: LatestArticleCmd},
			{Text: mainMenu[1], CallbackData: LatestMonthArticleCmd},
		},
	}

	if err := p.tg.SendMessagePost(msg); err != nil {
		return err
	}

	return nil
}

// sendByAUthor articles ...
func (p *Processor) sendByAuthor(chatID int, author string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send articles by category", err) }()
	userID, err := strconv.Atoi(author)
	if err != nil {
		return err
	}
	articles, err := p.storage.GetByAuthor(context.Background(), userID)
	if err != nil {
		return err
	}

	for _, a := range articles {
		createdTime := a.CreatedAt
		text := fmt.Sprintf(`%s
<b>%s</b>
<a href="%s">Читать...</a>`, createdTime.Format("2006-01-02"), a.Title, a.URL)
		if err := p.tg.SendMessage(chatID, text); err != nil {
			return err
		}
	}

	// send main menu buttons at the end of articles
	if err := p.sendHelpMenu(chatID); err != nil {
		return err
	}

	return nil
}

// sendByCategory articles ...
func (p *Processor) sendByCategory(chatID int, category string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send articles by category", err) }()
	categoryID, err := strconv.Atoi(category)
	if err != nil {
		return err
	}
	articles, err := p.storage.GetByCategory(context.Background(), categoryID)
	if err != nil {
		return err
	}

	for _, a := range articles {
		createdTime := a.CreatedAt
		text := fmt.Sprintf(`%s
<b>%s</b>
<a href="%s">Читать...</a>`, createdTime.Format("2006-01-02"), a.Title, a.URL)
		if err := p.tg.SendMessage(chatID, text); err != nil {
			return err
		}
	}

	// send main menu buttons at the end of articles
	if err := p.sendHelpMenu(chatID); err != nil {
		return err
	}

	return nil
}

// sendAll of articles ...
func (p *Processor) sendAll(chatID int) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send all articles", err) }()

	articles, err := p.storage.GetAllAPI(context.Background())
	if err != nil && !errors.Is(err, storage.ErrNoArticles) {
		return err
	}

	for _, a := range articles {
		createdTime := a.CreatedAt
		text := fmt.Sprintf(`%s
<b>%s</b>
<a href="%s">Читать...</a>`, createdTime.Format("2006-01-02"), a.Title, a.URL)
		if err := p.tg.SendMessage(chatID, text); err != nil {
			return err
		}
	}

	// send main menu buttons at the end of articles
	if err := p.sendHelpMenu(chatID); err != nil {
		return err
	}

	return nil
}

// sendLatest of articles by 7 days ...
func (p *Processor) sendLatest(chatID int) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send latest articles", err) }()

	articles, err := p.storage.GetLatest(context.Background())

	if err != nil && !errors.Is(err, storage.ErrNoArticles) {
		return err
	}

	msg := ""

	// if errors.Is(err, storage.ErrNoArticles) {
	if err == storage.ErrNoArticles {
		msg = `<b>Сводка за последние 7 дней</b>

За последние 7 дней в нашей «базе знаний» новых публикаций не было.
	
Список всех статей, инструкций и материалов в «Базе знаний» доступен по ссылке – <a href="https://sber-invest.kz/knowledgebase/list">«Перечень всех статей»</a>
			
С уважением, 
Сбережения и Инвестиции`
	} else {
		msg = `<b>Сводка за последние 7 дней</b>	

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

	if err := p.tg.SendMessageURL(chatID, msg); err != nil {
		return err
	}

	// send main menu buttons at the end
	if err := p.sendHelpMenuLatest(chatID); err != nil {
		return err
	}

	return nil
}

// sendLatest of articles by 30 days ...
func (p *Processor) sendLatestMonth(chatID int) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send latest articles", err) }()

	articles, err := p.storage.GetLatestMonth(context.Background())
	if err != nil && !errors.Is(err, storage.ErrNoArticles) {
		return err
	}

	if err == storage.ErrNoArticles {
		msg := `<b>Сводка за последние 30 дней</b>

За последние 30 дней в нашей «базе знаний» новых публикаций не было.
	
Список всех статей, инструкций и материалов в «Базе знаний» доступен по ссылке – <a href="https://sber-invest.kz/knowledgebase/list">«Перечень всех статей»</a>
			
С уважением, 
Сбережения и Инвестиции`
		if err := p.tg.SendMessageURL(chatID, msg); err != nil {
			return err
		}
	} else {
		msg := `<b>Сводка за последние 30 дней</b>

Список всех новых статей на тему налогообложения, инвестиций и финансовой грамотности за последние 30 дней ниже:
`

		for _, a := range articles {
			temp := fmt.Sprintf(`
<b>%s</b>
<a href="%s">Прочитать статью полностью></a>
`, a.Title, a.URL)
			msg += temp
		}
		msg += `
С уважением, 
Сбережения и Инвестиции`
		if err := p.tg.SendMessageURL(chatID, msg); err != nil {
			return err
		}
	}

	// send main menu buttons at the end
	if err := p.sendHelpMenuLatestMonth(chatID); err != nil {
		return err
	}

	return nil
}

// sendAuthor of articles ...
func (p *Processor) sendAuthor(chatID int) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send categories", err) }()

	authors, err := p.storage.GetAuthor(context.Background())
	if err != nil {
		return err
	}

	msg := &tgclient.SendMessage{
		ChatID:    chatID,
		Text:      "<b>Выберите автора</b>",
		ParseMode: "HTML",
	}

	for i := 0; i < len(authors); i++ {
		preslice := make([]tgclient.InlineKeyboardButton, 1)
		msg.ReplyMarkup.InlineKeyboard = append(msg.ReplyMarkup.InlineKeyboard, preslice)
	}

	for i, c := range authors {
		author := fmt.Sprintf("%s %s", c.AuthorName, c.AuthorSurname)
		callbackData := fmt.Sprintf("author-%v", c.UserID)
		msg.ReplyMarkup.InlineKeyboard[i][0] = tgclient.InlineKeyboardButton{Text: author, CallbackData: callbackData}
	}

	if err := p.tg.SendMessagePost(msg); err != nil {
		return err
	}
	return nil
}

// sendCategory for article filter...
func (p *Processor) sendCategory(chatID int) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send categories", err) }()

	categories, err := p.storage.GetCategory(context.Background())
	if err != nil {
		return err
	}

	msg := &tgclient.SendMessage{
		ChatID:    chatID,
		Text:      "<b>Выберите рубрику</b>",
		ParseMode: "HTML",
	}

	for i := 0; i < len(categories); i++ {
		preslice := make([]tgclient.InlineKeyboardButton, 1)
		msg.ReplyMarkup.InlineKeyboard = append(msg.ReplyMarkup.InlineKeyboard, preslice)
	}

	for i, c := range categories {
		categoryName := c.Category
		callbackData := fmt.Sprintf("category-%v", c.ID)
		msg.ReplyMarkup.InlineKeyboard[i][0] = tgclient.InlineKeyboardButton{Text: categoryName, CallbackData: callbackData}
	}

	if err := p.tg.SendMessagePost(msg); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendHelp(chatID int) error {
	err := p.tg.SendMessage(chatID, msgHelp)
	if err != nil {
		return err
	}
	// send main menu buttons at the end of articles
	if err := p.sendHelpMenu(chatID); err != nil {
		return err
	}
	return nil
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

// to send menu after articles to better navigation
func (p *Processor) sendHelpMenu(chatID int) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send articles by category", err) }()

	mainMenu := []string{
		"help",
		"База знаний",
	}

	msg := &tgclient.SendMessage{
		ChatID:    chatID,
		Text:      "<b>Menu</b>",
		ParseMode: "HTML",
	}

	msg.ReplyMarkup.InlineKeyboard = [][]tgclient.InlineKeyboardButton{
		{
			{Text: mainMenu[0], CallbackData: HelpCmd},
			{Text: mainMenu[1], CallbackData: ArticleCmd},
		},
	}

	if err := p.tg.SendMessagePost(msg); err != nil {
		return err
	}

	return nil
}

// to send menu after digest 7 days for better navigation
func (p *Processor) sendHelpMenuLatest(chatID int) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send articles by category", err) }()

	mainMenu := []string{
		"30 дней",
		"help",
		"База знаний",
	}

	msg := &tgclient.SendMessage{
		ChatID:    chatID,
		Text:      "<b>Menu</b>",
		ParseMode: "HTML",
	}

	msg.ReplyMarkup.InlineKeyboard = [][]tgclient.InlineKeyboardButton{
		{
			{Text: mainMenu[0], CallbackData: LatestMonthArticleCmd},
			{Text: mainMenu[1], CallbackData: HelpCmd},
			{Text: mainMenu[2], CallbackData: ArticleCmd},
		},
	}

	if err := p.tg.SendMessagePost(msg); err != nil {
		return err
	}

	return nil
}

// to send menu after digest 30 days for better navigation
func (p *Processor) sendHelpMenuLatestMonth(chatID int) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send articles by category", err) }()

	mainMenu := []string{
		"7 дней",
		"help",
		"База знаний",
	}

	msg := &tgclient.SendMessage{
		ChatID:    chatID,
		Text:      "<b>Menu</b>",
		ParseMode: "HTML",
	}

	msg.ReplyMarkup.InlineKeyboard = [][]tgclient.InlineKeyboardButton{
		{
			{Text: mainMenu[0], CallbackData: LatestArticleCmd},
			{Text: mainMenu[1], CallbackData: HelpCmd},
			{Text: mainMenu[2], CallbackData: ArticleCmd},
		},
	}

	if err := p.tg.SendMessagePost(msg); err != nil {
		return err
	}

	return nil
}

func isAddCmd(text string) bool { // TODO: to review and delete
	return isURL(text)
}

func isURL(text string) bool { // TODO: to review and delete
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
