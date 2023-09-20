package main

import (
	"context"
	"time"

	telegramHandlers "github.com/Malakhov-cmd/GPT-Telegram-Go.git/src/handler"
	configUtils "github.com/Malakhov-cmd/GPT-Telegram-Go.git/src/struct"
	logger "github.com/Malakhov-cmd/GPT-Telegram-Go.git/src/util"

	openai "github.com/sashabaranov/go-openai"
	telebot "github.com/tucnak/telebot"
)

var (
	log = logger.GetLogger()
)

func main() {
	log.Info("Инициализация программы")

	config := configUtils.GetConfig()

	bot, err := telebot.NewBot(telebot.Settings{
		Token: config.API_Keys.Telegram_Keys[0],
		Poller: &telebot.LongPoller{
			Timeout: 2 * time.Minute,
		},
	})
	if err != nil {
		log.Fatal("Не удалось инициализировать телеграм бота")
	}

	bot.Handle("/start", telegramHandlers.GetStartHandler(bot))
	bot.Handle("/help", telegramHandlers.GetHelpHandler(bot))

	// Подключение к API OpenAI с использованием вашего ключа API
	client := openai.NewClient(config.API_Keys.Openai_Keys[0])

	// Обработчик событий при получении сообщения
	bot.Handle(telebot.OnText, func(m *telebot.Message) {
		waitingMessage, err := bot.Send(m.Sender, "генерация ответа...")

		_, err = bot.Raw("sendChatAction", map[string]interface{}{
			"chat_id": m.Chat.ID,
			"action":  "typing",
		})

		if err != nil {
			log.Debug("Не удалось установить typing action")
		}

		// Отправка текста сообщения в GPT-3 для генерации ответа
		response, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: m.Text,
					},
				},
			},
		)

		bot.Delete(waitingMessage)

		if err != nil {
			log.Fatal("Не удалось подключиться к серверам OpenAI")
		}

		// Отправка ответа от GPT-3 обратно пользователю
		bot.Send(m.Sender, response.Choices[0].Message.Content)
	})

	// Запуск бота
	bot.Start()
}
