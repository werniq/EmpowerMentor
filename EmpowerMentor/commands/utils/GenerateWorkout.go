package utils

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

// GenerateWorkoutForParticularMuscleGroup creates workout for this muscle groups:
// []string{"back", "biceps", "chest", "triceps", "legs", "shoulders"}
func (App *Application) GenerateWorkoutForParticularMuscleGroup(update tgbotapi.Update) {
	args := strings.Split(update.Message.Text, " ")
	if len(args) <= 2 {
		App.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "please, provide muscle group and count of exercises you want to generate"))
		return
	}

	count := args[2]
	m := args[1]
	if StringInArray(m, muscles) == false {
		App.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Please, provide valid muscle group."))
		return
	}

	c, err := strconv.Atoi(count)
	if err != nil {
		App.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "error converting"+count+"to int"))
		return
	}
	exercises, err := App.DB.GenerateXRandomExercises(m, c)
	if err != nil {
		App.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "error generating workouts"+err.Error()))
		return
	}

	if len(exercises) < 2 {
		exercises, err = App.DB.GenerateXRandomExercises(m, c)
		if err != nil {
			App.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "error generating workouts"))
			return
		}
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	for i := 0; i <= len(exercises)-1; i++ {
		msg.Text += fmt.Sprintf("%s \n %s \n %s \n \n\n", exercises[i].Title, exercises[i].Technique, exercises[i].VideoURI)
	}

	_, _ = App.Bot.Send(msg)
}

// CreateWeekTrainingPlan creates a training plan for the whole week.
func (App *Application) CreateWeekTrainingPlan(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID

	workouts := []string{"back", "biceps", "chest", "triceps", "legs", "shoulders"}

	for i := 0; i+1 < len(workouts); i = i + 2 {
		// Generate exercises for the first workout
		exercises1, err := App.DB.GenerateXRandomExercises(workouts[i], 4)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "Something went wrong while generating the workout. Please try again later."+err.Error())
			App.Bot.Send(msg)
			return
		}

		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Exercises for %s and %s", workouts[i], workouts[i+1]))
		App.Bot.Send(msg)
		// Append exercises for the first workout to the message
		for _, exercise := range exercises1 {
			msg.Text = fmt.Sprintf("%s \n %s \n %s \n\n", exercise.Title, exercise.Technique, exercise.VideoURI)
			App.Bot.Send(msg)
		}

		// Generate exercises for the second workout
		exercises2, err := App.DB.GenerateXRandomExercises(workouts[i+1], 3)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "Something went wrong while generating the workout. Please try again later."+err.Error())
			App.Bot.Send(msg)
			return
		}

		// Append exercises for the second workout to the message
		for _, exercise := range exercises2 {
			msg.Text = fmt.Sprintf("%s\n%s\n%s\n\n", exercise.Title, exercise.Technique, exercise.VideoURI)
			App.Bot.Send(msg)
		}
	}
}

func (App *Application) RetrieveExerciseRecommendations(update tgbotapi.Update) {
	callbackData := update.CallbackQuery.Data

	if callbackData == "glutes" || callbackData == "chest" || callbackData == "back" || callbackData == "legs" || callbackData == "shoulders" || callbackData == "triceps" || callbackData == "biceps" {
		switch callbackData {

		case "chest":
			exercise, err := App.DB.GetOneRandomExercise("chest")
			if err != nil {
				App.Logger.Printf("error getting random exercise: %v\n", err)
				return
			}

			text := fmt.Sprintf(`
							One of the best exercises for chest is:
								%s	\n
								Technique for this exercise: %s
								Video: %s
						`, exercise.Title, exercise.Technique, exercise.VideoURI)

			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
			_, _ = App.Bot.Send(msg)

		case "back":
			exercise, err := App.DB.GetOneRandomExercise("back")
			if err != nil {
				App.Logger.Printf("error getting random back exercise: %v\n", err)
				return
			}

			text := fmt.Sprintf(`
						One of the best exercises for back is:
							%s
							Technique for this exercise is: %s
							Video: %s
					`, exercise.Title, exercise.Technique, exercise.VideoURI)

			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
			_, _ = App.Bot.Send(msg)

		case "legs":
			exercise, err := App.DB.GetOneRandomExercise("legs")
			if err != nil {
				App.Logger.Printf("error getting random exercise: %v\n", err)
				return
			}

			text := fmt.Sprintf(`
							One of the best exercises for legs is:
								%s	\n
								Technique for this exercise: %s
								Video: %s
						`, exercise.Title, exercise.Technique, exercise.VideoURI)

			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
			_, _ = App.Bot.Send(msg)

		case "shoulders":
			exercise, err := App.DB.GetOneRandomExercise("shoulders")
			if err != nil {
				App.Logger.Printf("error getting random exercise: %v\n", err)
				return
			}

			text := fmt.Sprintf(`
							One of the best exercises for shoulders is:
								%s	\n
								Technique for this exercise: %s
								Video: %s
						`, exercise.Title, exercise.Technique, exercise.VideoURI)

			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
			_, _ = App.Bot.Send(msg)

		case "triceps":
			exercise, err := App.DB.GetOneRandomExercise("triceps")
			if err != nil {
				App.Logger.Printf("error getting random exercise: %v\n", err)
				return
			}

			text := fmt.Sprintf(`
							One of the best exercises for triceps is:
								%s	\n
								Technique for this exercise: %s
								Video: %s
						`, exercise.Title, exercise.Technique, exercise.VideoURI)

			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
			_, _ = App.Bot.Send(msg)

		case "biceps":
			exercise, err := App.DB.GetOneRandomExercise("biceps")
			if err != nil {
				App.Logger.Printf("error getting random exercise: %v\n", err)
				return
			}

			text := fmt.Sprintf(`
							One of the best exercises for biceps is:
								%s	\n
								Technique for this exercise: %s
								Video: %s
						`, exercise.Title, exercise.Technique, exercise.VideoURI)

			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
			_, _ = App.Bot.Send(msg)

		case "glutes":
			exercise, err := App.DB.GetOneRandomExercise("glutes")
			if err != nil {
				App.Logger.Printf("error getting random exercise: %v\n", err)
				return
			}

			text := fmt.Sprintf(`
							One of the best exercises for glutes is:
								%s	\n
								Technique for this exercise: %s
								Video: %s		
						`, exercise.Title, exercise.Technique, exercise.VideoURI)
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
			_, _ = App.Bot.Send(msg)
		}
	}
	editMsg := tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, tgbotapi.InlineKeyboardMarkup{})
	_, _ = App.Bot.Send(editMsg)
	return
}
