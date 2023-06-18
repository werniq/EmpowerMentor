package utils

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"self-improvement-bot/models"
)

// CreateCustomMealPreparingPlan requires few arguments from user to create custom meal preparing plan
func (App *Application) CreateCustomMealPreparingPlan(update tgbotapi.Update) {
	var config models.SpoonocularConfiguration
	var exists bool

	config, exists = SpoonocularConfiguration[update.Message.From.ID]
	if !exists {
		config = models.SpoonocularConfiguration{}
		SpoonocularConfiguration[update.Message.From.ID] = config
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "2. What is your time frame? (day/week)")

	msg.Text = "It is time to create your custom meal preparing plan. Please, answer the following questions: \n\n1. What is your target calories? (e.g. 2000)"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("1500", "1500"),
			tgbotapi.NewInlineKeyboardButtonData("1750", "1750"),
			tgbotapi.NewInlineKeyboardButtonData("2000", "2000"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("2250", "2250"),
			tgbotapi.NewInlineKeyboardButtonData("2500", "2500"),
			tgbotapi.NewInlineKeyboardButtonData("3000", "3000"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("3250", "3250"),
			tgbotapi.NewInlineKeyboardButtonData("3500", "3500"),
		),
	)

	msg.ReplyMarkup = keyboard

	App.Bot.Send(msg)
	config.Step++
	SpoonocularConfiguration[update.Message.From.ID] = config
}

func (App *Application) CallbackMealPreparePlan(update tgbotapi.Update) {
	config, exists := SpoonocularConfiguration[update.CallbackQuery.From.ID]
	if !exists {
		config = models.SpoonocularConfiguration{}
		SpoonocularConfiguration[update.CallbackQuery.From.ID] = config
	}
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")

	if config.Step == 1 {
		config.TargetCalories = update.CallbackQuery.Data
		msg.Text = "2. What is your time frame? (day/week)"

		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("day", "day"),
				tgbotapi.NewInlineKeyboardButtonData("week", "week"),
			),
		)
		msg.ReplyMarkup = keyboard

		App.Bot.Send(msg)
		config.Step++
		SpoonocularConfiguration[update.CallbackQuery.From.ID] = config
		return
	}

	if config.Step == 2 {
		config.TimeFrame = update.CallbackQuery.Data
		msg.Text = "3. What type of diet you want to follow?"

		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(DietNames[0], DietNames[0]),
				tgbotapi.NewInlineKeyboardButtonData(DietNames[1], DietNames[1]),
				tgbotapi.NewInlineKeyboardButtonData(DietNames[2], DietNames[2]),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(DietNames[3], DietNames[3]),
				tgbotapi.NewInlineKeyboardButtonData(DietNames[4], DietNames[4]),
				tgbotapi.NewInlineKeyboardButtonData(DietNames[5], DietNames[5]),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(DietNames[6], DietNames[6]),
				tgbotapi.NewInlineKeyboardButtonData(DietNames[7], DietNames[7]),
				tgbotapi.NewInlineKeyboardButtonData(DietNames[8], DietNames[8]),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(DietNames[9], DietNames[9]),
				tgbotapi.NewInlineKeyboardButtonData(DietNames[10], DietNames[10]),
			),
		)

		msg.ReplyMarkup = keyboard
		App.Bot.Send(msg)

		config.Step++
		SpoonocularConfiguration[update.CallbackQuery.From.ID] = config
		return
	}

	if config.Step == 3 {
		config.Diet = update.CallbackQuery.Data
		msg.Text = "4. What do you want to exclude? (e.g. eggs, milk, nuts) P.S. For now, we support only those allergens: "
		for _, v := range Allergens {
			msg.Text += v + ", "
		}

		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(Allergens[0], Allergens[0]),
				tgbotapi.NewInlineKeyboardButtonData(Allergens[1], Allergens[1]),
				tgbotapi.NewInlineKeyboardButtonData(Allergens[2], Allergens[2]),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(Allergens[3], Allergens[3]),
				tgbotapi.NewInlineKeyboardButtonData(Allergens[4], Allergens[4]),
				tgbotapi.NewInlineKeyboardButtonData(Allergens[5], Allergens[5]),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(Allergens[6], Allergens[6]),
				tgbotapi.NewInlineKeyboardButtonData(Allergens[7], Allergens[7]),
				tgbotapi.NewInlineKeyboardButtonData(Allergens[8], Allergens[8]),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(Allergens[9], Allergens[9]),
				tgbotapi.NewInlineKeyboardButtonData(Allergens[10], Allergens[10]),
				tgbotapi.NewInlineKeyboardButtonData(Allergens[11], Allergens[11]),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(Allergens[12], Allergens[12]),
				tgbotapi.NewInlineKeyboardButtonData(Allergens[13], Allergens[13]),
				tgbotapi.NewInlineKeyboardButtonData(Allergens[14], Allergens[14]),
			),
		)
		msg.ReplyMarkup = keyboard

		App.Bot.Send(msg)
		config.Step++
		SpoonocularConfiguration[update.CallbackQuery.From.ID] = config
		return
	}

	if config.Step == 4 {
		config.Exclude = update.CallbackQuery.Data
		msg.Text = "Thank you! Your custom meal preparing plan is ready. Please, type /my-meal-plans to view it."
		SpoonocularConfiguration[update.CallbackQuery.From.ID] = config
		config.Step++
	}

	fmt.Println(App.Spoonocular.ApiKey)

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

	fmt.Println(uri)

	week, err := CreateMealPreparingPlan(uri)
	if err != nil {
		msg.Text = "Error creating meal preparing plan. Please, try again later."
		App.Bot.Send(msg)
		return
	}

	App.UserMealPlan(update, week)
}

