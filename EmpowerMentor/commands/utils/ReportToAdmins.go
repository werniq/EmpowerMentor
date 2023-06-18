package utils

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

// ReportToAdmins function to report issues with bot
func (App *Application) ReportToAdmins(update tgbotapi.Update) {
	problem := strings.TrimPrefix("/report", update.Message.Text)

	adminId, err := App.DB.GetRandomAdmin()
	if err != nil {
		App.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Something went wrong.. Sending report to admins"))
		return
	}

	// TODO: create table for storing user issues
	App.Bot.Send(tgbotapi.NewMessage(adminId, fmt.Sprintf(
		`User %d with username [%s] has reported following issue: %s`,
		update.Message.From.ID,
		update.Message.From.UserName,
		problem)))
}
