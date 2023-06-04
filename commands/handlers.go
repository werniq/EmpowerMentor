package commands

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"math/rand"
	"self-improvement-bot/commands/utils"
	"self-improvement-bot/models"
	"strconv"
	"strings"
	"time"
)

// i am creating bot for tracking self improvement
// want to add subscription for this bot so that users can track their progress
// and with time they will be able to see their progress
// functionality of bot:
// motivational quote few times a day
// daily reminder to track progress
// daily trackers for:
// - water
// - food
// - sleep
// - meditation
// - exercise
// - reading
// additionally, there will be a possibility to add custom trackers
// and to see progress in graphs

// add functionality to send recommendations about book reading
// add functionality to send recommendations about meditation
// add functionality to send recommendations about exercises, for particular muscle groups
// add functionality to send recommendations about sleep

// ability to set up reminders for all of the above
// ability to set up custom goals and reminders for them

// each morning bot will send a message with a quote and a reminder to track progress
// + challenges each morning
// habit building:
// few times a day send a message with a reminder to build a habit

// 1. Daily Affirmations: Start your day with uplifting affirmations and positive quotes to boost your confidence and cultivate a positive mindset.
// 2. Goal Tracker: Set personal goals, track your progress, and receive motivational messages to stay focused and driven.
// 3. Habit Builder: Build positive habits with personalized reminders, tips, and habit streaks analysis for long-term success.
// 4. Meditation and Mindfulness: Experience tranquility and improve your mental well-being with guided meditation sessions and mindfulness exercises.
// 5. Productivity Booster: Unlock your productivity potential with effective time management techniques, productivity tips, and personalized recommendations.
// 6. Self-Care Planner: Prioritize self-care with tailored routines, self-reflection exercises, and curated self-care suggestions for enhanced well-being.
// 7. Personal Development Courses: Learn and grow with self-improvement courses on various topics, delivered by experts in personal development.
// 8. Motivational Challenges: Challenge yourself, step out of your comfort zone, and unlock new achievements to reach your full potential.
// 9. Daily Quotes: Get inspired by daily quotes from the world’s greatest minds and thought leaders.

type Config struct {
	ServerURI string
}

type Application struct {
	Bot            *tgbotapi.BotAPI
	DB             *models.DatabaseModel
	Logger         *log.Logger
	Config         Config
	UpdatesChannel tgbotapi.UpdatesChannel
	Timers         map[int64]map[string]time.Time
}

