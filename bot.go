package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
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
	Bot            *tgbotapi.BotAPI
	DB             *models.DatabaseModel
	Logger         *log.Logger
	Config         Config
	UpdatesChannel tgbotapi.UpdatesChannel
	Timers         map[int64]map[string]time.Time
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

	updatesChannel := bot.ListenForWebhook(cfg.ServerURI + "/webhook" + bot.Token)

	app := &Application{
		Bot:            bot,
		Logger:         Logger,
		DB:             models.NewDatabaseModel(db),
		Config:         cfg,
		UpdatesChannel: updatesChannel,
	}

	err = app.DB.CreateMuscleGroupTables()
	if err != nil {
		app.Logger.Printf("Error creating muscle group tables: %v\n", err)
		return
	}
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

	app.Bot.Debug = true

	r := gin.Default()

	r.POST("/webhook"+bot.Token, func(c *gin.Context) {
		var update tgbotapi.Update
		if err := c.BindJSON(&update); err != nil {
			Logger.Printf("error binding update: %v\n", err)
			return
		}
	})

	server := &http.Server{
		Addr:    app.Config.ServerURI,
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			Logger.Fatalf("error starting server: %v\n", err)
		}
	}()

	log.Printf("bot is listening on %s\n", app.Config.ServerURI+"/webhook"+bot.Token)
	log.Printf("bot is authorized on account %s\n", bot.Self.UserName)

	select {}
}
