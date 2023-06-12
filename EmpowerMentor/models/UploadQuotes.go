package models

var (
	quotes = []string{
		"Believe in yourself and all that you are. Know that there is something inside you that is greater than any obstacle. - Christian D. Larson",
		"The future belongs to those who believe in the beauty of their dreams.- Eleanor Roosevelt",
		"Don't watch the clock; do what it does. Keep going. - Sam Levenson",
		"Success is not final, failure is not fatal: It is the courage to continue that counts. - Winston S. Churchill",
		"Your time is limited, don't waste it living someone else's life. - Steve Jobs",
		"The harder you work for something, the greater you'll feel when you achieve it. - Unknown",
		"Challenges are what make life interesting and overcoming them is what makes life meaningful.  - Joshua J. Marine",
		"Believe you can and you're halfway there. - Theodore Roosevelt",
		"Dream big and dare to fail. - Norman Vaughan",
		"Your life does not get better by chance, it gets better by change. - Jim Rohn",
		"No matter how hard the battle gets, or no matter how many people don't believe in your dream, never give up!  - Eric Thomas ",
		"The best way to predict the future is to create it. - Peter Drucker",
		"Difficulties in life are intended to make us better, not bitter. - Dan Reeves ",
		"You have to fight through some bad days to earn the best days of your life.",
		"Opportunities don't happen. You create them. - Chris Grosser",
	}
)

// TruncateMotivationalQuotes truncates motivational quotes table
func (m *DatabaseModel) TruncateMotivationalQuotes() error {
	stmt := "TRUNCATE motivational_quotes;"

	_, err := m.DB.Exec(stmt)
	if err != nil {
		return err
	}

	return nil
}

// UploadMotivationalQuotes uploads motivational quotes to the database
func (m *DatabaseModel) UploadMotivationalQuotes() error {
	stmt := "INSERT INTO motivational_quotes (quote) VALUES ($1);"

	for i := 0; i < len(quotes); i++ {
		_, err := m.DB.Exec(stmt, quotes[i])
		if err != nil {
			return err
		}

	}

	return nil
}

// UploadOneMotivationalQuote function is used by admins to upload new motivational quotes
func (m *DatabaseModel) UploadOneMotivationalQuote(quote string) error {
	stmt := `
			INSERT INTO 
			    motivational_quotes 
			    (quote) 
			VALUES 
			    ($1);
			`

	_, err := m.DB.Exec(stmt, quote)
	if err != nil {
		return err
	}

	return nil
}

// GetRandomQuote returns a random motivational quote from the database
func (m *DatabaseModel) GetRandomQuote() (string, error) {
	stmt := `
			SELECT 
			    quote
			FROM 
			    motivational_quotes
			ORDER BY 
			    random() 
			        * 
			    (SELECT COUNT(*) 
			FROM 
			    motivational_quotes) + 1
			LIMIT 
			    1;`

	var quote string

	err := m.DB.QueryRow(stmt).Scan(&quote)

	if err != nil {
		return "", err
	}

	return quote, nil
}
