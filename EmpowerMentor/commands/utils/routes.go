package utils

import (
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"regexp"
	"strings"
)

func (App *Application) ConfigureRoutes(r *gin.Engine) {
	r.POST("/webhook"+App.Bot.Token, func(c *gin.Context) {
		var update tgbotapi.Update
		if err := c.BindJSON(&update); err != nil {
			App.Logger.Printf("error binding update: %v\n", err)
			return
		}

		if update.Message != nil {
			if update.Message.From.IsBot {
				return
			}
		}

		// Create a regular expression pattern for a four-digit number with the first digit as 1, 2, or 3
		pattern := "^[1-3][0-9]{3}$"

		// somehow I should track user responses to particular messages
		// I need to create database table "messages", which will store user messages
		if update.CallbackQuery != nil {
			// retrieve last message and further actions depends on that message
			if update.CallbackQuery.Data == "glutes" || update.CallbackQuery.Data == "chest" || update.CallbackQuery.Data == "back" || update.CallbackQuery.Data == "legs" || update.CallbackQuery.Data == "shoulders" || update.CallbackQuery.Data == "arms" {
				App.RetrieveExerciseRecommendations(update)
			} else if update.CallbackQuery.Data == "v" || update.CallbackQuery.Data == "d" {
				// approve or deny request
				App.HandleChallengeRequest(update)
			} else if update.CallbackQuery.Data == "health" || update.CallbackQuery.Data == "fitness" || update.CallbackQuery.Data == "personal" || update.CallbackQuery.Data == "psychology" || update.CallbackQuery.Data == "mindfulness" || update.CallbackQuery.Data == "self-Care" || update.CallbackQuery.Data == "motivation" || update.CallbackQuery.Data == "productivity" || update.CallbackQuery.Data == "happiness" || update.CallbackQuery.Data == "relationships" || update.CallbackQuery.Data == "career-development" || update.CallbackQuery.Data == "leadership" || update.CallbackQuery.Data == "education" || update.CallbackQuery.Data == "self-Help" || update.CallbackQuery.Data == "mental" || update.CallbackQuery.Data == "well-being" {
				App.AddSubscriptionOnMorningNewsCallback(update)
			} else if regexp.MustCompile(pattern).MatchString(update.CallbackQuery.Data) ||
				update.CallbackQuery.Data == "day" ||
				update.CallbackQuery.Data == "week" ||
				update.CallbackQuery.Data == "month" ||
				StringInArray(update.CallbackQuery.Data, DietNames) ||
				StringInArray(update.CallbackQuery.Data, Allergens) {

				App.CallbackMealPreparePlan(update)
			}
			return
		}

		args := strings.Split(update.Message.Text, " ")
		command := args[0]

		if command == "/setup" {
			App.BotConfiguration(update)
			return
		}

		go func() {
			// TODO: verify if user is subscribed to daily reminders
			App.DailyExerciseReminder(update.Message.Chat.ID)
			App.DailyMeditationReminder(update.Message.Chat.ID)
			App.DailyReadingReminder(update.Message.Chat.ID)
			App.DailyWaterReminder(update.Message.Chat.ID)
			App.DailyMeditationReminder(update.Message.Chat.ID)
			App.DailySleepReminder(update.Message.Chat.ID)
			App.SendChallenges(update)
			App.RemindForWhatItIsFor(update)
			App.WalkingReminder(update)
		}()

		switch command {
		case "/start":
			App.Start(update)

		case "/changemeditationt":
			App.ChangePreferableMeditationTime(update)

		case "/changereadingt":
			App.ChangePreferableReadingTime(update)

		case "/changeexerciset":
			App.ChangePreferableExerciseTime(update)

		case "/changewalkingt":
			App.WalkingReminder(update)

		case "/changesleepingt":
			App.ChangePreferableSleepingTime(update)

		case "/changewakingupt":
			App.ChangePreferableWakeUpTime(update)

		case "/supplement-categories":
			App.ListAvailableSupplementCategories(update)

		case "/supplements-recommendations":
			App.SendSupplementsRecommendations(update)

		case "/recommendationsforexercise":
			App.ManuallySendExerciseRecommendations(update)

		case "/recommendationsformentalhealth":
			App.ManuallySendRecommendationsForMentalHealth(update)

		case "/recommendationsforsleeping":
			App.SendSleepingRecommendations(update)

		case "/customgoal":
			App.CreateCustomGoal(update)

		case "/gen":
			App.GenerateWorkoutForParticularMuscleGroup(update)

		case "/supplementsplan":
			App.CreateCustomSupplementIntakePlan(update)

		case "/weektrainingplan":
			App.CreateWeekTrainingPlan(update)

		case "/motivateme":
			App.GetRandomMotivationalQuote(update)

		case "/upload-quote":
			App.UploadMotivationalQuote(update)

		case "/upload-challenge":
			App.UploadChallenge(update)

		case "/requesttouploadquote":
			App.RequestToUploadChallenge(update)

		case "/addadmin":
			App.AddAdmin(update)

		case "/report":
			App.ReportToAdmins(update)

		case "/customchallenge":
			App.CreateCustomChallenge(update)

		case "/aboutme":
			App.AboutAuthor(update)

		case "/addsubsonnews":
			App.AddSubscriptionOnMorningNews(update)

		case "/mealprepare":
			App.CreateCustomMealPreparingPlan(update)

		case "/custommotivation":
			App.CreateCustomMotivation(update)

		case "/diets":
			App.GetDifferentDietInfo(update)

			// TODO: gen challenge

		default:
			App.Default(update)
		}
	})
}
