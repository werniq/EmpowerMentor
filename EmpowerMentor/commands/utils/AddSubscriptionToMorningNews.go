package utils

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

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
