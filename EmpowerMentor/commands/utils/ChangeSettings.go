package utils

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"self-improvement-bot/models"
	"strconv"
	"strings"
	"time"
)

// BotConfiguration is a handler for /configure command
func (App *Application) BotConfiguration(upd tgbotapi.Update) {
	var userBotConfiguration models.UserBotConfiguration
	userBotConfiguration.Step = 0

	message := upd.Message
	chatID := message.Chat.ID

	userBotConfiguration.UserId = message.From.ID

	state, exists := ConfigureStates[chatID]
	if !exists {
		state = userBotConfiguration
		ConfigureStates[chatID] = state
	}

	msg := tgbotapi.NewMessage(chatID, "")

	switch state.Step {
	case 0:
		msg.Text = "How do you want me to call you?"
		_, err := App.Bot.Send(msg)
		if err != nil {
			App.Logger.Printf("error sending message: %v\n", err)
			return
		}
		state.Step = 2
		ConfigureStates[chatID] = state
	case 2:
		state.Username = message.Text
		_, err := App.Bot.Send(msg)
		msg.Text = "Certainly! Next, what is your gender?"
		if err != nil {
			App.Logger.Printf("Error sending message: %v\n", err)
			return
		}
		state.Step++
		ConfigureStates[chatID] = state
	case 3:
		if message.Text == "m" || message.Text == "M" {
			state.Gender = "M"
		}
		if message.Text == "f" || message.Text == "F" {
			state.Gender = "F"
		}
		msg.Text = "Great! How old are you?"
		_, err := App.Bot.Send(msg)
		if err != nil {
			App.Logger.Printf("Error sending message: %v\n", err)
			return
		}
		state.Step++
		ConfigureStates[chatID] = state
	case 4:
		state.Age, _ = strconv.Atoi(message.Text)
		msg.Text = "What is your weight?"
		_, err := App.Bot.Send(msg)
		if err != nil {
			App.Logger.Printf("Error sending message: %v\n", err)
			return
		}
		state.Step++
		ConfigureStates[chatID] = state
	case 5:
		state.Weight, _ = strconv.ParseFloat(message.Text, 32)
		msg.Text = "What is your height?"
		_, err := App.Bot.Send(msg)
		if err != nil {
			App.Logger.Printf("Error sending message: %v\n", err)
			return
		}
		state.Step++
		ConfigureStates[chatID] = state
	case 6:
		state.Height, _ = strconv.ParseFloat(message.Text, 32)
		msg.Text = "What is your preferred physical activity?"
		_, err := App.Bot.Send(msg)
		if err != nil {
			App.Logger.Printf("Error sending message: %v\n", err)
			return
		}
		state.Step++
		ConfigureStates[chatID] = state
	case 7:
		state.PreferredPhysicalActivity = message.Text
		msg.Text = "How many times a week do you workout?"
		_, err := App.Bot.Send(msg)
		if err != nil {
			App.Logger.Printf("Error sending message: %v\n", err)
			return
		}
		state.Step++
		ConfigureStates[chatID] = state
	case 8:
		state.WorkoutCount, _ = strconv.Atoi(message.Text)
		msg.Text = "How many pages do you read a month?"
		_, err := App.Bot.Send(msg)
		if err != nil {
			App.Logger.Printf("Error sending message: %v\n", err)
			return
		}
		state.Step++
		ConfigureStates[chatID] = state
	case 9:
		state.BooksCount, _ = strconv.Atoi(message.Text)
		msg.Text = "What supplements do you prefer? Or what supplements would you like to try? "
		_, err := App.Bot.Send(msg)
		if err != nil {
			App.Logger.Printf("Error sending message: %v\n", err)
			return
		}
		state.Step++
		ConfigureStates[chatID] = state
	case 10:
		state.PreferringSupplements = message.Text
		msg.Text = "What habits would you like to acquire?"
		_, err := App.Bot.Send(msg)
		if err != nil {
			App.Logger.Printf("Error sending message: %v\n", err)
			return
		}
		state.Step++
		ConfigureStates[chatID] = state
	case 11:
		state.HabitsToAcquire = message.Text
		msg.Text = "What news categories are you interested in? If you are not interested in news, just type 'none'."
		_, err := App.Bot.Send(msg)
		if err != nil {
			App.Logger.Printf("Error sending message: %v\n", err)
			return
		}
		state.Step++
		ConfigureStates[chatID] = state
	case 12:
		if message.Text == "none" {
			msg.Text = "What time do you usually wake up?"
			_, err := App.Bot.Send(msg)
			if err != nil {
				App.Logger.Printf("Error sending message: %v\n", err)
				return
			}
			state.Step++
			ConfigureStates[chatID] = state
			return
		}
		state.NewsCategories = message.Text
		msg.Text = "What time do you usually wake up?"
		_, err := App.Bot.Send(msg)
		if err != nil {
			App.Logger.Printf("Error sending message: %v\n", err)
			return
		}
		state.Step++
		ConfigureStates[chatID] = state
	case 13:
		state.WakeUpTime, _ = time.Parse("15:04", message.Text)
		msg.Text = "Last question! What time do you usually go to bed?"
		_, err := App.Bot.Send(msg)

		if err != nil {
			App.Logger.Printf("Error sending message: %v\n", err)
			return
		}
		state.Step++
		ConfigureStates[chatID] = state
	case 14:
		state.BedTime, _ = time.Parse("15:04", message.Text)
		msg.Text = "Thank you for your answers! I will send you a message every day to remind you to drink water, go to bed at time you provided, send you message at your wake up time, do physical exercises, and meditate. I will also send you news articles that you might be interested in. If you want to change your answers, you should re-configure bot from start."
		_, err := App.Bot.Send(msg)

		if err != nil {
			App.Logger.Printf("Error sending message: %v\n", err)
			return
		}
		return
	}
}

