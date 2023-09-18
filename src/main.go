package main

import (
	"context"
	"log"
	"time"

	configUtils "github.com/Malakhov-cmd/GPT-Telegram-Go.git/src/struct"
	openai "github.com/sashabaranov/go-openai"
	"github.com/tucnak/telebot"
)

func main() {

	config := configUtils.GetConfig()

	bot, err := telebot.NewBot(telebot.Settings{
		Token: config.API_Keys.Telegram_Keys[0],
		Poller: &telebot.LongPoller{
			Timeout: 2 * time.Minute,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Подключение к API OpenAI с использованием вашего ключа API
	client := openai.NewClient(config.API_Keys.Openai_Keys[0])

	// Обработчик событий при получении сообщения
	bot.Handle(telebot.OnText, func(m *telebot.Message) {
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

		if err != nil {
			log.Fatal(err)
		}

		// Отправка ответа от GPT-3 обратно пользователю
		bot.Send(m.Sender, response.Choices[0].Message.Content)
	})

	// Запуск бота
	bot.Start()
}
