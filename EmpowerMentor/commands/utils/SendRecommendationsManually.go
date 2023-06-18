package utils

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math/rand"
)

// ManuallySendExerciseRecommendations sends exercise recommendations to the user
func (App *Application) ManuallySendExerciseRecommendations(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "What muscles you want to focus on?")

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏋️Chest", "chest"),
			tgbotapi.NewInlineKeyboardButtonData("🏅Back", "back"),
			tgbotapi.NewInlineKeyboardButtonData("🦿Legs", "legs"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👨‍👩‍👦Shoulders", "shoulders"),
			tgbotapi.NewInlineKeyboardButtonData("🥊Triceps", "triceps"),
			tgbotapi.NewInlineKeyboardButtonData("🦾Biceps", "biceps"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🍑Glutes", "glutes"),
		),
	)
	msg.ReplyMarkup = keyboard
	_, err := App.Bot.Send(msg)
	if err != nil {
		App.Logger.Printf("error sending message: %v\n", err)
	}
}

// ManuallySendSleepingRecommendations sends sleeping recommendations to the user
func (App *Application) ManuallySendSleepingRecommendations(update tgbotapi.Update) {
	// generate random number between 1 and len(sleepingRecommendations)

	num := rand.Intn(len(sleepingRecommendations)-1) + 1

	text := fmt.Sprintf("Sleeping recommendation for today: %s", sleepingRecommendations[num])

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)

	_, _ = App.Bot.Send(msg)
}

// ManuallySendRecommendationsForMentalHealth sends recommendations for mental health to the user
func (App *Application) ManuallySendRecommendationsForMentalHealth(update tgbotapi.Update) {
	num := rand.Intn(len(mentalHealthRecommendations)-1) + 1

	text := fmt.Sprintf("Mental health recommendation for today: %s", mentalHealthRecommendations[num])

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)

	_, _ = App.Bot.Send(msg)
}
