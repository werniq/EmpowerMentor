package utils

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

// SendSupplementsRecommendations is function which sends recommendations for intake supplements based on category
// provided by user
func (App *Application) SendSupplementsRecommendations(update tgbotapi.Update) {
	category := strings.TrimPrefix("/supplements ", update.Message.Text)
	if StringInArray(category, AvaliableSupplementCategories) == false {
		supplements, err := RequestToIHerb(category)
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

// CreateCustomSupplementIntakePlan function will help user to choose
func (App *Application) CreateCustomSupplementIntakePlan(update tgbotapi.Update) {
	if StringInArray(strings.TrimPrefix(update.Message.Text, "/supplements-plan "), AvaliableSupplementCategories) == false {
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

	if VerifyCategory(category) {
		items, err := RequestToIHerb(category)
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
