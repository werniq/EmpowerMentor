package commands

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"math/rand"
	"os"
	"self-improvement-bot/commands/utils"
	"self-improvement-bot/models"
	"strconv"
	"strings"
	"time"
)

type Application struct {
	Bot    *tgbotapi.BotAPI
	DB     *models.DatabaseModel
	Logger *log.Logger
	Config Config
	Stripe struct {
		PublishableKey string
		SecretKey      string
	}
	Spoonocular struct {
		ApiKey string
	}
}

type Config struct {
	ServerURI string
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
	AvailableMuscleGroup     = []string{"back", "biceps", "chest", "triceps", "legs", "shoulders"}
	AvaliableNewsCategories  = utils.RetrieveCategories()
	ConfigureStates          = make(map[int64]models.UserBotConfiguration)
	SpoonocularConfiguration = make(map[int64]models.SpoonocularConfiguration)
	muscles                  = []string{"back", "chest", "biceps", "triceps", "shoulders", "legs", "glutes"}
)

func NewApplication(bot *tgbotapi.BotAPI, db *models.DatabaseModel, logger *log.Logger, cfg Config) *Application {
	return &Application{
		Bot:    bot,
		DB:     db,
		Logger: logger,
		Config: cfg,
		Stripe: struct {
			PublishableKey string
			SecretKey      string
		}{
			PublishableKey: os.Getenv("STRIPE_PUBLIC_KEY"),
			SecretKey:      os.Getenv("STRIPE_SECRET_KEY")},
		Spoonocular: struct {
			ApiKey string
		}{
			ApiKey: os.Getenv("SPOONOCULAR_API_KEY"),
		},
	}
}

// Start is a handler for /start command
func (App *Application) Start(upd tgbotapi.Update) {
	App.AboutAuthor(upd)
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "I am glad to see you here! Hope this bot will help in your journey to a better life!")
	if _, err := App.Bot.Send(msg); err != nil {
		App.Logger.Printf("error sending message: %v\n", err)
	}
	err := App.DB.SaveMessage(upd)
	if err != nil {
		App.Logger.Printf("error saving message into database: %v\n", err)
		return
	}
}

func (App *Application) Help(upd tgbotapi.Update) {
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "Help")
	if _, err := App.Bot.Send(msg); err != nil {
		App.Logger.Printf("error sending message: %v\n", err)
	}
	err := App.DB.SaveMessage(upd)
	if err != nil {
		App.Logger.Printf("error saving message into database: %v\n", err)
		return
	}
}

func (App *Application) AboutAuthor(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.Text = `My name is Oleksandr Matviienko, link to my social medias you can find at https://werniq.github.io
				This 6-months I have gained ~15 kilograms of clean muscle mass and I want to share my experience with you.
				The main purpose of this bot is to help you to achieve your goals in fitness and health.
				It will motivate you through the day and will help you to track your progress.
				If you will listen to his advices, you will be able to achieve your goals in fitness and health.
				A bit about those habits.
					1. Meditation. It will help you to calm down and to focus on your goals. One time a day you should take a deep breath, calm down, and think about your goals. It will help you to stay motivated.
					2. Drink water. It will help you to stay hydrated and to keep your body in a good shape. As well as it will your cognitive abilities.
					3. Sleep. It will help you to recover after a hard day and to stay productive during the day.
					4. Exercise. It will help you to stay in a good shape and to keep your body healthy. Exercising creates some discipline in your life and it will help you to achieve your goals.
					5. Reading. Reading benefits include acquiring knowledge, stimulating the mind, expanding vocabulary, improving focus, reducing stress, fostering empathy, enhancing writing skills, boosting memory, stimulating imagination, and promoting personal growth.
				Hope, that this bot will help you to become more disciplined, organized, healthy and motivated`

	if _, err := App.Bot.Send(msg); err != nil {
		App.Logger.Printf("error sending message: %v\n", err)
	}
	err := App.DB.SaveMessage(update)
	if err != nil {
		App.Logger.Printf("error saving message into database: %v\n", err)
		return
	}
}

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

// WalkingReminder sends message to user if it's time to go for a walk
func (App *Application) WalkingReminder(upd tgbotapi.Update) {
	if time.Now().Hour() == 12 && time.Now().Minute() == 30 {
		App.Bot.Send(tgbotapi.NewMessage(upd.Message.Chat.ID, "It's time to go for a walk!"))
	}
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
	category := strings.TrimPrefix("/supplements ", update.Message.Text)
	if utils.StringInArray(category, AvaliableSupplementCategories) == false {
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
	} else {
		App.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Please provide valid category, Available categories you can view by typing /supplement-categories"))
	}
}