var (
	sleepingRecommendations = []string{
		"Stick to a consistent sleep schedule: Try to go to bed and wake up at the same time every day, even on weekends. This helps regulate your body's internal clock.",
		"Create a bedtime routine: Establish a relaxing routine before bed. This could include activities such as reading a book, taking a warm bath, or practicing relaxation exercises like deep breathing or meditation.",
		"Create a sleep-friendly environment: Make your bedroom comfortable, cool, dark, and quiet. Consider using blackout curtains, earplugs, or a white noise machine to block out any disruptions.",
		"Limit exposure to electronic devices: Avoid using electronic devices, such as smartphones, tablets, or computers, right before bed. The blue light emitted by these devices can interfere with your sleep.",
		"Avoid caffeine and stimulants: Limit your intake of caffeine, nicotine, and other stimulants, especially in the afternoon and evening. These substances can interfere with your ability to fall asleep.",
		"Regular exercise: Engage in regular physical activity during the day. However, avoid vigorous exercise close to bedtime, as it can make it harder to fall asleep. Aim to finish exercising at least a few hours before bed.",
		"Watch your diet: Avoid heavy meals and spicy or acidic foods close to bedtime, as they can cause discomfort and disrupt sleep. Instead, opt for a light snack if needed.",
		"Manage stress: Practice stress management techniques, such as relaxation exercises, journaling, or talking to a friend or therapist. High stress levels can interfere with your ability to sleep well.",
		"Limit naps: If you have trouble sleeping at night, try to limit daytime napping or keep it short (around 20-30 minutes) and avoid napping too close to bedtime.",
		"Consider your sleep environment: Choose a comfortable mattress, pillows, and bedding that suit your preferences and support your body.",
	}

	mentalHealthRecommendations = []string{
		"Practice self-care: Take time for yourself and engage in activities that you enjoy. This could include hobbies, relaxation techniques, spending time in nature, or engaging in creative outlets.",
		"Prioritize sleep: Ensure you get enough quality sleep as it plays a crucial role in maintaining good mental health.",
		"Stay physically active: Regular exercise can have a positive impact on mental health by reducing stress, improving mood, and increasing overall well-being. Find physical activities that you enjoy and make them a part of your routine.",
		"Maintain a balanced diet: Eat a nutritious diet that includes fruits, vegetables, whole grains, lean proteins, and healthy fats. Avoid excessive consumption of processed foods, sugary snacks, and drinks, as they can impact your mood and energy levels.",
		"Practice mindfulness and relaxation techniques: Mindfulness meditation, deep breathing exercises, and progressive muscle relaxation can help reduce stress and promote a sense of calm and well-being.",
		"Connect with others: Maintain healthy relationships and seek social support from friends, family, or support groups. Connecting with others can provide a sense of belonging, support, and perspective.",
		"Set realistic goals: Break larger goals into smaller, achievable steps. Celebrate your accomplishments along the way and be kind to yourself if you experience setbacks.",
		"Limit exposure to negative news and social media: Constant exposure to negative news and social media can impact your mental well-being. Take breaks from media and choose to engage in positive and uplifting content instead.",
		"Seek professional help if needed: If you're struggling with your mental health, don't hesitate to reach out to a mental health professional. They can provide guidance, support, and appropriate treatments tailored to your needs.",
		"Practice gratitude: Take time each day to reflect on the things you are grateful for. Keeping a gratitude journal or sharing your gratitude with others can help foster a positive mindset."}
	AvaliableSupplementCategories = []string{
		"vitamins",
		"supplements",
		"minerals",
		"digestive-support",
		"antioxidants",
		"bone-joint-cartilage",
		"sleep",
		"fish-oil-omegas-epa-dha",
		"brain-cognitive",
		"hair-skin-nails",
		"greens-superfoods",
		"amino-acids",
		"bee-products",
		"childrens-health",
		"mens-health",
		"mushrooms",
		"weight-loss",
		"phospholipids",
		"protein",
		"womens-health",
	}
	AvailableMuscleGroup = []string{"back", "biceps", "chest", "triceps", "legs", "shoulders"}
)

func (App *Application) Start(upd tgbotapi.Update) {
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "MindBoost Pro is your personal companion for self-improvement and personal growth. Unlock your potential, enhance your well-being, and achieve your goals with our comprehensive range of tools, resources, and expert guidance. Whether you're seeking motivation, mindfulness, or mastery, MindBoost Pro is here to empower you on your journey to becoming the best version of yourself.")
	if _, err := App.Bot.Send(msg); err != nil {
		App.Logger.Printf("error creating message: %v\n", err)
		return
	}
}

func (App *Application) Default(upd tgbotapi.Update) {
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "I don't know what to do with this command")
	if _, err := App.Bot.Send(msg); err != nil {
		App.Logger.Printf("error creating message: %v\n", err)
		return
	}
}