// UserMealPlan executes when user types /my-mealplans
func (App *Application) UserMealPlan(update tgbotapi.Update, week models.Week) {
	if week.Monday.Meals == nil {
		App.Logger.Printf("Error: %s", "UserMealPlan: week.Monday.Meals is nil")
		return
	}

	msg := tgbotapi.NewMessage(update.CallbackQuery.From.ID, "Your meal preparing plan: \n Monday: ")
	for _, meal := range week.Monday.Meals {
		msg.Text += fmt.Sprintf("%s\n", meal.Title)
		msg.Text += fmt.Sprintf("Ready in minutes: %d\n", meal.ReadyInMinutes)
		msg.Text += fmt.Sprintf("Servings: %d\n", meal.Servings)
		msg.Text += fmt.Sprintf("Source url: %s\n\n", meal.SourceURL)
	}
	App.Bot.Send(msg)

	msg.Text = "Tuesday: "
	for _, meal := range week.Tuesday.Meals {
		msg.Text += fmt.Sprintf("%s\n", meal.Title)
		msg.Text += fmt.Sprintf("Ready in minutes: %d\n", meal.ReadyInMinutes)
		msg.Text += fmt.Sprintf("Servings: %d\n", meal.Servings)
		msg.Text += fmt.Sprintf("Source url: %s\n\n", meal.SourceURL)
	}
	App.Bot.Send(msg)

	msg.Text = "Wednesday: "
	for _, meal := range week.Wednesday.Meals {
		msg.Text += fmt.Sprintf("%s\n", meal.Title)
		msg.Text += fmt.Sprintf("Ready in minutes: %d\n", meal.ReadyInMinutes)
		msg.Text += fmt.Sprintf("Servings: %d\n", meal.Servings)
		msg.Text += fmt.Sprintf("Source url: %s\n\n", meal.SourceURL)
	}
	App.Bot.Send(msg)

	msg.Text = "Thursday: "
	for _, meal := range week.Thursday.Meals {
		msg.Text += fmt.Sprintf("%s\n", meal.Title)
		msg.Text += fmt.Sprintf("Ready in minutes: %d\n", meal.ReadyInMinutes)
		msg.Text += fmt.Sprintf("Servings: %d\n", meal.Servings)
		msg.Text += fmt.Sprintf("Source url: %s\n\n", meal.SourceURL)
	}
	App.Bot.Send(msg)

	msg.Text += "Friday: "
	for _, meal := range week.Friday.Meals {
		msg.Text += fmt.Sprintf("%s\n", meal.Title)
		msg.Text += fmt.Sprintf("Ready in minutes: %d\n", meal.ReadyInMinutes)
		msg.Text += fmt.Sprintf("Servings: %d\n", meal.Servings)
		msg.Text += fmt.Sprintf("Source url: %s\n\n", meal.SourceURL)
	}
	App.Bot.Send(msg)

	msg.Text = "Saturday: "
	for _, meal := range week.Saturday.Meals {
		msg.Text += fmt.Sprintf("%s\n", meal.Title)
		msg.Text += fmt.Sprintf("Ready in minutes: %d\n", meal.ReadyInMinutes)
		msg.Text += fmt.Sprintf("Servings: %d\n", meal.Servings)
		msg.Text += fmt.Sprintf("Source url: %s\n\n", meal.SourceURL)
	}
	App.Bot.Send(msg)

	msg.Text = "Sunday: "
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
