package models

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type Workout struct {
	ID              string        `json:"id"`
	Title           string        `json:"title"`
	Gender          string        `json:"gender"`
	MotivationQuote string        `json:"motivation_quote"`
	Exercises       []Exercise    `json:"exercises"`
	Duration        time.Duration `json:"duration"`
	CreatedAt       time.Time     `json:"created_at"`
}

type Exercise struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Technique string `json:"technique"`
	VideoURI  string `json:"videoURI"`
}

type CardioWorkout struct {
	ID       int           `json:"id"`
	Title    string        `json:"title"`
	Gender   string        `json:""`
	Duration time.Duration `json:"duration"`
}

var (
	MuscleToCount = map[string]int{
		"back":    15,
		"chest":   16,
		"biceps":  15,
		"legs":    15,
		"glutes":  18,
		"triceps": 15,
	}
	muscles = []string{"back", "chest", "biceps", "triceps", "shoulders", "legs", "glutes"}
)

// back -> 16
// chest -> 17
// biceps -> 15
// legs -> 15
// glutes -> 13
// in progress: biceps & shoulders

func (m *DatabaseModel) CreateMuscleGroupTables() error {
	var stmt string
	var err error
	for i := 0; i <= len(muscles)-1; i++ {
		stmt = fmt.Sprintf(`CREATE TABLE "%s" (
    		id bigserial not null primary key,
    		title varchar(255) not null,
    		technique text not null,
    		videoURI text not null
		);`, muscles[i])

		_, err = m.DB.Exec(stmt)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *DatabaseModel) DropTables() error {
	var stmt string
	var err error
	for i := 0; i <= len(muscles)-1; i++ {
		stmt = fmt.Sprintf(`DROP TABLE "%s";`, muscles[i])

		_, err = m.DB.Exec(stmt)
		if err != nil {
			return err
		}
	}

	return nil
}

// TruncateTable removes all records from table
func (m *DatabaseModel) TruncateTable(tableName string) error {
	stmt := `TRUNCATE TABLE $1`

	row := m.DB.QueryRow(stmt, tableName)
	if row.Err() != nil {
		return row.Err()
	}

	return nil
}

// GetExerciseById retrieves exercise by given id from table name, granted in arguments
func (m *DatabaseModel) GetExerciseById(id int, table string) (Exercise, error) {
	var exercise Exercise
	var title string
	var technique string
	var videoURI string

	stmt := fmt.Sprintf(`SELECT title, technique, videoURI FROM "%s" WHERE id = $1`, table)

	row, err := m.DB.Query(stmt, id)
	if err != nil {
		fmt.Printf("Error getting exercise from TABLE [%s] by ID: %v", table, err)
		return exercise, err
	}
	if row.Next() {
		if err := row.Scan(&title, &technique, &videoURI); err != nil {
			fmt.Printf("Error retrieveing data from table: %v", err)
			return exercise, err
		}
	}

	exercise = Exercise{
		Title:     title,
		Technique: technique,
		VideoURI:  videoURI,
	}

	return exercise, nil
}

// TruncateTables removes all records from all existed tables, in my database
func (m *DatabaseModel) TruncateTables() error {
	stmt := `TRUNCATE TABLE $1`
	for i := 0; i <= len(muscles)-1; i++ {
		if row := m.DB.QueryRow(stmt, muscles[i]); row.Err() != nil {
			fmt.Printf("Error truncating table %s: %v", muscles[i], row.Err())
			return row.Err()
		}
	}
	return nil
}

// SaveAllLegsExercises stores all legs exercises in database
func (m *DatabaseModel) SaveAllLegsExercises() error {
	var legs []string
	var technique []string
	var videoURI []string

	legs = append(legs, "Back Squat")
	technique = append(technique, `
		Set up the barbell: Start by positioning the barbell across your upper traps and then secure it in place with both hands. Make sure the bar is balanced and evenly distributed.
		1. Stance: Stand with your feet slightly wider than shoulder-width apart, with your toes pointing slightly outwards.
		2. Descent: Lower yourself by bending at the hips and knees, while keeping your back straight and chest up. As you descend, keep your weight on your heels, and your knees tracking over your toes.
		3. Bottom position: Once your thighs are parallel to the ground, pause for a moment before you begin the ascent. Keep your core tight and make sure that your lower back does not round.
		4. Ascent: Drive through your heels to stand back up, keeping your torso upright and maintaining a stable bar position.
		5. Repeat: Repeat the movement for the desired number of repetitions.
		It's important to keep your technique and form in check throughout the entire movement, as improper form can lead to injury. Start with lighter weights until you are comfortable with the movement, and then gradually increase the weight.`)
	videoURI = append(videoURI, "https://youtu.be/Uv_DKDl7EjA")

	legs = append(legs, "Front Squat")
	technique = append(technique, `
		 1. Set the rack to the appropriate level for your height
		 2. Place the barbell on the front deltoid
		 3. Select a grip that feels comfortable
		 4. Set your elbow position
		 5. Breath and brace your core before take-off
		 6. Walk back from the rack with the minimal distance possible
		 7. Set your squat stance width
		 8. Breathe and brace your core again before squatting down
		 9. Crack at your hips and knees at the same time to initiate movement
		 10. Use a tempo that allows you to maintain tightness and control
		 11. Ensure your knees are tracking over your toes
		 12. Maintain an upright torso
		 13. Go as deep as your mobility allows
		 14. Push your knees forward in the bottom of the squat
		 15. Drive your feet through the floor and use your quads to drive up
		 16. Drive your elbows up and forward to prevent bar slippage
		 17. Accelerate through the entire range of motion to standing
	`)
	videoURI = append(videoURI, "https://youtu.be/7pyxT5hqmQY")

	legs = append(legs, "Bulgarian Split Squat")
	technique = append(technique, `
		1. Stand 4-5 inches in front of a bench with your feet shoulder-width apart from each other. The body should be facing forwards away from the bench.
		2. Lift the barbell onto your chest with a pronated grip. Lift the barbell overhead and then rest it on your shoulders.
		3. One foot should be moved backwards so that it rests on the bench, whilst the other foot is positioned in front.
		4. Gradually lower the leg until it is low enough to feel a contraction. Ensure that the knee doesn’t hit the floor and the knee should be over your toes.
		5. As you exhale your breath, lift your leg back to the starting position.
		6. Once performed, return to the starting position and perform with the opposing leg.
	`)
	videoURI = append(videoURI, "")

	legs = append(legs, "Leg Press")
	technique = append(technique, `Load your desired amount of weight onto the machine and sit down on the leg press seat. Place your legs on the pressing platform with a shoulder-width stance in front of your torso in the centre of your body. Disengage the safety levers on either side and take hold of the handles to your sides.`)
	videoURI = append(videoURI, "https://www.youtube.com/watch?v=s9-zeWzPUmA&embeds_euri=https%3A%2F%2Fwww.bing.com%2F&embeds_origin=https%3A%2F%2Fwww.bing.com&feature=emb_logo")

	legs = append(legs, "Hack Squat")
	technique = append(technique, `
		1. Retract your scapular and use the shoulder pads to cushion your shoulders.
		2. Keep your head upright, take hold of the safety bars and unlock.
		3. Straighten your legs (without locking them) and stand with your feet at shoulder width apart.
		4. Slightly face your toes outwards.
		5. Inhale when you perform the eccentric part of the movement (descent) and exhale when exploding on the concentric.
		6. Maintain a tight posture, descend until you break parallel and explode back up with control (without bouncing up from the bottom of the movement).
		7. Drive through the ball of your heels.
	`)
	videoURI = append(videoURI, "https://youtu.be/bhfyY8F8F24")

	legs = append(legs, "Romanian Deadlift")
	technique = append(technique, `
		Stand with feet hip-width apart and hold a barbell in front of your thighs. Slight bend in the knees, roll shoulders back, and keep them back throughout the exercise. Inhale and press hips back, hinging forward from the hips, keeping the barbell close. Stop when you feel a stretch in your hamstrings. Exhale, use hamstrings and glutes to pull back to standing. Complete set and replace barbell carefully on rack
	`)
	videoURI = append(videoURI, "https://youtu.be/_oyxCn2iSjU")

	legs = append(legs, "Nordic Hamstring Curl")
	technique = append(technique, `
		Start on knees with pad and have partner hold legs or anchor under equipment. Feet and ankles in line with knees, shoulders over hips, chin tucked. Arms at sides, pre-tension shoulders and hips, engage core, squeeze glutes and hamstrings. Slowly lower self while maintaining straight line from knees to head, using legs, then hands to catch self. Squeeze hamstrings to pull back to start, using hands if needed. Shoulders over hips at end of each repetition.
	`)
	videoURI = append(videoURI, "")

	legs = append(legs, "Landmine Goblet Squat")
	technique = append(technique, ``)
	videoURI = append(videoURI, "")

	legs = append(legs, "Reverse Lunge")
	technique = append(technique, ``)
	videoURI = append(videoURI, "https://youtu.be/kFSnvwvc5ac")

	legs = append(legs, "Barbell Hip Thrust")
	technique = append(technique, `
		1. Sit on the ground with a bench or box behind you.
		2. Roll a barbell up your legs until it rests on your hips.
		3. Lean back against the bench and bend your knees so that your feet are flat on the floor.
		4. Drive through your heels to lift your hips until they are in line with your shoulders and knees.
		5. Lower your hips back to the starting position and repeat.
		6. Keep your head, neck, and spine neutral throughout the exercise and avoid rounding your lower back.
		7. Maintain control and focus on using your glutes to lift and lower your hips.
		8. When finished, carefully roll the barbell off your hips and return it to the starting position.
	`)
	videoURI = append(videoURI, "https://youtu.be/xDmFkJxPzeM")

	legs = append(legs, "Leg Extension")
	technique = append(technique, `This technique is pretty simple: you lift the weight using two limbs and you lower the weight with one limb. For the leg extension, you'll lift the weight up pretty fast with both legs, lower slowly with one, then repeat, alternating legs. Accentuating the eccentric (negative) stress will lead to more strength gains.`)
	videoURI = append(videoURI, "https://youtu.be/9nmAtebIwy8")

	legs = append(legs, "Seated Leg Curl")
	technique = append(technique, ``)
	videoURI = append(videoURI, "")

	legs = append(legs, "Lying Leg Curl")
	technique = append(technique, `
		1. Lie facedown on a leg curl machine with your ankles under the pads.
		2. Adjust the pads so they are in line with your ankle joints.
		3. Grasp the handles on the side of the machine for stability.
		4. Without moving your upper body, curl your legs up towards your glutes.
		5. Pause briefly at the top, then lower your legs back to the starting position.
		6. Repeat for desired reps.
		7. Maintain control and focus on using your hamstrings to lift and lower your legs.
	`)
	videoURI = append(videoURI, "https://youtu.be/q1cKTmaeQWo")

	legs = append(legs, "Standing Calf Raise")
	technique = append(technique, `
		1. Stand with your feet hip-width apart on a calf raise machine or a raised platform.
		2. Place the balls of your feet on the edge of the platform, with your heels hanging off.
		3. Grasp the handles on the sides of the machine for stability, or use a barbell or dumbbells for added resistance.
		4. Raise your heels as high as you can, while keeping your legs straight.
		5. Pause briefly at the top, then lower your heels back to the starting position.
		6. Repeat for desired reps.
		7. Maintain control and focus on using your calf muscles to lift and lower your heels.
	`)
	videoURI = append(videoURI, "https://youtu.be/K_jsGgztcGU")

	legs = append(legs, "Prowler Push")
	technique = append(technique, `
		1. Stand at the back of a Prowler sled facing forward, with your feet shoulder-width apart.
		2. Grasp the handles of the sled, with your arms straight and your shoulders pulled back.
		3. Lower your body into a half-squat position, keeping your back straight.
		4. Drive your legs to push the sled forward, using your legs and hips.
		5. Keep your arms straight and your shoulders pulled back throughout the movement.
		6. Push the sled for a set distance or for a set time, then walk back to the starting position.
		7. Repeat for desired reps.
		8. Maintain control and focus on using your legs and hips to drive the sled forward, without straining your lower back.
	`)
	videoURI = append(videoURI, "https://youtu.be/9XRRXaUpnLk")

	legs = append(legs, "Assault Bike")
	technique = append(technique, `
		1. Sit on the seat of the Assault Bike, with your feet securely on the pedals.
		2. Grasp the handlebars, with your arms slightly bent and your shoulders pulled back.
		3. Start pedaling, using both your legs and arms to generate power.
		4. Keep a steady pace, using your legs to drive the pedals in a circular motion.
		5. Maintain control and focus on using your legs and arms to generate power, without straining your lower back.
		6. Adjust the resistance as needed to increase or decrease the difficulty of the workout.
		7. Pedal for a set time or distance, then stop and rest as needed.
		8. Repeat for desired reps.
	`)
	videoURI = append(videoURI, "https://youtu.be/RPY7HTGfOiU")

	stmt := `
		insert into 
			legs (title, technique, videoURI)
		values ($1, $2, $3)
	`

	for i := 0; i < len(legs)-1; i++ {
		if row := m.DB.QueryRowContext(context.Background(), stmt, legs[i], technique[i], videoURI[i]); row.Err() != nil {
			fmt.Printf("Error inserting values in database: %v", row.Err())
			return row.Err()
		}
	}

	fmt.Println("All leg exercises inserted into database.")
	return nil
}

// StoreAllGlutesExercises stores all exercises for glutes in database
func (m *DatabaseModel) StoreAllGlutesExercises() error {
	var glutes []string
	var technique []string
	var videoURI []string
	glutes = append(glutes, "Reverse Lunges")
	technique = append(technique, "Stand tall with your feet hip-width apart. Take a step backward with one leg, lowering your body until your front knee is bent at a 90-degree angle. Push through your front heel to return to the starting position. Repeat on the other leg.")
	videoURI = append(videoURI, "https://www.youtube.com/watch?v=QD6S8XIfOEM&ab_channel=TheBodyCoachTV")

	glutes = append(glutes, "Conventional Deadlift")
	technique = append(technique, ``)
	videoURI = append(videoURI, "https://youtu.be/kH-GjbMfi1o")

	glutes = append(glutes, "Belt Squat")
	technique = append(technique, `Set the belt around your hip, brace your core, and lift the weight. Release the weight pin and get yourself into position. With your chest up, squat down until the bottoms of your thighs are parallel to the floor. Now, drive back up by pushing your feet through the platform.`)
	videoURI = append(videoURI, "https://www.youtube.com/watch?v=oAK7QmugOzU&t=2s&ab_channel=MarkBell-SuperTrainingGym")

	glutes = append(glutes, "Banded Barbell Romanian Deadlift")
	technique = append(technique, `Load a barbell up with less weight than you would for a traditional deadlift, but get in the same deadlift position — feet shoulder-width apart and hands gripping the bar just outside the knees. Step your weight into the band and assume a slight forward lean. Push your hips back and lower the bar until it’s in the middle of your shins. You should feel a stretch in your hamstrings and glutes. Finish by driving your hips forward, bringing the weight back up, and driving into the band to the starting position.`)
	videoURI = append(videoURI, "https://www.youtube.com/watch?v=bSYqD-PJGS8&ab_channel=PhysiqueDevelopment")

	glutes = append(glutes, "Rear Foot Elevated Split Squat")
	technique = append(technique, `Stand with a dumbbell in each hand, and take a step forward, placing one foot, toes down, on an elevated surface behind you. It should be low enough to not shift your hips to one side. The distance of your step should be about one foot. Keep your chest up and squat down until both of your legs bend to around 90 degrees. Stand back up by driving your front foot through the floor.`)
	videoURI = append(videoURI, "https://www.youtube.com/watch?v=dWHkwgY9I08")

	glutes = append(glutes, "Sumo Deadlift")
	technique = append(technique, `To get into position, set your feet just outside of hip-width, with your toes pointed out around 30 degrees. From there, push your hips back as you reach down for the bar. Allows the knees to bend naturally as you reach down for the bar. Secure your grip, engage your abdominals, and maintain a neutral head position. Engage the upper back muscles along with the abdominals to help protect the spine and support the torso; then, drive through the floor. As you drive through the floor, drive your hips forward as you reach the top. Return safely back to the starting position by hinging at the hip and controlling the bar as it returns back to the floor.`)
	videoURI = append(videoURI, "https://www.youtube.com/watch?v=_vTsPOC6ZFc")

	glutes = append(glutes, "Modified Curtsy Lunge")
	technique = append(technique, `Stand with a dumbbell in each hand, and place one foot, toes down, on an elevated surface behind you. It should be low enough to not shift your hips to one side. The lead foot rotated internally around 10-20 degrees, aligning with the toes of the lead foot with the knee of the back leg. With a slight lean forward in your torso, squat down until both of your legs bend to around 90 degrees. Stand back up by driving your front foot through the floor.`)
	videoURI = append(videoURI, "https://www.youtube.com/watch?v=bQV7qr8tt4M")

	glutes = append(glutes, "Walking Lunge")
	technique = append(technique, `Stand with your feet together, and then take a step forward roughly 18 to 24 inches and plant your foot firmly to the ground. From there, you will allow your front knee to track forward — aiming between the first and second toe — while your back knee drops straight down to the ground. Then, while driving through the floor with your front foot, move your body forward to a standing position, where your back foot will meet the position of the front one.`)
	videoURI = append(videoURI, "https://www.youtube.com/watch?v=YYWhkctnP2o&embeds_euri=https%3A%2F%2Fwww.bing.com%2F&embeds_origin=https%3A%2F%2Fwww.bing.com&feature=emb_logo")

	glutes = append(glutes, "Cable Glute Kickback")
	technique = append(technique, `Place a strap attached to the cable around the ankle. Keep your back in a neutral position with your abs engaged. The body will be positioned off-center with the working leg in line with the cable attachment. Tilt your body forward and kick your leg out behind you while maintaining a very slight bend in the knee. Move your leg by squeezing the glute, not arching the lower back.`)
	videoURI = append(videoURI, "https://www.youtube.com/watch?v=dJa_Nf4zdik&t=31s")

	glutes = append(glutes, "Step Down")
	technique = append(technique, `Start with one foot close to the edge of a step up box or bench — ensuring the whole foot is in contact with the surface — with the other foot hanging off. Drop the foot to the ground, controlling your body weight with the opposite leg. Tap the heel of the foot to the ground and drive through the step with the working leg to return back to the starting position.`)
	videoURI = append(videoURI, "https://www.youtube.com/watch?v=ee3lC7iLfss&embeds_euri=https%3A%2F%2Fwww.bing.com%2F&embeds_origin=https%3A%2F%2Fwww.bing.com&feature=emb_logo")

	glutes = append(glutes, "Smith Machine Reverse Lunge")
	technique = append(technique, `Stand with your feet together on an elevated surface, and position your body under the bar in the Smith machine. Unrack the weight and take a step back with one leg until it’s behind you and your knee is an inch or so above the floor. Your front leg should bend at a 90-degree angle as well. Drive through your front foot and stand back up with control.`)
	videoURI = append(videoURI, "https://www.youtube.com/watch?v=JzrwCj4dVc0&t=2s")

	glutes = append(glutes, "Lateral Lunge")
	technique = append(technique, `Take a step out to the side roughly 18 to 24 inches and plant your foot firmly to the ground. From there, allow your front knee to track forward while your body tracks outward to one side; then, while driving through the floor with your lead foot, move your body back to the starting position.`)
	videoURI = append(videoURI, "https://www.youtube.com/watch?v=XjBb6sR4Bz8&ab_channel=PhysiqueDevelopment")

	glutes = append(glutes, "Goblet Squat")
	technique = append(technique, `Stand with your feet around shoulder width apart with a dumbbell or kettlebell and hold the weight directly under your chin with your elbows tucked in. Brace your core, tense your back, and ensure that you feel stable. Keep your chest up and squat down until both of your legs bend to around 90 degrees. Stand back up by driving through the floor.`)
	videoURI = append(videoURI, "https://www.youtube.com/watch?v=CkFzgR55gho&ab_channel=PhysiqueDevelopment")

	glutes = append(glutes, "Banded Glute Bridge")
	technique = append(technique, "Place a resistance band just above your knees and lie on your back with your knees bent and feet flat on the ground. Push through your heels and lift your hips off the ground, squeezing your glutes at the top. Lower your hips back down and repeat.")
	videoURI = append(videoURI, "https://www.youtube.com/watch?v=7a7jZj6YO80&ab_channel=AnabolicAliens")

	glutes = append(glutes, "Dumbbell Glute Bridge")
	technique = append(technique, "Lie on your back with your knees bent and feet flat on the ground. Hold a dumbbell on your hips and push through your heels to lift your hips off the ground, squeezing your glutes at the top. Lower your hips back down and repeat.")
	videoURI = append(videoURI, "https://www.youtube.com/watch?v=8bbE64NuDTU&ab_channel=StrongLifts")

	glutes = append(glutes, "Glute Bridge March")
	technique = append(technique, "Lie on your back with your knees bent and feet flat on the ground. Lift your hips off the ground into a glute bridge position. While maintaining the bridge, lift one foot off the ground and bring your knee towards your chest. Lower your foot back down and repeat on the other side.")
	videoURI = append(videoURI, "https://www.youtube.com/watch?v=7nzXdrjgeHg&ab_channel=HeatherRobertson")

	glutes = append(glutes, "Gluteral bridge")
	technique = append(technique, "Lie on your back with your knees bent and feet flat on the ground. Lift your hips off the ground into a glute bridge position. While maintaining the bridge, lift one foot off the ground and bring your knee towards your chest. Lower your foot back down and repeat on the other side.")
	videoURI = append(videoURI, "https://youtu.be/S_uZP4UH6J0")

	glutes = append(glutes, "Hip Thrust(specifically for my lovely girlfriend <3 )")
	technique = append(technique, `Things to keep in mind: 1. Make Sure the Bench Hits Your Mid Back Find a bench that comes to your middle back when sitting down on the floor. If you're using a tall bench, sit on a pad or some mats to raise your body. If your bench is too high on your back, you won't be able to get as much leverage through your hips, says Jason Pak, CPT, a certified personal trainer, USA Weightlifting–certified sports performance coach and co-owner of Achieve Fitness Boston. 2. Play With Foot Placement. 3. Look forward. Instead of letting your head fall back in line with the bench, keep your chin slightly lowered toward your chest. Throughout the entire exercise, keep your eyes locked on the wall in front of you and your head will naturally stay in place. Technique:	1. Sit on the floor with your mid back against the edge of a bench or box and rest a barbell on top of your hips, holding it with both hands.	2. Squeeze your glutes and push through your heels to raise your hips and the weight up toward the ceiling.	3. Pause for a few seconds, then slowly lower your hips down to return to the starting position.`)
	videoURI = append(videoURI, "https://youtu.be/S_uZP4UH6J0")

	stmt := `
		insert into 
			glutes(title, technique, videouri)
		values  ($1, $2, $3)
	`

	for i := 0; i <= len(glutes)-1; i++ {
		if row := m.DB.QueryRowContext(context.Background(), stmt, glutes[i], technique[i], videoURI[i]); row.Err() != nil {
			fmt.Printf("Error inserting values in database: %v", row.Err())
			return row.Err()
		}
	}

	fmt.Println("All GLUTES exercises inserted into database.")
	return nil
}

// StoreAllBicepsExercises stores all exercises for glutes in database
func (m *DatabaseModel) StoreAllBicepsExercises() error {
	var biceps []string
	var technique []string
	var videoURI []string

	biceps = append(biceps, "Barbell Curl")
	technique = append(technique, `
		Grab a barbell with an underhand grip, slightly wider than the shoulders. With the chest up and shoulder blades pulled tightly together, expose the front of your biceps by pulling the shoulders back into the socket. The elbows should reside under the shoulder joint, or slightly in front by the ribs. Curl the barbell up using the biceps, making sure not to let the torso lean forward, shoulder collapse forward, or the elbows slide backward to the side of the body (they should stay slightly in front of the shoulders).
	`)
	videoURI = append(videoURI, "https://youtu.be/kwG2ipFRgfo")

	biceps = append(biceps, "Chin-Up")
	technique = append(technique, "Hang from a bar with palms facing you and the hands about shoulder-width apart, or slightly wider. From a dead hang, squeeze your shoulder blades together and pull your body up, making sure not to let the body fold inwards (so many people do this) until your chin is at or above the bar.")
	videoURI = append(videoURI, "https://www.youtube.com/watch?v=oqDpaZkfV0o&ab_channel=BarBend")

	biceps = append(biceps, "EZ-Bar Preacher Curl")
	technique = append(technique, `
		Sit down on a preacher bench and rest the back of your triceps on the pad. Set your body in the same position as the standard barbell biceps curl (chest up, shoulders back, and elbows slightly forward). Grasp the EZ-bar handle on the inner angled pieces. This will place your hands slightly narrower than shoulder-width and on a semi-supinated angle. With the body locked in place, curl the bar upwards as you flex the biceps, briefly pausing at the top of the curl to flex the biceps. Lower the weight under control.
	`)
	videoURI = append(videoURI, "https://www.youtube.com/watch?v=I7PoPjFF2Ww&ab_channel=NikkieyStott")

	biceps = append(biceps, "Hammer Curl")
	technique = append(technique, `Hold a dumbbell in each hand while standing. Turn your wrists so that they’re facing each other. Keep your arms tucked in at your sides and flex your elbows to curl the dumbbells up towards your shoulders. Lower them back down with control.
	`)
	videoURI = append(videoURI, "https://youtu.be/TwD-YGVP4Bk")

	biceps = append(biceps, "Incline Dumbbell Curl")
	technique = append(technique, `Lay back on an incline bench, angled at about 60 degrees, with a dumbbell in each hand. Let your arms hang so they’re fully extended. Without moving your shoulders, curl the weight up to your shoulders. Hold the top of the movement for about a second, and then slowly lower the dumbbells with control.
	`)
	videoURI = append(videoURI, "https://youtu.be/soxrZlIl35U")

	biceps = append(biceps, "Facing-Away Cable Curl")
	technique = append(technique, `Set the handles of the cable pulleys to the lowest setting and attach D-handles to each pulley. Pick up a handle in each hand. Tense your upper back and let your arms hang so they’re fully extended. Without moving your shoulders, curl the weight up toward your shoulders. Hold the top of the movement for about a second, and then slowly lower the handles with control.
	`)
	videoURI = append(videoURI, "https://www.youtube.com/watch?v=fV9BpknCjGM&ab_channel=PhysiqueDevelopment")

	biceps = append(biceps, "Cable Curl")
	technique = append(technique, `Attach the desired handle to the pulley of a cable machine set to the lowest height. Grab the handle in both hands and take a few steps back so there’s constant tension on the cable (the weight stack should be elevated the entire time). Curl the bar up to your chest and then slowly lower it back down.
	`)
	videoURI = append(videoURI, "https://www.youtube.com/watch?v=NFzTWp2qpiE&ab_channel=FitFatherProject-FitnessForBusyFathers")

	biceps = append(biceps, "Concentration Curl")
	technique = append(technique, `Sit on a bench with your feet set wide enough to allow your arm to hang in the middle, with your elbow resting on the inside of the thigh. With a dumbbell in hand, slowly curl the dumbbell upward at a controlled tempo, concentrating on contracting the biceps to move the load. At the top of the movement, flex as hard as possible, then slowly lower the load. The key is not to lose tension on the biceps at any point in the range of motion.
	`)
	videoURI = append(videoURI, "https://www.youtube.com/watch?v=0AUGkch3tzc")

	biceps = append(biceps, "Cable Concentration Curl ")
	technique = append(technique, `Stand in front of a single cable on a functional trainer or cable tower. With the cable set around chest height, grab the handle with a supinated grip (palm facing up), and slightly lean your torso forward. Your working arm should be angled across the body while you curl the handle toward the opposite ear. Keep tension on the biceps all the way to the top of the movement, then slowly lower the load. 
	`)
	videoURI = append(videoURI, "https://youtu.be/ry2LiWeKtoo")

	biceps = append(biceps, "High Cable Curl")
	technique = append(technique, `Set a cable pulley to about shoulder height and attach D-handles to each cable pulley. Grab the bar with a supinated grip (palm facing up). Keep tension on the biceps all the way to the top of the movement, then slowly lower the load back to the starting position. The key to this exercise is maintaining your shoulder position throughout the range of motion, not allowing your elbows to dip — making it easier. Maintaining tension in the upper back will help keep shoulders stable and arm position constant, driving up tension in the biceps.
	`)
	videoURI = append(videoURI, "https://www.youtube.com/watch?v=Okl64XegSu0&ab_channel=PhysiqueDevelopment")

	biceps = append(biceps, "Cable Rope Supinating Curl")
	technique = append(technique, `Stand in front of a single cable on a functional trainer or cable tower. Set the cable to a lower setting with each side of the rope attachment in your hands (palms facing each other). To complete the curling motion, start by flexing your elbows (bringing the hands up toward the shoulders), and around one-third of the way up, rotate (supinate) your hands to face up as you continue to curl up. The key to increasing the challenge of supination of this exercise is waiting to supinate until you’re a third the way into the rep — this ensures the muscles responsible for supination can kick in while the resistance is on them.
	`)
	videoURI = append(videoURI, "https://www.youtube.com/watch?v=wW-lmxnc2vs")

	biceps = append(biceps, "Cable Hammer Curl")
	technique = append(technique, `With the cables positioned lower on the cable tower, the lifter will grip each handle with a neutral grip (palms facing each other), take a step back, engage the upper back to add stability to the upper body and curl the weight up. Maintain a neutral grip (palms facing each other). Squeeze and contract once you reach the top of the movement, then slowly lower the load back to the starting position.
	`)
	videoURI = append(videoURI, "https://www.youtube.com/watch?v=Xm05VEYT09Y&ab_channel=PhysiqueDevelopment")

	biceps = append(biceps, "Dual Cable Preacher Curl")
	technique = append(technique, `Set up a preacher bench roughly three to five feet away from a cable tower with two cable pulleys. Set the pulleys, so they’re slightly lower than the bench. Sit on the preacher bench and have a training partner hand you both handles. Position your elbows so that they rest over the pad. Lower the arms until your elbows are nearly locked out, and then curl the weight back up. 
	`)
	videoURI = append(videoURI, "https://www.youtube.com/watch?v=28frFOO6kI8")

	biceps = append(biceps, "TRX Suspension Curl")
	technique = append(technique, `Once the TRX suspension has been secured, grab hold of the handles, take a few steps forward, lean back, and curl your body weight up. To increase the difficulty of the exercise, you simply adjust your body position. The further you lean back, the more of your body weight you will resist during the movement. If you want to make the exercise easier, you can position your body to be more upright.
	`)
	videoURI = append(videoURI, "https://www.youtube.com/watch?v=c2DHtVkK1Fc&ab_channel=PhysiqueDevelopment")

	biceps = append(biceps, "EZ-Bar Reverse Curl")
	technique = append(technique, `Grip an EZ-bar with each hand while standing. Turn your wrists, so your palms are facing down (or best fit to the slanted part on the bar). Keep your arms tucked in at your sides and flex your elbows to curl the bar up towards your shoulders. Lower the bar back down with control.
	`)
	videoURI = append(videoURI, "https://www.youtube.com/watch?v=_ZHQMD_CdcQ")

	stmt := `
		INSERT INTO
			biceps(title, technique, videouri)
		VALUES  ($1, $2, $3)
	`

	for i := 0; i < len(biceps)-1; i++ {
		if row := m.DB.QueryRowContext(context.Background(), stmt, biceps[i], technique[i], videoURI[i]); row.Err() != nil {
			fmt.Printf("Error inserting into BICEPS table: %v", row.Err())
			return row.Err()
		}
	}

	return nil
}

// StoreAllBackExercises stores all exercises for back in database
func (m *DatabaseModel) StoreAllBackExercises() error {
	var back []string
	var technique []string
	var videoURI []string

	back = append(back, "Dumbbell Single-Arm Row")
	technique = append(technique, `
	Step back into a lunge position, with a soft bend in front knee & straight back leg. Leaning forward, rest hand on thigh, tighten core. Lower dumbbell to floor with full extension, maintain proper posture. Begin upward motion by sliding shoulder blade & lifting weight with elbow to ceiling, squeezing blade at end. Repeat for reps, switch sides, perform 2-3 sets with 1 min rest.
	`)
	videoURI = append(videoURI, "https://youtu.be/dFzUjzfih7k")

	back = append(back, "Inverted Row")
	technique = append(technique, `Set bar/rings around waist height, lower for more difficulty. Lie face up under bar, grip overhand wider than shoulder-width. Contract abs & butt, maintain straight line from ears to feet. Pull up to bar until chest touches, lower with proper form.`)
	videoURI = append(videoURI, "https://youtu.be/5W8F6MzZ8Rk")

	back = append(back, "Dumbbell Pullover")
	technique = append(technique, `Extend arms over chest, palms facing, elbows slightly bent. Inhale, extend weights back & over head with strong back & core. Take 3-4 secs to reach full extension behind head. Exhale, return to starting position.`)
	videoURI = append(videoURI, "https://youtu.be/ar-eEXkHLO4")

	back = append(back, "Single-Arm Eccentric Pulldown")
	technique = append(technique, `Set high pulley, kneel 3 ft from machine, grasp rope with left arm, square hips & shoulders, squeeze shoulder blades & abs. Use heavy weight. Pull weight down with both arms, tighten glutes & abs, release right arm, slowly raise weight with slight elbow bend for 3-5 secs.`)
	videoURI = append(videoURI, "https://youtu.be/r68DrqhNSis")

	back = append(back, "Lat Pulldown")
	technique = append(technique, `Set bar height so you can comfortably grasp with extended arms, allowing full range of motion. Adjust thigh pad if station has one. Use wide overhand knuckles-up grip. Pull bar down to chin level while exhaling, keeping upper torso stationary and feet flat. Squeeze shoulder blades and slowly return bar to starting position, controlling ascent. Repeat 8-12 reps, rest, then continue sets.`)
	videoURI = append(videoURI, "https://youtu.be/CAwf7n6Luuc")

	back = append(back, "Towel Grip Landmine Row")
	technique = append(technique, `Setup barbell properly, use landmine attachment or stick in corner with towel to protect wall. Feet hip-width apart, perpendicular to barbell. Brace core, exhale. Shove butt back, slight knee bend, hinge at hips. Grab barbell, pull shoulder blade back, bring elbow back. Repeat prescribed reps. Feel pumped lats & appreciate landmine setup.`)
	videoURI = append(videoURI, "https://www.youtube.com/watch?v=tZ6N_xSmVCg")

	back = append(back, "Bear Row to Gorilla Row")
	technique = append(technique, `
		Start in bear-plank position, with hands gripping kettlebells and back flat
		Drive left hand into kettlebell and row right bell to rib cage
		Squeeze abs and shoulder blades, keeping hips and shoulders square
		Jump forward to gorilla stance, row both bells to rib cage
		Jump back to bear-plank, repeat for 40 sec, rest for 40 sec.`)
	videoURI = append(videoURI, "Do not have video. Sorry :c")

	back = append(back, "Bent-Over Barbell Rows")
	technique = append(technique, `Stand with a shoulder-width stance.Grab the barbell, wider than shoulder-width, with an overhand grip.Bending your knees slightly, and your core tight, bend over at the waist keeping your lower back tight.Bending over until your upper body is at a 45-degree bend or lower, pull the bar up towards your lower chest.Keep your elbows as close to your sides as possible.At the top of the movement, you should feel like you are pinching your shoulder blades towards each other.`)
	videoURI = append(videoURI, "https://youtu.be/FWJR5Ve8bnQ")

	back = append(back, "Seated Cable Row")
	technique = append(technique, `Pull the handle and weight back toward the lower abdomen while trying not to use the momentum of the row too much by moving the torso backward with the arms. Target the middle to upper back by keeping your back straight and squeezing your shoulder blades together as you row, keeping your chest out. Return the handle forward under tension to full stretch, remembering to keep that back straight even though flexed at the hips. Repeat the exercise for the desired number of repetitions.`)
	videoURI = append(videoURI, "https://youtu.be/sP_4vybjVJs")

	back = append(back, "Pullup")
	technique = append(technique, `Exhale while pulling yourself up so your chin is level with the bar. Pause at the top. Lower yourself (inhaling as you go down) until your elbows are straight. Repeat the movement without touching the floor. Complete the number of repetitions your workout requires.`)
	videoURI = append(videoURI, "https://youtu.be/eGo4IYlbE5g")

	// T-bar row in the simulator

	back = append(back, "Deadlfit")
	technique = append(technique, `Stand behind the loaded barbell, with feet shoulder-width apart. Bar should be over shoelaces and touching shins. Sit back, keeping chest up and back straight. Grab the barbell with an "over-under" grip. Push hips forward to stand up, pressing feet into floor for stability and keeping core strong and back flat. Send hips back to return to starting position, bar on floor and against legs, chest up and looking forward. Release bar and stand up.`)
	videoURI = append(videoURI, "https://youtu.be/r4MzxtBKyNE")

	back = append(back, "Kettlebell Swings")
	technique = append(technique, `Give yourself enough space (4-5ft in front and a couple of feet behind) to perform the kettlebell swing. Avoid any breakable items in front of you. Place the kettlebell on the ground in front of you, with your feet hip-width apart and toes angled out. Keep your abs engaged and shoulders rolled back while bending your knees slightly. Reach for the kettlebell handle by pressing your hips back and tipping your torso forward, keeping your back straight. Grasp the handle with both hands and roll your shoulders back to control momentum. Exhale and use your glutes and hamstrings to rise to an upright position. The kettlebell should swing naturally to shoulder height. Inhale and swing it back toward the floor by pressing your hips back. Keep your neck aligned with your spine. Continue swinging, using your hips and glutes for power. Reduce the power gradually to return the kettlebell safely to the floor. Remember to hinge your hips, not use your arms or quads, for the movement.`)
	videoURI = append(videoURI, "https://youtu.be/sSESeQAir2M")

	back = append(back, "Kettlebelt Snatch")
	technique = append(technique, `Grip the kettlebell with your fingers, sit back to load your hips, and keep your feet hip-to-shoulder distance apart. The kettlebell swings back between your legs as you stand and exhale sharply from the mouth. Keep arm close to the body, extend hips and knees to allow kettlebell to pull your arm upward, and accelerate it vertically with a rapid hip pull and shrug of the trapezoids. Release fingers, insert palm, and exhale as the arm locks out overhead. Shift weight to opposite leg and inhale on downswing. Connect arm to torso and pull hand towards you to change to hook grip. Repeat for desired reps, inhaling on downswing and exhaling on backswing.`)
	videoURI = append(videoURI, "https://youtu.be/Pm-b2XFeABA")

	back = append(back, "Hyperextension")
	technique = append(technique, `Setup in a hyperextension machine with your feet anchored and torso roughly perpendicular to your legs at a 45 degree angle. Begin in a hinged position with your arms crossed and initiate the movement by flexing your glutes. Extend the hips and finish with your body in a straight line. Repeat for the desired number of repetitions.`)
	videoURI = append(videoURI, "https://www.youtube.com/watch?v=ph3pddpKzzw")

	back = append(back, "Block pull behind the head")
	technique = append(technique, `To perform the upper block pull down exercise: Adjust the leg fixing roller to prevent the weight from pulling your body up. Place your thighs under the bolsters and press your feet to the floor at knee level, position your pelvis on the seat so the upper block with the bar is over your head. Grip the bar with a medium grip, keeping your little finger at the beginning of the scar. Exhale and pull the bar down behind your head, keeping the forearms perpendicular to the floor and parallel to each other. Use your back muscles to perform the exercise, not your arms. Avoid straining your shoulders and forearms, and make sure your back muscles are being worked. Pause at the bottom point, inhale at the top point but don't fully relax your shoulder girdle. Do not stretch your back or shoulders too much to prevent injury. Finish the workout with a series of inhalation and exhalation when your elbows are extended at the top point.`)
	videoURI = append(videoURI, "Do not have video, sorry :c")

	back = append(back, "T-bar row in the simulator")
	technique = append(technique, `Stand on the plate, lean forward, and grasp the bar with an overhand grip in straight arms. Inhale and pull the bar towards you as high as possible. With control, lower the bar back to the starting position. T-bar rows is a variant of standing rows where the weight and bar path is fixed.`)
	videoURI = append(videoURI, "https://youtu.be/j3Igk5nyZE4")

	stmt := `INSERT INTO back(title, technique, videouri) VALUES ($1, $2, $3)`

	for i := 0; i < len(back)-1; i++ {
		if row := m.DB.QueryRowContext(context.Background(), stmt, back[i], technique[i], videoURI[i]); row.Err() != nil {
			fmt.Printf("Error inserting into BACK table: %v", row.Err())
			return row.Err()
		}
	}

	return nil
}

// StoreAllChestExercises stores all exercises for chest in database
func (m *DatabaseModel) StoreAllChestExercises() error {
	m.TruncateTables()
	var chest []string
	var technique []string
	var videoURI []string

	chest = append(chest, "Barbell Bench Press")
	technique = append(technique, `Lie on a flat bench with your eyes directly under the bar. Grasp the barbell with a grip slightly wider than shoulder-width apart. Lift the bar off the rack and lower it to your chest, keeping your elbows at a 45-degree angle to your body. Your arms should be perpendicular to the floor. Pause briefly, then press the barbell back up to the starting position, locking out your arms at the top. Keep your core tight and your back flat against the bench throughout the exercise. Repeat for the desired number of repetitions.`)
	videoURI = append(videoURI, "https://youtu.be/rT7DgCr-3pg")

	chest = append(chest, "Dumbbell Bench Press")
	technique = append(technique, `Lie on a flat bench with your eyes directly under the dumbbells. Grasp the dumbbells with a neutral grip, palms facing each other. Lift the dumbbells up and position them above your chest with your arms straight. Lower the dumbbells to your chest, keeping your elbows at a 45-degree angle to your body. Your arms should be perpendicular to the floor. Pause briefly, then press the dumbbells back up to the starting position, locking out your arms at the top. Keep your core tight and your back flat against the bench throughout the exercise. Repeat for the desired number of repetitions.`)
	videoURI = append(videoURI, "https://youtu.be/VmB1G1K7v94")

	chest = append(chest, "Cable Crossover")
	technique = append(technique, `Stand in the center of a cable crossover machine with your feet shoulder-width apart. Attach two D-handle attachments to the high pulleys.
	Grasp each handle and step forward to create tension in the cables. Your arms should be straight out in front of you, with your palms facing each other.
	Keeping your arms straight, bring your hands together in front of your chest, crossing one arm over the other. Your palms should touch or almost touch at the top of the movement.
	Pause briefly, then return to the starting position with your arms straight. Repeat for the desired number of repetitions.`)
	videoURI = append(videoURI, "https://youtu.be/taI4XduLpTk")

	chest = append(chest, "Incline Dumbbell Press")
	technique = append(technique, `Lie on an incline bench set to a 45-degree angle with your eyes directly under the dumbbells. Grasp the dumbbells with a neutral grip, palms facing each other.
	Lift the dumbbells up and position them above your chest with your arms straight.
	Lower the dumbbells to your chest, keeping your elbows at a 45-degree angle to your body. Your arms should be perpendicular to the floor.
	Pause briefly, then press the dumbbells back up to the starting position, locking out your arms at the top. Keep your core tight and your back flat against the bench throughout the exercise.
	Repeat for the desired number of repetitions.`)
	videoURI = append(videoURI, "https://youtu.be/8iPEnn-ltC8")

	chest = append(chest, " Dumbbell Half Flye")
	technique = append(technique, `
	Lie on a flat bench with your eyes directly under the dumbbells. Grasp the dumbbells with a neutral grip, palms facing each other.
	Lift the dumbbells up and position them above your chest with your arms straight.
	Keeping your arms slightly bent, lower the dumbbells out to your sides in a semicircular motion, until they are at chest level. Your palms should be facing each other.
	Pause briefly, then return to the starting position, lifting the dumbbells back up in a semicircular motion.
	Repeat for the desired number of repetitions.`)
	videoURI = append(videoURI, "https://youtu.be/eozdVDA78K0")

	chest = append(chest, "Dumbbell Flye")
	technique = append(technique, `
	Lie on a flat bench with your eyes directly under the dumbbells. Grasp the dumbbells with a neutral grip, palms facing each other.
	Lift the dumbbells up and position them above your chest with your arms straight.
	Keeping your arms slightly bent, lower the dumbbells out to your sides in a semicircular motion, until they are at chest level. Your palms should be facing each other.
	Pause briefly, then return to the starting position, lifting the dumbbells back up in a semicircular motion.
	Repeat for the desired number of repetitions.`)
	videoURI = append(videoURI, "https://youtu.be/QENKPHhQVi4")

	chest = append(chest, "Incline Dumbbell Flye")
	technique = append(technique, `
		Lie on an incline bench set to a 45-degree angle with your eyes directly under the dumbbells. Grasp the dumbbells with a neutral grip, palms facing each other.
		Lift the dumbbells up and position them above your chest with your arms straight.
		Keeping your arms slightly bent, lower the dumbbells out to your sides in a semicircular motion, until they are at chest level. Your palms should be facing each other.
		Pause briefly, then return to the starting position, lifting the dumbbells back up in a semicircular motion.
		Repeat for the desired number of repetitions.`)
	videoURI = append(videoURI, "https://youtu.be/ajdFwa-qM98")

	chest = append(chest, "Landmine Press")
	technique = append(technique, `Load one end of a barbell into a landmine attachment or secure it in a corner. Stand facing the barbell with your feet shoulder-width apart.
		Grasp the barbell with both hands, one hand on each side of the bar, with a neutral grip (palms facing each other).
		Hold the bar at chest height, with your elbows bent at 90 degrees.
		Keeping your core tight, press the bar up and away from your chest until your arms are straight.
		Pause briefly, then return to the starting position. Repeat for the desired number of repetitions.`)
	videoURI = append(videoURI, "https://youtu.be/7i64SnEJv6A")

	chest = append(chest, "Pullover")
	technique = append(technique, `
		Lie on a flat bench with your head and shoulders hanging off the end. Grasp a dumbbell with both hands and hold it above your chest.
		Keeping your arms straight, lower the dumbbell behind your head, keeping it close to your chest.
		Pause briefly when the dumbbell is directly behind your head, then return to the starting position, lifting the dumbbell back up to the starting position above your chest.
		Repeat for the desired number of repetitions.`)
	videoURI = append(videoURI, "https://youtu.be/2ut5wBTORJY")

	chest = append(chest, "Plate Pressout")
	technique = append(technique, `
		Stand holding a weight plate with both hands, at chest height, with your palms facing each other.
		Keeping your arms straight, press the weight plate out in front of you until your arms are fully extended.
		Pause briefly, then return to the starting position, pulling the weight plate back towards your chest.
		Repeat for the desired number of repetitions.`)
	videoURI = append(videoURI, "https://youtu.be/2XmNv4T_Jfo")

	chest = append(chest, "Pushup")
	technique = append(technique, `
		Start in a plank position, with your hands placed slightly wider than shoulder-width apart and your feet close together.
		Lower your body towards the ground, keeping your core tight and your back straight. Your elbows should bend and form a 90-degree angle as you descend.
		Pause briefly when your chest is close to the ground, then push yourself back up to the starting position, straightening your arms.
		Repeat for the desired number of repetitions.`)
	videoURI = append(videoURI, "https://youtu.be/IODxDxX7oi4")

	chest = append(chest, "Band-Resisted Flye technique")
	technique = append(technique, `
		Attach a resistance band to a sturdy anchor point, such as a power rack or a heavy object, at about chest height. Stand facing the anchor point, holding the handles of the band in each hand.
		Take a step or two away from the anchor point, so that the band is taut. Stand with your feet shoulder-width apart and your arms extended in front of you, holding the handles of the band.
		Keeping your arms straight, raise the handles of the band out to your sides in a semicircular motion, until they are at chest level. Your palms should be facing each other.
		Pause briefly, then return to the starting position, lowering the band handles back in a semicircular motion.
		Repeat for the desired number of repetitions.`)
	videoURI = append(videoURI, "https://youtu.be/gfS03SvURHc")

	chest = append(chest, "Pec Deck")
	technique = append(technique, `
		Sit down on the Pec Deck machine and adjust the seat height so that the handles are level with your chest.
		Grasp the handles with your palms facing each other and your arms extended.
		Keeping your arms straight, bring the handles together in front of your chest.]
		Pause briefly, then return to the starting position, extending your arms out to the sides.
		Repeat for the desired number of repetitions.`)
	videoURI = append(videoURI, "https://youtu.be/u56jywgbvE4")

	chest = append(chest, "Wide-Grip Dips")
	technique = append(technique, `
		Find a pair of parallel bars and position yourself between them, with your palms grasping the bars and your arms extended. Your hands should be positioned wider than shoulder-width apart.
		Lower your body towards the ground by bending your arms, keeping your core tight and your back straight.
		Pause briefly when your upper arms are parallel to the ground, then press yourself back up to the starting position, straightening your arms.
		Repeat for the desired number of repetitions.`)
	videoURI = append(videoURI, "https://youtu.be/yN6Q1UI_xkE")

	chest = append(chest, "Smith Machine Incline Press")
	technique = append(technique, `
		Set the barbell of the Smith Machine to an incline position, usually around 30-45 degrees. Lie down on the bench, positioning yourself underneath the bar.
		Grasp the bar with your palms facing away from you, slightly wider than shoulder-width apart. Your feet should be flat on the ground.
		Unrack the barbell and lower it to your chest, keeping your core tight and your back straight.
		Pause briefly, then press the bar back up to the starting position, straightening your arms.
		Repeat for the desired number of repetitions.`)
	videoURI = append(videoURI, "https://youtu.be/b8DqTO6ak0k")

	chest = append(chest, "Low-Cable Crossover")
	technique = append(technique, `
		Stand between two cable machines, with the pulleys set at about hip height. Grasp the handles of each cable with your palms facing each other.
		Take a step forward, so that the cables are taut, and stand with your feet shoulder-width apart. Your arms should be extended in front of you.
		Keeping your arms straight, bring the cables across your body towards each other, crossing in front of your chest. Your palms should come together in front of your chest.
		Pause briefly, then return to the starting position, extending your arms back out in front of you.
		Repeat for the desired number of repetitions.`)
	videoURI = append(videoURI, "https://youtu.be/M1N804yWA-8")

	chest = append(chest, "Chest Press Machine")
	technique = append(technique, `
		Sit down on the Chest Press Machine and adjust the seat height so that the handles are level with your chest.
		Grasp the handles with your palms facing away from you, slightly wider than shoulder-width apart.
		Unrack the weight by pressing the handles away from your chest. Your arms should be extended in front of you.
		Lower the handles towards your chest, keeping your core tight and your back against the seat.
		Pause briefly, then press the handles back up to the starting position, extending your arms.
		Repeat for the desired number of repetitions.`)
	videoURI = append(videoURI, "https://youtu.be/xUm0BiZCWlQ")

	stmt := `INSERT INTO chest(title, technique, videouri) VALUES ($1, $2, $3)`

	for i := 0; i < len(chest)-1; i++ {
		if len(chest[i]) >= 50 {
			chest[i] = chest[i][:50]
		}
		if _, err := m.DB.Query(stmt, chest[i], technique[i], videoURI[i]); err != nil {
			fmt.Printf("\nError inserting into CHEST table: %v\n", err)
			return err
		}
		//fmt.Println("Value ", i, " inserted")
	}

	return nil
}

func (m *DatabaseModel) SaveAllShouldersExercises() error {
	var title []string
	var technique []string
	var videoURI []string

	title = append(title, "Seated Dumbbell Shoulder Press")
	technique = append(technique, `
		1. Sit on a bench with back support and hold a dumbbell in each hand at shoulder level, palms facing forward.
		2. Press the dumbbells upward until your arms are fully extended overhead.
		3. Slowly lower the dumbbells back to the starting position.
		4. Repeat for the desired number of repetitions.
		5. Keep your core engaged and maintain proper posture throughout the exercise.
	`)
	videoURI = append(videoURI, "https://youtu.be/-urozvDwN7I")

	title = append(title, "Barbell Overhead Press")
	technique = append(technique, `
		1. Stand with your feet shoulder-width apart and grip the barbell with your hands slightly wider than shoulder-width apart, palms facing forward.
		2. Lift the barbell off the rack and position it at shoulder level with your elbows bent and the bar resting on your clavicles.
		3. Press the barbell upward until your arms are fully extended overhead.
		4. Lower the barbell back to the starting position, keeping your core engaged and maintaining proper posture throughout the exercise.
		5. Repeat for the desired number of repetitions.
		6. Be mindful of your breathing and avoid arching your lower back excessively.
	`)
	videoURI = append(videoURI, "https://youtu.be/CnBmiBqp-AI")

	title = append(title, "Dumbbell Lateral Raise")
	technique = append(technique, `
		1. Stand with your feet shoulder-width apart and hold a dumbbell in each hand at your sides, palms facing inward.
		2. Keeping a slight bend in your elbows, raise your arms out to the sides until they are parallel to the floor.
		3. Pause briefly at the top of the movement, then lower the dumbbells back to the starting position.
		4. Repeat for the desired number of repetitions.
		5. Focus on using your shoulder muscles to lift the dumbbells, rather than relying on momentum.
		6. Keep your core engaged and maintain proper posture throughout the exercise.
	`)
	videoURI = append(videoURI, "https://youtu.be/d32l9mSIoYc")

	title = append(title, "Arnold Press")
	technique = append(technique, `
		1. Sit on a bench with back support and hold a dumbbell in each hand at shoulder level, palms facing your body.
		2. Rotate your palms as you press the dumbbells overhead, so that at the top of the movement your palms are facing forward.
		3. Reverse the motion and rotate your palms back to the starting position as you lower the dumbbells back to shoulder level.
		4. Repeat for the desired number of repetitions.
		5. Keep your core engaged and maintain proper posture throughout the exercise.
	`)
	videoURI = append(videoURI, "https://youtu.be/cn1bYxO2oAE")

	title = append(title, "Front Raise")
	technique = append(technique, `
		1. Stand with your feet shoulder-width apart and hold a dumbbell in each hand with your palms facing your thighs.
		2. Keeping your arms straight and maintaining a slight bend in your elbows, lift the dumbbells directly in front of you until they are at shoulder level.
		3. Pause briefly at the top of the movement, then lower the dumbbells back to the starting position.
		4. Repeat for the desired number of repetitions.
		5. Focus on using your shoulder muscles to lift the dumbbells, rather than relying on momentum.
		6. Keep your core engaged and maintain proper posture throughout the exercise.
	`)
	videoURI = append(videoURI, "https://youtu.be/tb-XZ1Cfah4")

	stmt := `INSERT INTO shoulders(title, technique, videouri) VALUES ($1, $2, $3)`

	for i := 0; i < len(title)-1; i++ {
		if len(title[i]) >= 50 {
			title[i] = title[i][:50]
		}
		if _, err := m.DB.Query(stmt, title[i], technique[i], videoURI[i]); err != nil {
			fmt.Printf("\nError inserting into CHEST table: %v\n", err)
			return err
		}
	}

	return nil
}

//  StoreAllTricepsExercises stores all exercises for glutes in database
//  func (db *DatabaseModel) StoreAllTricepsExercises() error {
//	var triceps []string
//	var technique []string
//	var videoURI []string
//
//	return nil
//}

// StoreAllShoulderExercises stores all exercises for glutes in database
//func (db *DatabaseModel) StoreAllShoulderExercises() error {
//	var shoulder []string
//	var technique []string
//	var videoURI []string
//
//	return nil
//}

// GetAllExercisesFromTable function is used for getting all exercises records from table
// when called with argument "chest" returns all exercises from chest table
func (m *DatabaseModel) GetAllExercisesFromTable(tableName string) ([]Exercise, error) {
	var exercises []Exercise
	var exercise Exercise
	stmt := "SELECT * FROM $1"

	row, err := m.DB.Query(stmt, tableName)
	if err != nil {
		fmt.Printf("Error querying row: %v", err)
		return exercises, err
	}

	// otherwise, exercises are received
	err = row.Scan(
		&exercise.Title,
		&exercise.Technique,
		&exercise.VideoURI)
	if err != nil {
		fmt.Printf("Error scanning row body: %v", err)
		return exercises, err
	}

	return exercises, nil
}

// GetOneRandomExercise generates random number in range(0-MuscleToCount[muscle]), where MuscleToCount[muscle]
// mapping represents the last index of exercise in table.
func (m *DatabaseModel) GetOneRandomExercise(table string) (Exercise, error) {
	var emptyExercise Exercise
	var exercise Exercise
	var err error

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	id := r.Intn(MuscleToCount[table] + 1)

	exercise, err = m.GetExerciseById(id, table)
	if err != nil {
		return emptyExercise, fmt.Errorf("error getting exercise from table %s by id[%d]: %w", table, id, err)
	}
	if exercise == emptyExercise {
		exercise, err = m.GetExerciseById(id-1, table)
		if err != nil {
			return emptyExercise, fmt.Errorf("error getting exercise from table %s by id[%d]: %w", table, id-1, err)
		}
	}

	return exercise, nil
}

func (m *DatabaseModel) GenerateXRandomExercises(table string, num int) ([]Exercise, error) {
	if StringInArray(table, muscles) == false {
		return nil, fmt.Errorf("table %s not found in MuscleGroups", table)
	}

	var exercises []Exercise
	for i := 0; i < num; i++ {
		ex, err := m.GetOneRandomExercise(table)
		if err != nil {
			return nil, err
		}
		exercises = append(exercises, ex)
	}
	return exercises, nil
}

func StringInArray(s string, arr []string) bool {
	for i := 0; i <= len(arr)-1; i++ {
		if s == arr[i] {
			return true
		}
	}
	return false
}
