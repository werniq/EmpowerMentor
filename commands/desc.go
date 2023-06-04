package commands

import "fmt"

func desc() {
	fmt.Println(`
			Bot which specified on self-improvement.
			What it should do? 

			Firstly, send daily reminders for following actions: drinking water, doing exercises, reading books, sleeping.

			Plus, if user wants, it should send motivational quotes, and recommendations for books, exercises, etc.
			Secondly, there should be an ability to create custom plans for user, and send reminders for them.
			And as well as automatic recommendations sending, user may ask bot for recommendations, and bot should send them.

			User may ask bot to create meal preparation plan, supplements intake plan, exercises plan, etc.
			As well as custom plans, user may create his own meal preparation plan, supplements intake plan.
		
			That is plan for acquiring habits. 

			Summary:
				1. Daily reminders for following actions: drinking water, doing exercises, reading books, sleeping.
				2. Ability to create custom plans for user, and send reminders for them.
				3. Ability to ask bot for recommendations, and bot should send them.
				4. Ability to create custom plans for user, and send reminders for them.
				5. Ability to create custom meal preparation plan, supplements intake plan, exercises plan, etc.
				6. Create own morning routine, and follow it
				7. Acquiring stoic philosophy habits
				8. Everyday challenges
				9. News about topics which user interested in, and which he wants to follow
			
		
				Notes:
					Create workout with keyboard 

	`)
}

func process() {
	fmt.Println(`
				1. Bot starts 
				2. Bot sets up database connection, and after that creates few necessary database tables
				3. Bot starts listening for updates
				4. Bot receives update, and depending on message text, it calls corresponding function
				5. Bot sends appropriate response to user 
			`)
}

func functions() {
	fmt.Printf(`
		/add-quote (only admins) - add motivational quote
		/add-quote-for-review - for users, create request for admins to review quote and if verified, add it to database
		/report <!-{isuse}-> - reporting issues
		`)
}
