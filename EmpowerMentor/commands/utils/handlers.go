package utils

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"self-improvement-bot/models"
	"strings"
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
	AvaliableNewsCategories  = RetrieveCategories()
	ConfigureStates          = make(map[int64]models.UserBotConfiguration)
	SpoonocularConfiguration = make(map[int64]models.SpoonocularConfiguration)
	muscles                  = []string{"back", "chest", "biceps", "triceps", "shoulders", "legs", "glutes"}
	DietNames                = []string{
		"Gluten Free",
		"Ketogenic",
		"Vegetarian",
		"Lacto-Vegetarian",
		"Ovo-Vegetarian",
		"Vegan",
		"Pescetarian",
		"Paleo",
		"Primal",
		"Low FODMAP",
		"Whole30",
	}
	Allergens = []string{
		"Milk",
		"Eggs",
		"Fish",
		"Crustacean shellfish",
		"Tree nuts",
		"Peanuts",
		"Wheat",
		"Soybeans",
		"Mustard",
		"Sesame",
		"Sulfites",
		"Celery",
		"Lupin",
		"Mollusks",
		"Gluten",
	}
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

// GetRandomMotivationalQuote retrieves random motivational quote from database
func (App *Application) GetRandomMotivationalQuote(update tgbotapi.Update) {
	quote, err := App.DB.GetRandomQuote()
	if err != nil {
		App.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "error generating motivational quote"))
		return
	}

	App.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, quote))
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

// GetCustomUserChallenge retrieves random user custom challenge
func (App *Application) GetCustomUserChallenge(update tgbotapi.Update) {
	challenge, err := App.DB.RetrieveUserRandomCustomChallenge(update.Message.From.ID)
	if err != nil {
		App.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error getting your challenge. Please, report error using command /report <issue> "))
		return
	}

	App.Bot.Send(App.TranslateToUkrainian(tgbotapi.NewMessage(update.Message.Chat.ID, challenge)))
}

// TranslateToUkrainian translates message to Ukrainian language
func (App *Application) TranslateToUkrainian(message tgbotapi.MessageConfig) tgbotapi.MessageConfig {
	message.Text = TranslateToUkrainian(message.Text)
	return message
}