func (App *Application) BotConfiguration(upd tgbotapi.Update) {
	var userBotConfiguration models.UserBotConfiguration
	userBotConfiguration.Step = 1

	states := make(map[int64]models.UserBotConfiguration)

	userBotConfiguration.UserId = upd.Message.From.ID

	for update := range App.UpdatesChannel {
		if update.Message == nil {
			continue
		}

		message := update.Message
		chatId := message.Chat.ID

		state, exists := states[chatId]
		if !exists {
			state = models.UserBotConfiguration{}
			states[chatId] = state
		}

		msg := tgbotapi.NewMessage(chatId, "")

		switch state.Step {
		case 1:
			msg.Text = "How do you want me to call you?"
			App.Bot.Send(msg)
			state.Step++
		case 2:
			state.Username = message.Text
			msg.Text = "Certainly! Next, what is your gender?"
			App.Bot.Send(msg)
			state.Step++
		case 3:
			if message.Text == "m" || message.Text == "M" {
				state.Gender = "M"
			}
			if message.Text == "f" || message.Text == "F" {
				state.Gender = "F"
			}
			msg.Text = "Great! How old are you?"
			App.Bot.Send(msg)
			state.Step++
		case 4:
			state.Age, _ = strconv.Atoi(message.Text)
			msg.Text = "What is your weight?"
			App.Bot.Send(msg)
			state.Step++
		case 5:
			state.Weight, _ = strconv.ParseFloat(message.Text, 32)
			msg.Text = "What is your height?"
			App.Bot.Send(msg)
			state.Step++
		case 6:
			state.Height, _ = strconv.ParseFloat(message.Text, 32)
			msg.Text = "What is your preferred physical activity?"
			App.Bot.Send(msg)
			state.Step++
		case 7:
			state.PreferredPhysicalActivity = message.Text
			msg.Text = "How many times a week do you workout?"
			App.Bot.Send(msg)
			state.Step++
		case 8:
			state.WorkoutCount, _ = strconv.Atoi(message.Text)
			msg.Text = "How many pages do you read a month?"
			App.Bot.Send(msg)
			state.Step++
		case 9:
			state.BooksCount, _ = strconv.Atoi(message.Text)
			msg.Text = "What supplements do you prefer? Or what supplements would you like to try? "
			App.Bot.Send(msg)
			state.Step++
		case 10:
			state.PreferringSupplements = message.Text
			msg.Text = "What habits would you like to acquire?"
			App.Bot.Send(msg)
			state.Step++
		case 11:
			state.HabitsToAcquire = message.Text
			msg.Text = "What news categories are you interested in? If you are not interested in news, just type 'none'."
			App.Bot.Send(msg)
			state.Step++
		case 12:
			if message.Text == "none" {
				msg.Text = "What time do you usually wake up?"
				App.Bot.Send(msg)
				state.Step++
				continue
			}
			state.NewsCategories = message.Text
			msg.Text = "What time do you usually wake up?"
			App.Bot.Send(msg)
			state.Step++
		case 13:
			state.WakeUpTime, _ = time.Parse("15:04", message.Text)
			msg.Text = "Last question! What time do you usually go to bed?"
			App.Bot.Send(msg)
			state.Step++
		case 14:
			state.BedTime, _ = time.Parse("15:04", message.Text)
			msg.Text = "Thank you for your answers! I will send you a message every day to remind you to drink water, go to bed at time you provided, send you message at your wake up time, do physical exercises, and meditate. I will also send you news articles that you might be interested in. If you want to change your answers, you should re-configure bot from start."
			App.Bot.Send(msg)
		}
	}
}

func (App *Application) ChangePreferableMeditationTime(upd tgbotapi.Update) {}
func (App *Application) ChangePreferableReadingTime(upd tgbotapi.Update)    {}
func (App *Application) ChangePreferableExerciseTime(upd tgbotapi.Update)   {}
func (App *Application) ChangePreferableSleepingTime(upd tgbotapi.Update)   {}

// DailyWaterReminder sends a message to the user with a reminder to drink water
func (App *Application) DailyWaterReminder(channelId int64) {
	// TODO: add preferable amount and time for drinking water (each x hours)

	msg := tgbotapi.NewMessage(channelId, "Don't forget to drink water today!")
	if _, err := App.Bot.Send(msg); err != nil {
		App.Logger.Printf("error sending message: %v\n", err)
		return
	}
}

// DailySleepReminder sends a message to the user with a reminder to sleep
func (App *Application) DailySleepReminder(channelId int64) {
	// TODO: add sleep recommendations
	// TODO: add user preferable time for sleep, and send message exactly at that time

	msg := tgbotapi.NewMessage(channelId, "Don't forget to sleep well today!")
	if _, err := App.Bot.Send(msg); err != nil {
		App.Logger.Printf("error sending message: %v\n", err)
		return
	}
}

// DailyMeditationReminder sends a message to the user with a reminder to meditate
func (App *Application) DailyMeditationReminder(channelId int64) {
	// TODO: add meditation recommendations
	// TODO: add user preferable time for meditation, and send message exactly at that time

	msg := tgbotapi.NewMessage(channelId, "Don't forget to meditate today!")
	if _, err := App.Bot.Send(msg); err != nil {
		App.Logger.Printf("error sending message: %v\n", err)
		return
	}
}

// DailyExerciseReminder sends a message to the user with a reminder to exercise
func (App *Application) DailyExerciseReminder(channelId int64) {
	// TODO: add exercise plans and recommendations
	// TODO: add user preferable time for training, create custom plan, and send message exactly at that time

	msg := tgbotapi.NewMessage(channelId, "Don't forget to exercise today!")
	if _, err := App.Bot.Send(msg); err != nil {
		App.Logger.Printf("error sending message: %v\n", err)
		return
	}
}

