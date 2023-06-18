package utils

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

// HandleChallengeRequest function handles challenge request
func (App *Application) HandleChallengeRequest(update tgbotapi.Update) {
	if update.CallbackQuery.Data == "v" {
		// now we need to get challenge from database
		userId := update.CallbackQuery.From.ID
		challenge, err := App.DB.GetChallengeUploadRequest(userId)
		if err != nil {
			App.Logger.Printf("error retrieving challange from request_challenge table: %s", err.Error())
			return
		}

		App.ApproveChallengeRequest(userId, challenge)
	} else if update.CallbackQuery.Data == "d" {
		userId := update.CallbackQuery.From.ID
		challenge, err := App.DB.GetChallengeUploadRequest(userId)
		if err != nil {
			App.Logger.Printf("error retrieving challange from request_challenge table: %s", err.Error())
			return
		}

		App.DenyChallengeRequest(userId, challenge)
	}
}

// ApproveChallengeRequest function approves challenge request
func (App *Application) ApproveChallengeRequest(userId int64, challenge string) {
	err := App.DB.UploadChallenge(challenge)
	if err != nil {
		App.Logger.Printf("error uploading challenge to database: %s", err.Error())
		return
	}

	msg := tgbotapi.NewMessage(userId, "Your challenge"+challenge+"has been approved. Thank you for your contribution to this bot. ")
	App.Bot.Send(msg)
}

// DenyChallengeRequest function denies challenge request
func (App *Application) DenyChallengeRequest(userId int64, challenge string) {
	// assume that update.CallbackQuery.Data == "d" -> deny
	err := App.DB.DeleteChallengeRequest(userId)
	if err != nil {
		App.Logger.Printf("error deleting challenge request: %s", err.Error())
		return
	}

	msg := tgbotapi.NewMessage(userId, "Your challenge"+challenge+"has been denied. Please, find new one and try one more time :) ")
	App.Bot.Send(msg)
}

// IsAdmin function checks if user is administrator
func (App *Application) IsAdmin(userId int64) bool {
	if App.DB.IsAdmin(userId) != nil {
		return false
	}
	return true
}

// AddAdmin function adds administrator to this bot
func (App *Application) AddAdmin(update tgbotapi.Update) {
	if App.DB.IsAdmin(update.Message.From.ID) != nil {
		return
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	err := App.DB.AddAdmin(update.Message.From.ID, update.Message.From.UserName)
	if err != nil {
		msg.Text = "error occured while inserting new admin to database"
		_, _ = App.Bot.Send(msg)
		return
	}
	msg.Text = "Successfully added new admin."
	App.Bot.Send(msg)
}

// UploadChallenge function inserts new challenge record into database
func (App *Application) UploadChallenge(update tgbotapi.Update) {
	if !App.IsAdmin(update.Message.From.ID) {
		return
	}

	if ok, err := App.DB.ChallengeExists(strings.TrimPrefix(update.Message.Text, "/upload-challenge")); !ok {
		App.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Something went wrong. Please, report to @werniq issue using /report <-{issue}->: "+err.Error()))
		return
	}

	App.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Successfully  uploaded challenge!"))
}

// UploadMotivationalQuote is function for updating database which contains motivational quotes with new portions
// of motivation
func (App *Application) UploadMotivationalQuote(update tgbotapi.Update) {
	if !App.IsAdmin(update.Message.From.ID) {
		App.Default(update)
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	text := strings.TrimPrefix("/add-quote", update.Message.Text)

	err := App.DB.UploadOneMotivationalQuote(text)
	if err != nil {
		msg.Text = "Error uploading quote to database: " + err.Error()
		return
	}

	msg.Text = "Successfully uploaded quote to database! soon it will motivate people <33"
	_, _ = App.Bot.Send(msg)
}
