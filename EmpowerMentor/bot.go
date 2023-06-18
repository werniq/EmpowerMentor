package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"self-improvement-bot/commands/utils"
	"self-improvement-bot/driver"
	"self-improvement-bot/models"
	"time"
)

var (
	Logger = log.New(os.Stdout, "ERROR: \t", log.Ldate|log.Ltime|log.Lshortfile)
)

type Config struct {
	ServerURI string
}

type Application struct {
	Bot    *tgbotapi.BotAPI
	DB     *models.DatabaseModel
	Logger *log.Logger
	Config Config
	Timers map[int64]map[string]time.Time
	Stripe struct {
		PublishableKey string
		SecretKey      string
	}
}

func init() {
	app := &Application{}
	db, err := driver.OpenDb()
	if err != nil {
		Logger.Printf("error opening database: %v\n", err)
		return
	}

	defer db.Close()

	app.DB = models.NewDatabaseModel(db)

	err = app.DB.TruncateTableChallenges()
	if err != nil {
		Logger.Printf("error truncating challenges: %v\n", err)
		return
	}
	err = app.DB.TruncateTableMotivationQuotes()
	if err != nil {
		Logger.Printf("error truncating motivational quotes: %v\n", err)
		return
	}

	err = app.DB.TruncateTables()
	if err != nil {
		Logger.Printf("error truncating tables: %v\n", err)
		return
	}

	err = app.DB.TruncateMotivationalQuotes()
	if err != nil {
		Logger.Printf("error truncating motivational quotes: %v\n", err)
		return
	}
	fmt.Printf("Tables truncated successfully\n")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		Logger.Printf("error loading .env file: %v\n", err)
		return
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		Logger.Printf("error creating new bot: %v\n", err)
		return
	}

	cfg := Config{
		ServerURI: os.Getenv("SERVER_URI"),
	}

	db, err := driver.OpenDb()
	if err != nil {
		Logger.Printf("error opening database: %v\n", err)
		return
	}

	webhookURL := cfg.ServerURI + "/webhook" + bot.Token

	// sets up webhook for bot
	_, err = http.Post(fmt.Sprintf("https://api.telegram.org/bot%s/setWebhook?url=%s", bot.Token, webhookURL), "application/json", nil)
	if err != nil {
		Logger.Printf("error setting webhook: %v\n", err)
		return
	}

	app := utils.NewApplication(bot, models.NewDatabaseModel(db), Logger, utils.Config(cfg))

	err = app.DB.UploadMotivationalQuotes()
	if err != nil {
		app.Logger.Printf("Error uploading motivational quotes: %v\n", err)
		return
	}

	err = app.DB.UploadDifferentChallenges()
	if err != nil {
		app.Logger.Printf("Error uploading challenges to database: %v\n", err)
		return
	}

	err = app.DB.DropTables()
	if err != nil {
		Logger.Printf("error dropping tables: %v\n", err)
		return
	}

	err = app.DB.CreateMuscleGroupTables()
	if err != nil {
		Logger.Printf("error creating muscle group tables: %v\n", err)
		return
	}

	err = app.DB.StoreAllChestExercises()
	if err != nil {
		app.Logger.Printf("Error uploading data to %s table", "chest")
		return
	}
	err = app.DB.StoreAllBackExercises()
	if err != nil {
		app.Logger.Printf("Error uploading data to %s table", "back")
		return
	}
	err = app.DB.SaveAllLegsExercises()
	if err != nil {
		app.Logger.Printf("Error uploading data to %s table", "legs")
		return
	}
	err = app.DB.StoreAllBicepsExercises()
	if err != nil {
		app.Logger.Printf("Error uploading data to %s table", "biceps")
		return
	}
	err = app.DB.StoreAllGlutesExercises()
	if err != nil {
		app.Logger.Printf("Error uploading data to %s table", "glutes")
		return
	}
	err = app.DB.SaveAllShouldersExercises()
	if err != nil {
		app.Logger.Printf("Error uploading data to %s table", "shoulders")
		return
	}

	app.Bot.Debug = true

	r := gin.Default()

	app.ConfigureRoutes(r)

	fmt.Printf("bot is listening on %s\n", app.Config.ServerURI+"/webhook"+bot.Token)

	if err = r.Run(":8080"); err != nil {
		app.Logger.Printf("error running server at port :8080")
		return
	}

	log.Printf("bot is listening on %s\n", app.Config.ServerURI+"/webhook"+bot.Token)
	log.Printf("bot is authorized on account %s\n", bot.Self.UserName)

	select {}
}