func (App *Application) DailyMorningMotivationalQuote(channelId int64) {
	quote, err := App.DB.GetRandomQuote()
	if err != nil {
		App.Logger.Printf("error getting random quote: %v\n", err)
		msg := tgbotapi.NewMessage(channelId, "Error getting random quote")
		_, _ = App.Bot.Send(msg)
	}

	msg := tgbotapi.NewMessage(channelId, quote)
	_, _ = App.Bot.Send(msg)
}

func (App *Application) AddSubscriptionOnMorningNews(update tgbotapi.Update) {
	// TODO: implement categories for news, and add ability to choose them

}

// Sending Recommendations (Manual, using commands)

func (App *Application) ManuallySendReadingRecommendations(update tgbotapi.Update) {
}

// ManuallySendExerciseRecommendations sends exercise recommendations to the user
func (App *Application) ManuallySendExerciseRecommendations(update tgbotapi.Update) {

	updates := App.UpdatesChannel

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
	)
	msg.ReplyMarkup = keyboard
	m, err := App.Bot.Send(msg)
	if err != nil {
		App.Logger.Printf("error sending message: %v\n", err)
	}

	for update := range updates {
		if update.CallbackQuery != nil {
			callbackData := update.CallbackQuery.Data
			messageID := update.CallbackQuery.Message.MessageID
			chatID := update.CallbackQuery.Message.Chat.ID

			if chatID == m.Chat.ID && messageID == m.MessageID {
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

				}
			}
			editMsg := tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, tgbotapi.InlineKeyboardMarkup{})
			_, _ = App.Bot.Send(editMsg)
			return
		}
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

//  this functions should be using chat gpt api however i do not want to pay for it
//func (App *Application) ManuallySendRecommendationsForCustomGoals(update tgbotapi.Update) {
//}
//func (App *Application) ManuallySendRecommendationsForCustomHabits(update tgbotapi.Update) {
//}

// Sending Recommendations (Automate)

// SendSupplementsRecommendations is function which sends recommendations for intake supplements based on category
// provided by user
func (App *Application) SendSupplementsRecommendations(update tgbotapi.Update) {
	category := strings.TrimPrefix("/supplements", update.Message.Text)
	supplements, err := utils.RequestToIHerb(category)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
		_, _ = App.Bot.Send(msg)
		return
	}

	text := fmt.Sprintf(`
				Based on provided category: %s
				We recommend you following supplements:
			`, category)

	for i := 0; i <= len(supplements)-1; i++ {
		text += fmt.Sprintf("Title: %s \n Price: %s \n", supplements[i].Title, supplements[i].Price)
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	_, _ = App.Bot.Send(msg)
}

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

// CreateCustomSupplementIntakePlan function will help user to choose
func (App *Application) CreateCustomSupplementIntakePlan(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "What categories of supplements would you like to take? Please, input 1 by 1 \n Here is list: \n")
	for i := 0; i <= len(AvaliableSupplementCategories)-1; i++ {
		msg.Text += AvaliableSupplementCategories[i] + "\n"
	}

	_, _ = App.Bot.Send(msg)

	updates := App.UpdatesChannel

	for update := range updates {
		if utils.VerifyCategory(update.Message.Text) {
			items, err := utils.RequestToIHerb(update.Message.Text)
			if err != nil {
				msg.Text = "something went wrong. please, try again later"
				_, _ = App.Bot.Send(msg)
				return
			}
			msg.Text = ""
			for i := 0; i <= len(items)-1; i++ {
				msg.Text += "Title: " + items[i].Title + "\n"
				msg.Text += "Price: " + items[i].Price + "\n"
				msg.Text += "\n"
			}
			_, _ = App.Bot.Send(msg)
			return
		} else {
			msg.Text = ""
			msg.Text = "please, choose from categories"
			for i := 0; i <= len(AvaliableSupplementCategories)-1; i++ {
				msg.Text += AvaliableSupplementCategories[i] + "\n"
			}
			return
		}
	}
}

