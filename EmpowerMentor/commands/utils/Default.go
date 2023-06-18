package utils

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// Default is a default handler for all commands that are not implemented
func (App *Application) Default(upd tgbotapi.Update) {
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "I don't know what to do with this command")
	if _, err := App.Bot.Send(msg); err != nil {
		App.Logger.Printf("error sending message: %v\n", err)
	}
	err := App.DB.SaveMessage(upd)
	if err != nil {
		App.Logger.Printf("error saving message into database: %v\n", err)
		return
	}
}