// ChangePreferableMeditationTime is a handler for /change_meditation_time command
func (App *Application) ChangePreferableMeditationTime(upd tgbotapi.Update) {
	args := strings.Split(upd.Message.Text, " ")
	inputedTime := args[1]
	t, err := time.Parse("15:15", inputedTime)

	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "")

	if err != nil {
		msg.Text = "Wrong time format. Please, input time in format HH:MM"
		App.Bot.Send(msg)
		return
	}

	err = App.DB.SetDailyMeditationReminderForUser(upd.Message.From.ID, t)
	if err != nil {
		msg.Text = "Something went wrong. Please, try again later."
		App.Bot.Send(msg)
		return
	}

	msg.Text = "Your daily meditation reminder has been changed to " + inputedTime
	App.Bot.Send(msg)
}

// ChangePreferableReadingTime changes preferable reading time for user
func (App *Application) ChangePreferableReadingTime(upd tgbotapi.Update) {
	args := strings.Split(upd.Message.Text, " ")
	inputedTime := args[1]
	t, err := time.Parse("15:15", inputedTime)

	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "")

	if err != nil {
		msg.Text = "Wrong time format. Please, input time in format HH:MM"
		App.Bot.Send(msg)
		return
	}

	err = App.DB.SetDailyReadingReminderForUser(upd.Message.From.ID, t)
	if err != nil {
		msg.Text = "Something went wrong. Please, try again later."
		App.Bot.Send(msg)
		return
	}

	msg.Text = "Your daily reading reminder has been changed to " + inputedTime
	App.Bot.Send(msg)
}

// ChangePreferableExerciseTime changes preferable exercise time for user
func (App *Application) ChangePreferableExerciseTime(upd tgbotapi.Update) {
	args := strings.Split(upd.Message.Text, " ")
	inputedTime := args[1]
	t, err := time.Parse("15:15", inputedTime)

	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "")

	if err != nil {
		msg.Text = "Wrong time format. Please, input time in format HH:MM"
		App.Bot.Send(msg)
		return
	}

	err = App.DB.SetDailyExerciseReminderForUser(upd.Message.From.ID, t)
	if err != nil {
		msg.Text = "Something went wrong. Please, try again later."
		App.Bot.Send(msg)
		return
	}

	msg.Text = "Your daily exercise reminder has been changed to " + inputedTime
	App.Bot.Send(msg)
}

// ChangePreferableSleepingTime changes preferable sleeping time for user
func (App *Application) ChangePreferableSleepingTime(upd tgbotapi.Update) {
	args := strings.Split(upd.Message.Text, " ")
	inputedTime := args[1]
	t, err := time.Parse("15:15", inputedTime)

	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "")

	if err != nil {
		msg.Text = "Wrong time format. Please, input time in format HH:MM"
		App.Bot.Send(msg)
		return
	}

	err = App.DB.SetDailySleepingReminderForUser(upd.Message.From.ID, t)
	if err != nil {
		msg.Text = "Something went wrong. Please, try again later."
		App.Bot.Send(msg)
		return
	}

	msg.Text = "Your daily sleep reminder has been changed to " + inputedTime
	App.Bot.Send(msg)
}

// ChangePreferableWakeUpTime changes preferable wake-up time for user
func (App *Application) ChangePreferableWakeUpTime(update tgbotapi.Update) {
	args := strings.Split(update.Message.Text, " ")
	inputedTime := args[1]
	t, err := time.Parse("15:15", inputedTime)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	if err != nil {
		msg.Text = "Wrong time format. Please, input time in format HH:MM"
		App.Bot.Send(msg)
		return
	}

	err = App.DB.SetDailyWakeUpReminderForUser(update.Message.From.ID, t)
	if err != nil {
		msg.Text = "Something went wrong. Please, try again later."
		App.Bot.Send(msg)
		return
	}

	msg.Text = "Your daily wake up reminder has been changed to " + inputedTime
	App.Bot.Send(msg)
}
