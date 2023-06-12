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

func LinkedInPost() {
	fmt.Println(`
		Warm welcome, wonderful beings assembled here.
		My name is Oleksandr Matviienko, and I am currently looking for Full-Time Golang Developer opportunity. 
		I would like to introduce a bit myself, and my projects.

		My journey as a developer began few years ago. Being jealous to my classmate, (which was already moving in programming with giant steps)
		I have created a game, which had a logic of simple attraction, somewhere in the spaces of the amusement park. After that, I have been focusing on web development,
		mainly, in Python. I have created News Website, tried to build crypto trading bot, and created a tic tac toe game. 

		One time, when I was driving in my fathers car, he tried to explain me such complex technology as blockchain in simple terms. When we arrived home, I become curious about it.
		Decentralized, immutable, secure, trustable and transparent database? "Why don't we use it everywhere?",  I thought. I started learning Solidity programming language, and 
		slowly became familiar with Ethereum ecosystem concepts.  
		Time passed, and I began looking for work. Being Solidity developer seemed too monotonously for me. I already knew Python but it also seemed something not like that.. 

		So, I decided to learn Golang. Go has pleasing to the eye syntax, concurrency interested my from the most beginning. Huge community support means that you won't stuck with
		your problem for a long time. As well as enormous standard library, which is enough to build anything. 
		My developer journey (for this moment) lasts for nearly 3 years, from which 1.5 i have been learning Golang. 
		I have built various websites, from which I can highlight: 
			- E commerce application:
				From functionality of this website i can hightlight following:
					1. User Management System: Allows administrators to manage users (create, read, update, delete)
					2. Product Management System: Allows administrators to manage products (create, read, update, delete)
					3. JWT Authentication: Provides secure authentication for users and administrators
					4. Pagination: Allows for easy browsing of large numbers of products
					5. Authorization: Requires authentication and authorization for all API calls
					6. Token-Based Authentication: Uses JWT to issue and verify access and refresh tokens
					7. PostgreSQL: Uses PostgreSQL for data storage and management
				Demo of this website, you can view here: https://youtu.be/ZilthxaHNeI

			- GymGo:	
				GymGo is website, for generating workouts for you.	
					I was on my way to Gym, and cannot remember some exercises for legs(one reason is because i do them rarely), but still..
					When I came home, I started building website, that can generate whole workout for you, with provided technique, and even link for proper
					explanation of the technique on YouTube. You can all exercises for particular muscle group (back, chest, biceps etc..), as well as choose
					muscle groups you want to target, enter count of exercises for each, and generate workout. After that, I never skipped leg day :D
				Link to the demonstration: https://youtu.be/i49wkoM_VDE

			- Saurfang Bot:
				As passionate about World of Warcraft game, I decided to create a bot which will express my 'connection' with this wonderful game. I even created an 
				article on Medium for it (https://medium.com/@qniwwwersss/how-to-code-discord-moderator-bot-in-golang-world-of-warcraft-style-e35c326da145)
				Nothing related to WoW I actually have not provided.. But greeting end with 'For the Horde'!!! And Saurfang image on avatar.. Like, is it enough????

			- AuthDigitBot:
				This discord bot is used for my custom authentication system in minecraft server. Initially starting with the network library, the MC TCP-based 
				network communication protocol was implemented with pure Go. A package is provided to facilitate the writing of command-line MC robot programs.

				Instead of relying on traditional password-based login methods, the bot generates unique codes consisting of digits only. 
				Users can obtain these codes through the Discord platform, where they are securely generated and distributed. To authenticate on the 
				Minecraft server, users simply need to enter the code instead of typing their password. This approach enhances security by eliminating 
				the need to remember or share passwords,

				This bot it best choice for Minecraft server, because it creates connection with AuthMe plugin database, and sends message to user, whenever any another IP address
				is logged into this account. As well as notifying about attempt to change password and gaining experience and server currency using custom quiz system.
			
			- S3 Compatible HTTP Service:
				An S3-compatible HTTP service refers to a service that implements an API that is compatible with Amazon S3 (Simple Storage Service). 
				Amazon S3 is a popular cloud storage service provided by Amazon Web Services (AWS) that allows users to store and retrieve large amounts of data.
				I decided to create my own one, and really enjoyed in process! I have implemented functionality to interact with that API using front-end side, and it works fine.
				
				It provides functionality such as:
					1. Create bucket
					2. Upload files to bucket
					3. GetFiles from bucket
					4. GetFile from bucket (using some key
					5. CheckBucketName for proper verification of creating bucket
				You can upload files from front-end as well.

			- EmpowerMentor:
					This is my current project. I am working on bot, which will focus on self-improvement. 
					It notifies you about drinking water, exercising, meditating and reading once a day (Drink water few times), and sends motivational quotes once a day (with ability to create your own)
					It may be used to create personal workouts, supplement intake plan, generating meal prepare plan for day (or whole week), based on maximum calories count,
					excluding products, which you have allergy on. It sends you recommendation for your mental health, for better sleeping, and physical health. As well, through
					day it sends you challenges. For example: run 110% of your normal distance, take a ice bath, or cold shower. It has two type of user: ordinary user and admin.
					AdminsId are stored in database, and admins may view issues related to bot functionality. Admins may upload new challenges, new motivational quotes, manage users,
					subscriptions. On user's configuration of bot, it requires you to provide wake up and bed time. At that time it will send you appropriate message. You can change it,
					in addition to preferable time for meditation, reading and exercise (which you input in setting up the bot). Hope, it will help someone to become more disciplined, and more stoic :D
		`)
}
