package commands

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go-multas-caba/browser"
	"go-multas-caba/validator"
)

type CommandHandler struct {
}

func New() *CommandHandler {
	return &CommandHandler{}
}

func (c *CommandHandler) TransitFines(msg *tgbotapi.Message) string {
	params := strings.Fields(msg.Text)
	if len(params) < 2 {
		return "Falta el dominio"
	}
	domain := params[1]
	errors := []error{validator.MaxLength(domain, 8), validator.MinLength(domain, 6)}
	err := validator.CheckErrors(errors)
	if err != nil {
		return err.Error()
	}
	sumOfTrafficFines := browser.TransitFines(domain)
	if sumOfTrafficFines == "" {
		return "Fallo la consulta"
	}
	return fmt.Sprintf("Debes: %s", sumOfTrafficFines)
}
