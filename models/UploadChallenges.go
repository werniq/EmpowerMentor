package models

var (
	Challenges = []string{
		"Take a cold shower",
		"Run 110% of your normal distance",
		"Study an hour more today",
		"Meditate 10 minutes more than normal",
		"Wake up an hour earlier",
		"Read a book in a day",
		"Cook a new recipe from scratch",
		"Try a new workout routine",
		"Write a handwritten letter to a friend",
		"Volunteer for a local charity",
		"Learn a new language phrase every day for a week",
		"Practice a musical instrument for 30 minutes",
		"Take a day off from social media",
		"Write a gratitude journal for a week",
		"Complete a puzzle or brain teaser",
		"Take a yoga or meditation class",
		"Start a 30-day fitness challenge",
		"Learn to juggle",
		"Write a poem or short story",
		"Try a new hobby or craft",
		"Take a hike in nature",
		"Have a technology-free day",
		"Try a new type of cuisine",
		"Complete a random act of kindness",
		"Take up a new sport",
		"Learn a magic trick",
		"Practice mindfulness for 15 minutes every day",
		"Take a photography walk and capture interesting scenes",
		"Learn to knit or crochet",
		"Explore a new neighborhood in your city",
		"Try a new type of dance class",
		"Do a digital declutter and organize your files",
		"Learn to play chess or a strategy game",
		"Write a letter to your future self",
		"Start a daily journaling habit",
		"Learn to code and build a simple website",
		"Take a cooking or baking class",
		"Visit a museum or art gallery",
		"Learn to solve a Rubik's Cube",
		"Practice deep breathing exercises for relaxation",
		"Start a 30-day decluttering challenge",
		"Try a new form of exercise, such as kickboxing or Pilates",
		"Have a screen-free evening and engage in offline activities",
		"Learn to play a new card game",
		"Volunteer at a local animal shelter",
		"Explore a new genre of music",
		"Take a dance or Zumba class",
		"Learn to make origami creations",
		"Start a daily sketch or doodle practice",
		"Try a new hairstyle or haircut",
		"Go for a bike ride in a new location",
		"Learn to solve Sudoku puzzles",
		"Try a new type of tea or coffee",
		"Plant a garden or grow your own herbs",
		"Write a list of personal goals for the next year",
		"Take a day trip to a nearby town or city",
		"Learn to juggle with three balls",
		"Practice positive affirmations daily",
		"Learn to do a handstand or a yoga pose",
		"Try a new form of art, like pottery or glassblowing",
		"Complete a 30-day meditation challenge",
		"Learn to do a backflip on a trampoline or in a safe environment",
		"Start a daily gratitude practice",
		"Try a new type of exercise class, such as aerial yoga or barre",
		"Learn to play a new board game or card game",
		"Have a technology-free weekend getaway",
		"Learn to solve a crossword puzzle",
		"Start a 30-day writing challenge",
		"Try a new type of outdoor adventure, like rock climbing or kayaking",
		"Learn to make homemade candles or soap",
		"Complete a home improvement project",
		"Explore a new genre of books or movies",
		"Try a new type of art, like watercolor painting or sculpture",
		"Learn to do a handstand push-up or a challenging yoga inversion",
		"Start a daily exercise routine, such as jogging or weightlifting",
		"Take a cooking challenge and make a three-course meal from scratch",
		"Learn to solve a puzzle cube, like a 4x4 Rubik's Cube or a Megaminx",
		"Try a new type of meditation, such as guided visualization or loving-kindness meditation",
		"Learn to make your own jewelry or accessories",
		"Start a daily stretching routine to improve flexibility",
		"Try a new type of martial arts class, like Brazilian Jiu-Jitsu or Krav Maga",
		"Take up a new form of dance, like salsa or hip-hop",
		"Learn to play a new musical instrument",
		"Complete a 30-day photography challenge, capturing different themes or subjects",
		"Try a new type of outdoor activity, like paddleboarding or rock climbing",
		"Learn to solve a challenging puzzle, like a cryptic crossword or a Sudoku variant",
		"Start a daily reading habit, aiming to finish a book each week",
		"Try a new type of DIY project, like woodworking or home decor",
		"Learn to make your own natural skincare or beauty products",
		"Start a daily mindfulness practice, incorporating meditation and mindful activities",
		"Try a new type of dance workout, like Zumba or Bollywood dance",
		"Learn to solve a challenging mathematical problem or puzzle",
		"Complete a 30-day gratitude challenge, expressing gratitude for something every day",
		"Try a new type of water sport, like surfing or paddleboarding",
		"Learn to make your own pottery or ceramics",
		"Start a daily journaling practice, reflecting on your thoughts and experiences",
		"Try a new type of aerial fitness, like aerial silks or aerial hoop",
		"Learn to solve a challenging logic puzzle, like a Sudoku variant or a logic grid",
		"Complete a 30-day self-care challenge, prioritizing self-care activities and practices",
		"Try a new type of cooking technique or cuisine, like sous vide or Thai cuisine",
	}
)

// UploadDifferentChallenges function inserts all prepared challenges to database
func (m *DatabaseModel) UploadDifferentChallenges() error {
	stmt := `
		INSERT INTO 
		    challenges
		    (challenge) 
		VALUES
		    ($1);`

	for i := 0; i <= len(Challenges)-1; i++ {
		_, err := m.DB.Exec(stmt, Challenges[i])
		if err != nil {
			return err
		}
	}
	return nil
}

// UploadChallenge inserts custom challenge to database
func (m *DatabaseModel) UploadChallenge(challenge string) error {
	stmt := `
		INSERT INTO 
		    challenges
		    (challenge) 
		VALUES
		    ($1);`

	_, err := m.DB.Exec(stmt, challenge)
	if err != nil {
		return err
	}

	return nil
}

// ChallengeExists verifies if challenge with this text exists
func (m *DatabaseModel) ChallengeExists(challenge string) (bool, error) {
	stmt := `
		SELECT 
		    id
		FROM 
		    challenges 
		WHERE 
		    challenge = $1;`

	var id int
	err := m.DB.QueryRow(stmt, challenge).Scan(&id)
	if err != nil {
		return false, err
	}

	return id > 0, nil
}