// ListAvailableSupplementCategories is function which sends list of available categories for supplements
func (App *Application) ListAvailableSupplementCategories(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	text := "Available categories: \n"
	for i := 0; i <= len(AvaliableSupplementCategories)-1; i++ {
		text += fmt.Sprintf("%s \n", AvaliableSupplementCategories[i])
	}

	_, err := App.Bot.Send(msg)
	if err != nil {
		App.Logger.Printf("error sending message: %v\n", err)
		return
	}
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
	if utils.StringInArray(strings.TrimPrefix(update.Message.Text, "/supplements-plan "), AvaliableSupplementCategories) == false {
		fmt.Println(strings.TrimPrefix(update.Message.Text, "/supplements-plan"))

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "What categories of supplements would you like to take? Please, input 1 by 1 \n Here is list: \n")
		for i := 0; i <= len(AvaliableSupplementCategories)-1; i++ {
			msg.Text += AvaliableSupplementCategories[i] + "\n\t\t"
		}

		_, _ = App.Bot.Send(msg)
		return
	}

	category := strings.TrimPrefix(update.Message.Text, "/supplements-plan")
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	if utils.VerifyCategory(category) {
		items, err := utils.RequestToIHerb(category)
		if err != nil {
			msg.Text = "something went wrong. please, try again later"
			_, _ = App.Bot.Send(msg)
			return
		}

		for i := 0; i <= len(items)-1; i++ {
			msg.Text += "Title: " + items[i].Title + "\n"
			msg.Text += "Price: " + items[i].Price + "\n"
			msg.Text += "\n"
		}
		_, _ = App.Bot.Send(msg)
		return
	}
}

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
	if utils.StringInArray(m, muscles) == false {
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
	msg.Text = "Successfully added new admin."
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

// RemindForWhatItIsFor functions periodically sends message to user for motivating him/her
func (App *Application) RemindForWhatItIsFor(update tgbotapi.Update) {
	if time.Now().Hour() == 12 && time.Now().Minute() == 0 {
		quote, err := App.DB.GetUserCustomMotivation(update.Message.From.ID)
		if err != nil {
			App.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Something went wrong: "+err.Error()+"\n Please, report to admins by using command /report {<-issue->}"))
			return
		}
		App.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Do not forget about: "+quote))
	}
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

	err = App.DB.ChallengeUploadRequest(update.Message.From.ID, challenge)
	if err != nil {
		App.Logger.Printf("error uploading challenge request: %s", err.Error())
		return
	}
}

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

// GetCustomUserChallenge retrieves random user custom challenge
func (App *Application) GetCustomUserChallenge(update tgbotapi.Update) {
	challenge, err := App.DB.RetrieveUserRandomCustomChallenge(update.Message.From.ID)
	if err != nil {
		App.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error getting your challenge. Please, report error using command /report <issue> "))
		return
	}

	App.Bot.Send(App.TranslateToUkrainian(tgbotapi.NewMessage(update.Message.Chat.ID, challenge)))
}

// SendChallenges sends challenges to users
func (App *Application) SendChallenges(update tgbotapi.Update) {
	if time.Now().Hour() == 11 && time.Now().Minute() == 0 {
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
	if time.Now().Hour() == 20 && time.Now().Minute() == 0 {
		num := rand.Intn(len(sleepingRecommendations))
		App.Bot.Send(App.TranslateToUkrainian(tgbotapi.NewMessage(update.Message.Chat.ID, sleepingRecommendations[num])))
	}
}

// SendRecommendationsForMentalHealth sends recommendations to users for mental health
func (App *Application) SendRecommendationsForMentalHealth(update tgbotapi.Update) {
	if time.Now().Hour() == 20 && time.Now().Minute() == 0 {
		num := rand.Intn(len(mentalHealthRecommendations))
		App.Bot.Send(App.TranslateToUkrainian(tgbotapi.NewMessage(update.Message.Chat.ID, mentalHealthRecommendations[num])))
	}
}

// AddSubscriptionOnMorningNews adds subscription on morning news
func (App *Application) AddSubscriptionOnMorningNews(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "What categories do you want to receive news about?")

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("health🏥", "health"),
			tgbotapi.NewInlineKeyboardButtonData("fitness💪", "fitness"),
			tgbotapi.NewInlineKeyboardButtonData("personal-growth🌱", "personal"),
			tgbotapi.NewInlineKeyboardButtonData("psychology 🧠 ", "psychology"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("mindfulness🧘‍♂️", "mindfulness"),
			tgbotapi.NewInlineKeyboardButtonData("self-care🛀", "self-Care"),
			tgbotapi.NewInlineKeyboardButtonData("motivation🛀", "motivation"),
			tgbotapi.NewInlineKeyboardButtonData("productivity 📊", "productivity"),
			tgbotapi.NewInlineKeyboardButtonData("happiness😊", "happiness"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("relationships💑", "relationships"),
			tgbotapi.NewInlineKeyboardButtonData("career-development🚀", "career-development"),
			tgbotapi.NewInlineKeyboardButtonData("leadership🎯", "leadership"),
			tgbotapi.NewInlineKeyboardButtonData("education📚", "education"),
			tgbotapi.NewInlineKeyboardButtonData("self-Help🤝", "self-Help"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("mental‍👩‍⚕️", "mental"),
			tgbotapi.NewInlineKeyboardButtonData("well-being🌞", "well-being"),
		),
	)

	msg.ReplyMarkup = keyboard
	msg = App.TranslateToUkrainian(msg)
	App.Bot.Send(msg)
}

// AddSubscriptionOnMorningNewsCallback adds subscription on morning news callback
// TODO: implement function to actually send morning newsletter based on user preferences
func (App *Application) AddSubscriptionOnMorningNewsCallback(update tgbotapi.Update) {
	// assume that update.CallbackQuery.Data == "health" || or any available category from list
	m := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")
	data := update.CallbackQuery.Data
	if data == "health" || data == "fitness" || data == "personal" || data == "growth" || data == "psychology" || data == "mindfulness" || data == "self-Care" || data == "motivation" || data == "productivity" || data == "happiness" || data == "relationships" || data == "career-development" || data == "leadership" || data == "education" || data == "self-Help" || data == "mental" || data == "well-being" {
		err := App.DB.StoreNewsCategories(data, update.CallbackQuery.From.ID)
		if err != nil {
			m.Text = "Error adding subscription. Please, try again later."
			App.Bot.Send(m)
			return
		}

		m.Text = "Your subscription has been added. Thank you!"
		App.Bot.Send(m)
		return
	}
}

// CreateCustomMealPreparingPlan requires few arguments from user to create custom meal preparing plan
func (App *Application) CreateCustomMealPreparingPlan(update tgbotapi.Update) {
	config, exists := SpoonocularConfiguration[update.Message.From.ID]
	if !exists {
		config = models.SpoonocularConfiguration{}
		SpoonocularConfiguration[update.Message.From.ID] = config
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "2. What is your time frame? (day/week)")

	if config.Step == 0 {
		msg.Text = "It is time to create your custom meal preparing plan. Please, answer the following questions: \n\n1. What is your target calories? (e.g. 2000)"
		App.Bot.Send(msg)
		config.Step++
		SpoonocularConfiguration[update.Message.From.ID] = config
		return
	}

	if config.Step == 1 {
		config.TargetCalories = update.Message.Text
		msg.Text = "2. What is your time frame? (day/week)"
		App.Bot.Send(msg)
		config.Step++
		SpoonocularConfiguration[update.Message.From.ID] = config
		return
	}

	if config.Step == 2 {
		config.TimeFrame = update.Message.Text
		msg.Text = "3. What is your diet? (list of diets you can view by typing /diets)"
		App.Bot.Send(msg)
		config.Step++
		SpoonocularConfiguration[update.Message.From.ID] = config
		return
	}

	if config.Step == 3 {
		config.Diet = update.Message.Text
		msg.Text = "4. What do you want to exclude? (e.g. eggs, milk, nuts)"
		App.Bot.Send(msg)
		config.Step++
		SpoonocularConfiguration[update.Message.From.ID] = config
		return
	}

	if config.Step == 4 {
		config.Exclude = update.Message.Text
		msg.Text = "Thank you! Your custom meal preparing plan is ready. Please, type /my-meal-plans to view it."
		SpoonocularConfiguration[update.Message.From.ID] = config
	}

	// send request to spoonacular api
	uri := "https://api.spoonacular.com/mealplanner/generate?apiKey=" + App.Spoonocular.ApiKey
	if config.TimeFrame != "" {
		uri += "&timeFrame=" + config.TimeFrame
	}
	if config.TargetCalories != "" {
		uri += "&targetCalories=" + config.TargetCalories
	}

	if config.Diet != "" {
		uri += "&diet=" + config.Diet
	}

	if config.Exclude != "" {
		uri += "&exclude=" + config.Exclude
	}

	week, err := utils.CreateMealPreparingPlan(uri)
	if err != nil {
		msg.Text = "Error creating meal preparing plan. Please, try again later."
		App.Bot.Send(msg)
		return
	}

	err = App.DB.InsertMealPreparePlan(week, update.Message.From.ID)
	if err != nil {
		msg.Text = "Error inserting meal preparing plan. Please, try again later."
		App.Bot.Send(msg)
		return
	}

	msg.Text = "Thank you! Your custom meal preparing plan is ready. Please, type /my-mealplans to view it."
	App.Bot.Send(msg)
}

// UserMealPlan executes when user types /my-mealplans
func (App *Application) UserMealPlan(update tgbotapi.Update) {
	week, err := App.DB.GetMealPlan(update.Message.From.ID)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Error getting meal preparing plan. Please, try again later.")
		App.Bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Your meal preparing plan: \n\n Monday: ")
	for _, meal := range week.Monday.Meals {
		msg.Text += fmt.Sprintf("%s\n", meal.Title)
		msg.Text += fmt.Sprintf("Ready in minutes: %d\n", meal.ReadyInMinutes)
		msg.Text += fmt.Sprintf("Servings: %d\n", meal.Servings)
		msg.Text += fmt.Sprintf("Source url: %s\n\n", meal.SourceURL)
	}

	msg.Text += "Tuesday: "
	for _, meal := range week.Tuesday.Meals {
		msg.Text += fmt.Sprintf("%s\n", meal.Title)
		msg.Text += fmt.Sprintf("Ready in minutes: %d\n", meal.ReadyInMinutes)
		msg.Text += fmt.Sprintf("Servings: %d\n", meal.Servings)
		msg.Text += fmt.Sprintf("Source url: %s\n\n", meal.SourceURL)
	}

	msg.Text += "Wednesday: "
	for _, meal := range week.Wednesday.Meals {
		msg.Text += fmt.Sprintf("%s\n", meal.Title)
		msg.Text += fmt.Sprintf("Ready in minutes: %d\n", meal.ReadyInMinutes)
		msg.Text += fmt.Sprintf("Servings: %d\n", meal.Servings)
		msg.Text += fmt.Sprintf("Source url: %s\n\n", meal.SourceURL)
	}

	msg.Text += "Thursday: "
	for _, meal := range week.Thursday.Meals {
		msg.Text += fmt.Sprintf("%s\n", meal.Title)
		msg.Text += fmt.Sprintf("Ready in minutes: %d\n", meal.ReadyInMinutes)
		msg.Text += fmt.Sprintf("Servings: %d\n", meal.Servings)
		msg.Text += fmt.Sprintf("Source url: %s\n\n", meal.SourceURL)
	}

	msg.Text += "Friday: "
	for _, meal := range week.Friday.Meals {
		msg.Text += fmt.Sprintf("%s\n", meal.Title)
		msg.Text += fmt.Sprintf("Ready in minutes: %d\n", meal.ReadyInMinutes)
		msg.Text += fmt.Sprintf("Servings: %d\n", meal.Servings)
		msg.Text += fmt.Sprintf("Source url: %s\n\n", meal.SourceURL)
	}

	msg.Text += "Saturday: "
	for _, meal := range week.Saturday.Meals {
		msg.Text += fmt.Sprintf("%s\n", meal.Title)
		msg.Text += fmt.Sprintf("Ready in minutes: %d\n", meal.ReadyInMinutes)
		msg.Text += fmt.Sprintf("Servings: %d\n", meal.Servings)
		msg.Text += fmt.Sprintf("Source url: %s\n\n", meal.SourceURL)
	}

	msg.Text += "Sunday: "
	for _, meal := range week.Sunday.Meals {
		msg.Text += fmt.Sprintf("%s\n", meal.Title)
		msg.Text += fmt.Sprintf("Ready in minutes: %d\n", meal.ReadyInMinutes)
		msg.Text += fmt.Sprintf("Servings: %d\n", meal.Servings)
		msg.Text += fmt.Sprintf("Source url: %s\n\n", meal.SourceURL)
	}

	App.Bot.Send(msg)
}

// GetDifferentDietInfo sends information about different diets
func (App *Application) GetDifferentDietInfo(update tgbotapi.Update) {
	dietName := "Gluten Free"
	desc := "Eliminating gluten means avoiding wheat, barley, rye, and other gluten-containing grains and foods made from them (or that may have been cross contaminated)."

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Diet name: %s\nDescription: %s", dietName, desc))
	App.Bot.Send(msg)

	dietName = "Ketogenic"
	desc = "keto diet is based more on the ratio of fat, protein, and carbs in the diet rather than specific ingredients. Generally speaking, high fat, protein-rich foods are acceptable and high carbohydrate foods are not. The formula we use is 55-80% fat content, 15-35% protein content, and under 10% of carbohydrates."

	msg.Text = fmt.Sprintf("Diet name: %s\nDescription: %s", dietName, desc)
	App.Bot.Send(msg)

	dietName = "Vegetarian"
	desc = "No ingredients may contain meat or meat by-products, such as bones or gelatin."

	msg.Text = fmt.Sprintf("Diet name: %s\nDescription: %s", dietName, desc)
	App.Bot.Send(msg)

	dietName = "Lacto-Vegetarian"
	desc = "All ingredients must be vegetarian and none of the ingredients can be or contain egg."

	msg.Text = fmt.Sprintf("Diet name: %s\nDescription: %s", dietName, desc)
	App.Bot.Send(msg)

	dietName = "Ovo-Vegetarian"
	desc = "All ingredients must be vegetarian and none of the ingredients can be or contain dairy."

	msg.Text = fmt.Sprintf("Diet name: %s\nDescription: %s", dietName, desc)
	App.Bot.Send(msg)

	dietName = "Vegan"
	desc = "No ingredients may contain meat or meat by-products, such as bones or gelatin, nor may they contain eggs, dairy, or honey."

	msg.Text = fmt.Sprintf("Diet name: %s\nDescription: %s", dietName, desc)
	App.Bot.Send(msg)

	dietName = "Pescetarian"
	desc = "Everything is allowed except meat and meat by-products - some pescetarians eat eggs and dairy, some do not."

	msg.Text = fmt.Sprintf("Diet name: %s\nDescription: %s", dietName, desc)
	App.Bot.Send(msg)

	dietName = "Paleo"
	desc = "Allowed ingredients include meat (especially grass fed), fish, eggs, vegetables, some oils (e.g. coconut and olive oil), and in smaller quantities, fruit, nuts, and sweet potatoes. We also allow honey and maple syrup (popular in Paleo desserts, but strict Paleo followers may disagree). Ingredients not allowed include legumes (e.g. beans and lentils), grains, dairy, refined sugar, and processed foods."

	msg.Text = fmt.Sprintf("Diet name: %s\nDescription: %s", dietName, desc)
	App.Bot.Send(msg)

	dietName = "Primal"
	desc = "Very similar to Paleo, except dairy is allowed - think raw and full fat milk, butter, ghee, etc."

	msg.Text = fmt.Sprintf("Diet name: %s\nDescription: %s", dietName, desc)
	App.Bot.Send(msg)

	dietName = "Low FODMAP"
	desc = `FODMAP stands for 'fermentable oligo-, di-, mono-saccharides and polyols". Our ontology knows which foods are considered high in these types of carbohydrates (e.g. legumes, wheat, and dairy products)`

	msg.Text = fmt.Sprintf("Diet name: %s\nDescription: %s", dietName, desc)
	App.Bot.Send(msg)

	dietName = "Whole30"
	desc = "Allowed ingredients include meat, fish/seafood, eggs, vegetables, fresh fruit, coconut oil, olive oil, small amounts of dried fruit and nuts/seeds. Ingredients not allowed include added sweeteners (natural and artificial, except small amounts of fruit juice), dairy (except clarified butter or ghee), alcohol, grains, legumes (except green beans, sugar snap peas, and snow peas)"

	msg.Text = fmt.Sprintf("Diet name: %s\nDescription: %s", dietName, desc)
	App.Bot.Send(msg)
}

// TranslateToUkrainian translates message to Ukrainian language
func (App *Application) TranslateToUkrainian(message tgbotapi.MessageConfig) tgbotapi.MessageConfig {
	message.Text = utils.TranslateToUkrainian(message.Text)
	return message
}