// GenerateWorkoutForParticularMuscleGroup creates workout for this muscle groups:
// []string{"back", "biceps", "chest", "triceps", "legs", "shoulders"}
func (App *Application) GenerateWorkoutForParticularMuscleGroup(update tgbotapi.Update) {
	muscle := strings.Split(update.Message.Text, " ")
	muscle = muscle[1:]
	count := muscle[1]

	c, err := strconv.Atoi(count)
	if err != nil {
		App.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "error converting"+count+"to int"))
		return
	}
	exercises, err := App.DB.GenerateXRandomExercises(muscle[0], c)
	if err != nil {
		App.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "error generating workouts"))
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	for i := 0; i <= len(exercises)-1; i++ {
		msg.Text = fmt.Sprintf("%s \n %s \n %s \n", exercises[i].Title, exercises[i].Technique, exercises[i].VideoURI)
	}

	_, _ = App.Bot.Send(msg)
}

// CreateWeekTrainingPlan actually creates training plan for the whole week!! <33
func (App *Application) CreateWeekTrainingPlan(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	workouts := []string{"back", "biceps", "chest", "triceps", "legs", "shoulders"}
	for i := 0; i < len(workouts); i = i + 2 {
		exercises, err := App.DB.GenerateXRandomExercises(workouts[i], 4)
		if err != nil {
			msg.Text = "something went wrong while generating workout. please, try again later"
			_, _ = App.Bot.Send(msg)
			return
		}
		msg.Text += "Exercises for " + workouts[i] + workouts[i+1]
		for i := 0; i <= len(exercises)-1; i++ {
			msg.Text += fmt.Sprintf("%s \n %s \n %s \n", exercises[i].Title, exercises[i].Technique, exercises[i].VideoURI)
		}

		exercises, err = App.DB.GenerateXRandomExercises(workouts[i+1], 3)
		if err != nil {
			msg.Text = "something went wrong while generating workout. please, try again later"
			_, _ = App.Bot.Send(msg)
			return
		}

		for i := 0; i <= len(exercises)-1; i++ {
			msg.Text += fmt.Sprintf("%s \n %s \n %s \n", exercises[i].Title, exercises[i].Technique, exercises[i].VideoURI)
		}
		_, _ = App.Bot.Send(msg)
	}
}

// GetRandomMotivationalQuote retrieves random motivational quote from database
func (App *Application) GetRandomMotivationalQuote(update tgbotapi.Update) {
	quote, err := App.DB.GetRandomQuote()
	if err != nil {
		App.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "error generating motivational quote"))
		return
	}

	App.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, quote))
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
	msg.Text = "Succesfully added new admin."
	App.Bot.Send(msg)
}

// ReportToAdmins function to report issues with bot
func (App *Application) ReportToAdmins(update tgbotapi.Update) {
	problem := strings.TrimPrefix("/report", update.Message.Text)

	adminId, err := App.DB.GetRandomAdmin()
	if err != nil {
		App.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Something went wrong.. Sending report to admins"))
		return
	}

	// TODO: create table for storing user issues
	App.Bot.Send(tgbotapi.NewMessage(adminId, fmt.Sprintf(
		`User %d with username [%s] has reported following issue: %s`,
		update.Message.From.ID,
		update.Message.From.UserName,
		problem)))
}

// CreateCustomMotivationalReminder saves user answer to the question 'What is your motivation?'
// it will be used to once a day remind user about his/her motivation
func (App *Application) CreateCustomMotivationalReminder(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, `
			Each of us has a someone, or something that push us to move forward. Not to give up. To keep going, fighting. 
			What is yours? This message will be displayed as 'Do not forget about <X>' or 'Do it because <X>'. Please, send reply to this message with your answer. `)

	sentMsg, _ := App.Bot.Send(msg)

	updates := App.UpdatesChannel

	var reason string

	for upd := range updates {
		if upd.Message != nil && upd.Message.ReplyToMessage != nil && upd.Message.ReplyToMessage.MessageID == sentMsg.MessageID {
			reason = upd.Message.Text
			break
		} else {
			msg.Text = "Please, reply to this message with your answer. (and call function again)"
			App.Bot.Send(msg)
			return
		}
	}

	err := App.DB.SaveUserCustomMotivation(update.Message.From.ID, reason)
	if err != nil {
		msg.Text = "Error saving your motivation. Please, try again later."
		App.Bot.Send(msg)
	}

	msg.Text = "Your motivation has been saved. Thank you!"
	App.Bot.Send(msg)
}

// RemindForWhatItIsFor functions periodically sends message to user for motivating him/her
func (App *Application) RemindForWhatItIsFor(update tgbotapi.Update) {
	quote, err := App.DB.GetUserCustomMotivation(update.Message.From.ID)
	if err != nil {
		App.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Something went wrong: "+err.Error()+"\n Please, report to admins by using command /report {<-issue->}"))
		return
	}

	App.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Do not forget about: "+quote))
}

