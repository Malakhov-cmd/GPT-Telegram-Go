package telegramHandlers

import (
	logger "github.com/Malakhov-cmd/GPT-Telegram-Go.git/src/util"
	telebot "github.com/tucnak/telebot"
)

var (
	log = logger.GetLogger()
)

func GetStartHandler(bot *telebot.Bot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		bot.Send(m.Sender, "Привет, я бот подключенный к OpenAI, просто напиши любой интересующий тебя вопрос и я отвечу на него.")
	}
}

func GetHelpHandler(bot *telebot.Bot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		bot.Send(m.Sender, "Пока не придумал")
	}
}
