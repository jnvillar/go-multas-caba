package commands

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"multas_caba/browser"
)

type CommandHandler struct {
}

func New() *CommandHandler {
	return &CommandHandler{}
}

func (c *CommandHandler) TransitFines(msg *tgbotapi.Message) string {
	sumOfTrafficFines := browser.TransitFines(msg.Text)
	if sumOfTrafficFines == "" {
		return "Fallo la consulta"
	}
	return fmt.Sprintf("Debes: %s", sumOfTrafficFines)
}
