package utils

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
	"time"
)

// CreateCustomGoal creates custom goal for user
func (App *Application) CreateCustomGoal(update tgbotapi.Update) {
	msg := update.Message
	goal := strings.TrimPrefix(msg.Text, "/create-goal")
	if err := App.DB.CreateCustomGoal(msg.From.ID, goal); err != nil {
		App.Bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "error saving data into database. try again later"))
		return
	}
	App.Bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Successfully created custom goal! We would help you achieve that <3"))
}

// CreateCustomReminder initializes new reminder for user, which will appear on provided by user time
func (App *Application) CreateCustomReminder(update tgbotapi.Update) {
	m := update.Message

	args := strings.Split(strings.TrimPrefix("/create-reminder", m.Text), " ")
	reminder := args[0]
	ti := args[1]

	t, err := time.Parse("15:15:15", ti)
	if err != nil {
		App.Bot.Send(tgbotapi.NewMessage(m.Chat.ID, "error parsing time: "+err.Error()))
		return
	}

	err = App.DB.CreateCustomReminder(m.From.ID, reminder, t)
	if err != nil {
		App.Bot.Send(tgbotapi.NewMessage(m.Chat.ID, "Something went wrong. Please, report to admins by typing /report"))
		return
	}
}

// CreateCustomMotivationalReminder saves user answer to the question 'What is your motivation?'
// it will be used to once a day remind user about his/her motivation
func (App *Application) CreateCustomMotivationalReminder(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, `
			Each of us has a someone, or something that push us to move forward. Not to give up. To keep going, fighting.
			What is yours? This message will be displayed as 'Do not forget about <X>' or 'Do it because <X>'. Please, send reply to this message with your answer. `)

	msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true, Selective: true}

	_, _ = App.Bot.Send(msg)
}

// CreateCustomMotivation function saves user answer to the question 'What is your motivation?'
func (App *Application) CreateCustomMotivation(update tgbotapi.Update) {
	// if command := /create-custom-motivation  -> create new custom motivation
	m := update.Message
	reason := m.Text

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	err := App.DB.SaveUserCustomMotivation(update.Message.From.ID, reason)
	if err != nil {
		msg.Text = "Error saving your motivation. Please, try again later."
		App.Bot.Send(msg)
	}

	msg.Text = "Your custom motivation has been saved. Thank you!"
	App.Bot.Send(msg)
}

// CreateCustomChallenge creates request
func (App *Application) CreateCustomChallenge(update tgbotapi.Update) {
	// if command := /create-custom-challenge  -> create new custom challenge
	challenge := strings.TrimPrefix(update.Message.Text, "/create-custom-challenge")

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	err := App.DB.SaveUserCustomChallenge(update.Message.From.ID, challenge)
	if err != nil {
		msg.Text = "Error saving your challenge. Please, try again later."
		App.Bot.Send(msg)
	}

	msg.Text = "Your custom challenge has been saved. Thank you!"
	_, err = App.Bot.Send(msg)
	if err != nil {
		App.Logger.Printf("error sending message: %s", err.Error())
		return
	}
}
