package utils

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math/rand"
	"time"
)

// SendChallenges sends challenges to users
func (App *Application) SendChallenges(update tgbotapi.Update) {
	if time.Now().Hour() == 11 && time.Now().Minute() == 0 && time.Now().Second() == 0 {
		challenge, err := App.DB.RetrieveUserRandomCustomChallenge(update.Message.From.ID)
		if err != nil {
			App.Bot.Send(App.TranslateToUkrainian(tgbotapi.NewMessage(update.Message.Chat.ID, "Error getting challenge. Please, report error using command /report <issue> ")))
			return
		}

		App.Bot.Send(App.TranslateToUkrainian(tgbotapi.NewMessage(update.Message.Chat.ID, challenge)))
	}
}

// SendSleepingRecommendations sends recommendations to users for sleeping well and healthy
func (App *Application) SendSleepingRecommendations(update tgbotapi.Update) {
	if time.Now().Hour() == 20 && time.Now().Minute() == 0 && time.Now().Second() == 0 {
		num := rand.Intn(len(sleepingRecommendations))
		App.Bot.Send(App.TranslateToUkrainian(tgbotapi.NewMessage(update.Message.Chat.ID, sleepingRecommendations[num])))
	}
}

// SendRecommendationsForMentalHealth sends recommendations to users for mental health
func (App *Application) SendRecommendationsForMentalHealth(update tgbotapi.Update) {
	if time.Now().Hour() == 20 && time.Now().Minute() == 0 && time.Now().Second() == 0 {
		num := rand.Intn(len(mentalHealthRecommendations))
		App.Bot.Send(App.TranslateToUkrainian(tgbotapi.NewMessage(update.Message.Chat.ID, mentalHealthRecommendations[num])))
	}
}

// RemindForWhatItIsFor functions periodically sends message to user for motivating him/her
func (App *Application) RemindForWhatItIsFor(update tgbotapi.Update) {
	if time.Now().Hour() == 12 && time.Now().Minute() == 0 && time.Now().Second() == 0 {
		quote, err := App.DB.GetUserCustomMotivation(update.Message.From.ID)
		if err != nil {
			App.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Something went wrong: "+err.Error()+"\n Please, report to admins by using command /report {<-issue->}"))
			return
		}
		App.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Do not forget about: "+quote))
	}
}

// DailyWaterReminder sends a message to the user with a reminder to drink water
func (App *Application) DailyWaterReminder(channelId int64) {
	if time.Now().Hour()%4 == 0 && time.Now().Hour() < 22 && time.Now().Hour() > 7 && time.Now().Minute() == 0 && time.Now().Second() == 0 {
		msg := tgbotapi.NewMessage(channelId, "Don't forget to drink water today!")
		if _, err := App.Bot.Send(msg); err != nil {
			App.Logger.Printf("error sending message: %v\n", err)
			return
		}
	}
}

// DailySleepReminder sends a message to the user with a reminder to sleep
func (App *Application) DailySleepReminder(channelId int64) {
	// TODO: add sleep recommendations
	// TODO: add user preferable time for sleep, and send message exactly at that time
	if time.Now().Hour() == 22 && time.Now().Minute() == 30 && time.Now().Second() == 0 {
		msg := tgbotapi.NewMessage(channelId, "Don't forget to sleep well today!")
		if _, err := App.Bot.Send(msg); err != nil {
			App.Logger.Printf("error sending message: %v\n", err)
			return
		}
	}
}

// DailyMeditationReminder sends a message to the user with a reminder to meditate
func (App *Application) DailyMeditationReminder(channelId int64) {
	// TODO: add meditation recommendations
	// TODO: add user preferable time for meditation, and send message exactly at that time
	if time.Now().Hour() == 7 && time.Now().Minute() == 30 && time.Now().Second() == 0 && time.Now().Second() == 0 {
		msg := tgbotapi.NewMessage(channelId, "Don't forget to meditate today!")
		if _, err := App.Bot.Send(msg); err != nil {
			App.Logger.Printf("error sending message: %v\n", err)
			return
		}
	}
}

// DailyReadingReminder sends a message to the user with a reminder to read
func (App *Application) DailyReadingReminder(channelId int64) {
	if time.Now().Hour() == 15 && time.Now().Minute() == 30 && time.Now().Second() == 0 {
		msg := tgbotapi.NewMessage(channelId, "Don't forget to read today!")
		if _, err := App.Bot.Send(msg); err != nil {
			App.Logger.Printf("error sending message: %v\n", err)
			return
		}
	}
}

// DailyExerciseReminder sends a message to the user with a reminder to exercise
func (App *Application) DailyExerciseReminder(channelId int64) {
	// TODO: add exercise plans and recommendations
	// TODO: add user preferable time for training, create custom plan, and send message exactly at that time
	if time.Now().Hour() == 12 && time.Now().Minute() == 30 && time.Now().Second() == 0 {
		msg := tgbotapi.NewMessage(channelId, "Don't forget to exercise today!")
		if _, err := App.Bot.Send(msg); err != nil {
			App.Logger.Printf("error sending message: %v\n", err)
			return
		}
	}
}

// DailyMorningMotivationalQuote sends a message to the user with a motivational quote
func (App *Application) DailyMorningMotivationalQuote(channelId, userId int64) {
	var ok bool
	var err error
	if ok, err = App.DB.UserExists(userId); err != nil {
		App.Logger.Printf("error checking if user exists: %v\n", err)
		return
	}
	if !ok {
		App.Bot.Send(tgbotapi.NewMessage(channelId, "You are not registered. Please, use /setup command to register."))
		return
	}

	wakeUpTime, err := App.DB.GetUserWakeupTime(channelId)
	if err != nil {
		App.Bot.Send(tgbotapi.NewMessage(channelId, "Something went wrong. Please, try again later."))
		return
	}

	if time.Now().Hour() == wakeUpTime.Hour() && time.Now().Minute() == wakeUpTime.Minute() {
		quote, err := App.DB.GetRandomQuote()
		if err != nil {
			App.Logger.Printf("error getting random quote: %v\n", err)
			msg := tgbotapi.NewMessage(channelId, "Error getting random quote")
			_, _ = App.Bot.Send(msg)
		}

		msg := tgbotapi.NewMessage(channelId, quote)
		_, _ = App.Bot.Send(msg)
	}
}

// WalkingReminder sends message to user if it's time to go for a walk
func (App *Application) WalkingReminder(upd tgbotapi.Update) {
	if time.Now().Hour() == 12 && time.Now().Minute() == 30 && time.Now().Second() == 0 {
		App.Bot.Send(tgbotapi.NewMessage(upd.Message.Chat.ID, "It's time to go for a walk!"))
	}
}