// RequestToUploadChallenge function sends request to admins approve or deny custom challenge to be stored in database
// and used in future
func (App *Application) RequestToUploadChallenge(update tgbotapi.Update) {
	adminId, err := App.DB.GetRandomAdmin()
	if err != nil {
		App.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error getting admins id. Please, report error using command /report <issue> "))
		return
	}
	challenge := strings.TrimPrefix(update.Message.Text, "/upload-challenge")
	msg := tgbotapi.NewMessage(adminId, "Request to upload custom challenge. "+challenge)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅Verify", "v"),
			tgbotapi.NewInlineKeyboardButtonData("❌Deny", "d"),
		),
	)
	msg.ReplyMarkup = keyboard
	_, _ = App.Bot.Send(msg)

	updates := App.UpdatesChannel

	for upd := range updates {
		if upd.CallbackQuery != nil && upd.Message.Chat.ID == adminId {
			callbackData := upd.CallbackQuery.Data

			switch callbackData {
			case "v":
				err = App.DB.UploadChallenge(challenge)
				if err != nil {
					App.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error uploading challenge. Please, report error using command /report <issue> "))
					return
				}
			case "d":
				App.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Challenge has been denied."))
			}
		}
	}

}

// CreateCustomChallenge creates request
func (App *Application) CreateCustomChallenge(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, `
			Please, send reply to this message with your challenge. 
			Example: 'Do 10 push-ups every day' or 'Do not eat sweets for a week'`)

	sentMsg, _ := App.Bot.Send(msg)

	updates := App.UpdatesChannel

	var challenge string

	for upd := range updates {
		if upd.Message != nil && upd.Message.ReplyToMessage != nil && upd.Message.ReplyToMessage.MessageID == sentMsg.MessageID {
			challenge = upd.Message.Text
			break
		} else {
			msg.Text = "Please, reply to this message with your challenge. (and call function again)"
			App.Bot.Send(msg)
			return
		}
	}

	err := App.DB.SaveUserCustomChallenge(update.Message.From.ID, challenge)
	if err != nil {
		msg.Text = "Error saving your challenge. Please, try again later."
		App.Bot.Send(msg)
	}

	msg.Text = "Your challenge has been saved. Thank you!"
	App.Bot.Send(msg)
}

func (App *Application) ReviewChallenge(update tgbotapi.Update) {}

func (App *Application) CreateCustomHabit(update tgbotapi.Update)                    {}
func (App *Application) SendChallenges(update tgbotapi.Update)                       {}
func (App *Application) SendReadingRecommendations(update tgbotapi.Update)           {}
func (App *Application) SendExerciseRecommendations(update tgbotapi.Update)          {}
func (App *Application) SendSleepingRecommendations(update tgbotapi.Update)          {}
func (App *Application) SendRecommendationsForMentalHealth(update tgbotapi.Update)   {}
func (App *Application) SendRecommendationsForPhysicalHealth(update tgbotapi.Update) {}
func (App *Application) ListCustomGoals(update tgbotapi.Update)                      {}
func (App *Application) ListCustomChallenges(update tgbotapi.Update)                 {}
func (App *Application) ListCustomHabits(update tgbotapi.Update)                     {}
func (App *Application) ListCustomReminders(update tgbotapi.Update)                  {}
func (App *Application) CreateCustomRoutine(update tgbotapi.Update)                  {}
func (App *Application) CreateCustomMealPreparingPlan(update tgbotapi.Update)        {}
func (App *Application) CreateCustomExerciseForWorkoutPlan(update tgbotapi.Update)   {}
func (App *Application) AddCustomMealToExistingPlan(update tgbotapi.Update)          {}
func (App *Application) SendRandomMorningChallenge(update tgbotapi.Update)           {}
func (App *Application) RemindAboutDoingCustomHabit(update tgbotapi.Update)          {}

// Subscription management
func (App *Application) BuySubscription(update tgbotapi.Update)     {}
func (App *Application) RemoveSubscription(update tgbotapi.Update)  {}
func (App *Application) ManageUsers(update tgbotapi.Update)         {}
func (App *Application) ManageSubscriptions(update tgbotapi.Update) {}

// Optional
// author info
// donations (crypto donations) + statistic
