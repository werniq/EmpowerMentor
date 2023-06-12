//package commands
//
//import (
//	"fmt"
//	"github.com/gin-gonic/gin"
//	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
//	"strings"
//)
//
//func (App *Application) ConfigureRoutes(r *gin.Engine) {
//	r.POST("/webhook"+App.Bot.Token, func(c *gin.Context) {
//		fmt.Println("webhook")
//		var update tgbotapi.Update
//		if err := c.BindJSON(&update); err != nil {
//			App.Logger.Printf("error binding update: %v\n", err)
//			return
//		}
//		if update.Message == nil {
//			return
//		}
//
//		go func() {
//			//times, reminders, err := App.DB.GetUserReminders(update.Message.Chat.ID)
//			//if err != nil {
//			//	App.Logger.Printf("error getting user reminders: %v\n", err)
//			//	return
//			//}
//			// preferable_time_to_meditate, preferable_time_to_exercise, preferable_time_to_read
//			if ok, err := App.DB.UserExists(update.Message.From.ID); err != nil {
//				if !ok {
//					return
//				}
//
//				// if user exists, send reminders
//				go func() {
//					App.DailyWaterReminder(update.Message.From.ID)
//				}()
//
//				go func() {
//					App.DailySleepReminder(update.Message.From.ID)
//				}()
//
//				go func() {
//					App.DailyExerciseReminder(update.Message.From.ID)
//				}()
//
//				go func() {
//					App.DailyMeditationReminder(update.Message.From.ID)
//				}()
//
//				go func() {
//					App.DailyReadingReminder(update.Message.From.ID)
//				}()
//
//				go func() {
//					App.DailyMorningMotivationalQuote(update.Message.From.ID)
//				}()
//
//				go func() {
//					App.SendChallenges(update)
//				}()
//
//				go func() {
//					App.SendSleepingRecommendations(update)
//				}()
//
//				go func() {
//					App.SendSupplementsRecommendations(update)
//				}()
//
//				go func() {
//					App.SendRecommendationsForMentalHealth(update)
//				}()
//
//				// concurrency functions:
//				// DailyWaterReminder()
//				// DailySleepReminder()
//				// DailyExerciseReminder()
//				// DailyMeditationReminder()
//				// DailyExerciseReminder()
//
//				// DailyMorningMotivationalQuote()
//				// MorningNewsSender()
//				// SendChallenges()
//				go func() {
//					App.RemindForWhatItIsFor(update)
//				}()
//
//				// SendRecommendationsForPhysicalHealth()
//
//				// if time.Now() is equal to any of the times in times, send a message
//			}
//		}()
//
//		args := strings.Split(update.Message.Text, " ")
//		command := args[0]
//		fmt.Println(command)
//		fmt.Println(update.Message.Chat.ID)
//
//		if update.Message.ReplyToMessage != nil {
//			m, err := App.DB.RetrieveLastUserMessage(update.Message.Chat.ID)
//			if err != nil {
//				App.Logger.Printf("error retrieving last user message: %v\n", err)
//				App.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error retrieving last user message. Please try again later, and contact the admin if the problem persists."))
//				return
//			}
//			args := strings.Split(m.Text, " ")
//			command = args[0]
//
//			switch command {
//			case "/add-quote":
//			}
//		}
//
//		switch command {
//
//		case "/setup":
//			App.BotConfiguration(update)
//
//		case "/start":
//			App.Start(update)
//
//		case "/create-supplement-plan":
//			App.CreateCustomSupplementIntakePlan(update)
//
//		case "/exercise-recommendations":
//			App.ManuallySendExerciseRecommendations(update)
//
//		case "/mental-health-recommendations":
//			App.ManuallySendRecommendationsForMentalHealth(update)
//
//		case "/sleeping-recommendations":
//			App.ManuallySendSleepingRecommendations(update)
//
//		case "/report":
//			App.ReportToAdmins(update)
//
//		case "/add-quote":
//			App.UploadMotivationalQuote(update)
//
//		case "/add-challenge-for-review":
//			App.RequestToUploadChallenge(update)
//
//		case "/generate-custom-workout":
//			App.GenerateWorkoutForParticularMuscleGroup(update)
//
//		case "/generate-week-training-plan":
//			App.CreateWeekTrainingPlan(update)
//
//		case "/change-meditation-time":
//			App.ChangePreferableMeditationTime(update)
//
//		case "/change-reading-time":
//			App.ChangePreferableReadingTime(update)
//
//		case "/change-exercise-time":
//			App.ChangePreferableExerciseTime(update)
//
//		case "/change-sleeping-time":
//			App.ChangePreferableSleepingTime(update)
//
//		case "/change-wake-up-time":
//			App.ChangePreferableWakeUpTime(update)
//
//		case "/add-admin":
//			App.AddAdmin(update)
//
//		case "/create-custom-goal":
//			App.CreateCustomGoal(update)
//
//		case "/create-custom-challenge":
//			App.CreateCustomChallenge(update)
//
//		case "/create-custom-reminder":
//			App.CreateCustomReminder(update)
//
//		case "/create-custom-habit":
//			App.CreateCustomHabit(update)
//
//		case "/diets":
//			App.GetDifferentDietInfo(update)
//
//			// todo: finish
//		case "/create-custom-meal-preparing-plan":
//			App.CreateCustomMealPreparingPlan(update)
//
//		case "/buy-subscription":
//			// TODO: finish
//			App.BuySubscription(update)
//
//		case "/add-subscription-to-morning-newsletter":
//			App.AddSubscriptionOnMorningNews(update)
//
//		case "/manage-users":
//			App.ManageUsers(update)
//
//		case "/manage-subscriptions":
//			App.ManageSubscriptions(update)
//
//		case "/remove-subs":
//			App.RemoveSubscription(update)
//
//		case "/create-custom-motivation-reminder":
//			App.CreateCustomMotivationalReminder(update)
//
//		case "/receive-crypto-donation":
//			// Todo: finish
//			App.ReceiveCryptoDonation(update)
//
//		default:
//			App.Default(update)
//		}
//	})
//}
