package utils

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type RegisterState struct {
	Step int
	Name string
}

func (App *Application) Register(update tgbotapi.Update) {

	// TODO: check if user is already registered

	ok, err := App.DB.IsUserRegistered(update.Message.From.ID)
	if err != nil {
		App.Logger.Printf("Error while checking if user is registered: %v", err)
		return
	}

	if ok {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You are already registered")
		App.Bot.Send(msg)
		return
	}

}
