package main

import (
	"log"
	"os"

	"multas_caba/commands"
	"multas_caba/validator"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	pb := NewBot()
	var token = os.Getenv("TELEGRAMFINESTOKEN")
	bot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		logMessage(update.Message)

		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			msg := pb.handleCommand(update.Message)
			bot.Send(msg)
		}
	}
}

func logMessage(message *tgbotapi.Message) {
	log.Printf("[%s] %s", message.From.UserName, message.Text)
}

type Bot struct {
	commandHandlers map[string]func(message *tgbotapi.Message) string
}

func NewBot() Bot {
	return Bot{
		commandHandlers: loadHandlers(),
	}
}

func (b *Bot) handleCommand(message *tgbotapi.Message) *tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(message.Chat.ID, "")

	ok, err := validator.Length(msg.Text)
	if err != nil {
		msg.Text = err.Error()
		return &msg
	}

	command := message.Command()
	handler, ok := b.commandHandlers[command]

	if !ok {
		msg.Text = "I don't know that command"
	} else {
		msg.Text = handler(message)
	}
	return &msg
}

func loadHandlers() map[string]func(message *tgbotapi.Message) string {
	ch := commands.New()
	return map[string]func(message *tgbotapi.Message) string{
		"multas": ch.TransitFines,
	}
}
