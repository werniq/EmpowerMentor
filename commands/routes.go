package commands

import (
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"self-improvement-bot"
	"time"
)

func (App *Application) ConfigureRoutes(r *gin.Engine) {
	r.POST("/webhook"+App.Bot.Token, func(c *gin.Context) {
		var update tgbotapi.Update
		if err := c.BindJSON(&update); err != nil {
			main.Logger.Printf("error binding update: %v\n", err)
			return
		}
		if update.Message == nil {
			return
		}

		go func() {
			times, reminders, err := App.DB.GetUserReminders(update.Message.Chat.ID)
			if err != nil {
				App.Logger.Printf("error getting user reminders: %v\n", err)
				return
			}

			// if time.Now() is equal to any of the times in times, send a message
			for i, t := range times {
				if time.Now().Hour() == t.Hour() && time.Now().Minute() == t.Minute() {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "It's time to do "+reminders[i])
					_, err := App.Bot.Send(msg)
					if err != nil {
						App.Logger.Printf("error sending message: %v\n", err)
						return
					}
				}
			}
		}()

		for update := range App.UpdatesChannel {
			switch update.Message.Text {
			case "/setup":
				App.BotConfiguration(update)
			case "/start":
				App.Start(update)
			//case "":
			default:
				App.Default(update)
			}
		}
	})
}
