package chat

import (
	"errors"
	"math/rand"
	"strconv"
	"time"
)

var botNames = []string{"Sophia", "Noah", "Emma", "Liam", "Olivia", "Jacob", "Isabella", "Mason", "Ava", "William", "Mia", "Ethan", "Emily", "Michael", "Abigail", "Alexander", "Madison", "Jayden", "Elizabeth", "Daniel", "Charlotte", "Elijah", "Avery", "Aiden", "Sofia", "James", "Chloe", "Benjamin", "Ella", "Matthew", "Harper", "Jackson", "Amelia", "Logan", "Aubrey", "David", "Addison", "Anthony", "Evelyn", "Joseph", "Natalie", "Joshua", "Grace", "Andrew", "Hannah", "Lucas", "Zoey", "Gabriel", "Victoria", "Samuel", "Lillian", "Christopher", "Lily", "John", "Brooklyn", "Dylan", "Samantha", "Isaac", "Layla", "Ryan", "Zoe", "Nathan", "Audrey", "Carter", "Leah", "Caleb", "Allison", "Luke", "Anna", "Christian", "Aaliyah", "Hunter", "Savannah", "Henry", "Gabriella", "Owen", "Camila", "Landon", "Aria", "Jack", "Kaylee", "Wyatt", "Scarlett", "Jonathan", "Hailey", "Eli", "Arianna", "Isaiah", "Riley", "Sebastian", "Alexis", "Jaxon", "Nevaeh", "Julian", "Sarah", "Brayden", "Claire", "Gavin", "Sadie", "Levi", "Peyton", "Aaron", "Aubree", "Oliver", "Serenity", "Jordan", "Ariana", "Nicholas", "Genesis", "Evan", "Penelope", "Connor", "Alyssa", "Charles", "Bella", "Jeremiah", "Taylor", "Cameron", "Alexa", "Adrian", "Kylie", "Thomas", "Mackenzie", "Robert", "Caroline", "Tyler", "Kennedy", "Colton", "Autumn", "Austin", "Lucy", "Jace", "Ashley", "Angel", "Madelyn", "Dominic", "Violet", "Josiah", "Stella", "Brandon", "Brianna", "Ayden", "Maya", "Kevin", "Skylar", "Zachary", "Ellie", "Parker", "Julia", "Blake", "Sophie", "Jose", "Katherine", "Chase", "Mila", "Grayson", "Khloe", "Jason", "Paisley", "Ian", "Annabelle", "Bentley", "Alexandra", "Adam", "Nora", "Xavier", "Melanie", "Cooper", "London", "Justin", "Gianna", "Nolan", "Naomi", "Hudson", "Eva", "Easton", "Faith", "Jase", "Madeline", "Carson", "Lauren", "Nathaniel", "Nicole", "Jaxson", "Ruby", "Kayden"}

func addBot(args []string, subscription *Subscription) (string, error) {

	if len(args) < 2 {
		return "", errors.New("Expected additional arguments")
	}

	number, err := strconv.Atoi(args[0])  // number of bots
	timeout, err := strconv.Atoi(args[1]) // timeout

	if err != nil {
		return "", err
	}

	for i := 0; i < number; i++ {
		go func(num int) {
			name := botNames[rand.Intn(len(botNames))]
			bot := &User{Id: name, Name: name, IsBot: true}

			botSubscription := subscription.Zone.Subscribe(bot)

			// Bot event handler
			go func() {
				timer := time.NewTimer(time.Duration(timeout) * time.Second)
				for {
					select {
					case <-timer.C:
						println("Unsubscribe")
						botSubscription.Zone.Unsubscribe(botSubscription)
						return
					case <-subscription.Events:
					}
				}
			}()
		}(i)
	}

	return "ok", nil
}
