package telegram

const msgHelp = `@sber_invest_bot - Сервис для работы с платформой <a href="https://sber-invest.kz"><b>«Сбережения и Инвестиции»</b></a>.

<b>Сервис может:</b>
- отсортировать публикаций по рубрикам; 
- отсортировать публикаций по авторам. 
- найти заголовки всех публикаций в базе знаний;   
- найти заголовки всех статей опубликованных за последние 7 и 30 дней. 

<b>Главное меню:</b>
/articles - База знаний <a href="https://sber-invest.kz/knowledgebase"><b>«Сбережения и Инвестиции»</b></a>.
/help - показывает это сообщение.
`

const msgHello = "<b>Привет!</b>\n\n" + msgHelp

const (
	msgUnknownCommand = "Незнакомая команда 🤔"
	msgNoSavedPages   = "You have no saved pages 🙊"
	msgSaved          = "Saved! 👌"
	msgAlreadyExists  = "You have already have this page in your list 🤗"
	msgStatusChanged  = "Статус изменен ✅"
)
