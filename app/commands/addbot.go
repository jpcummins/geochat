package commands

import (
	"github.com/jpcummins/geochat/app/types"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var botNames = []string{"Sophia", "Noah", "Emma", "Liam", "Olivia", "Jacob", "Isabella", "Mason", "Ava", "William", "Mia", "Ethan", "Emily", "Michael", "Abigail", "Alexander", "Madison", "Jayden", "Elizabeth", "Daniel", "Charlotte", "Elijah", "Avery", "Aiden", "Sofia", "James", "Chloe", "Benjamin", "Ella", "Matthew", "Harper", "Jackson", "Amelia", "Logan", "Aubrey", "David", "Addison", "Anthony", "Evelyn", "Joseph", "Natalie", "Joshua", "Grace", "Andrew", "Hannah", "Lucas", "Zoey", "Gabriel", "Victoria", "Samuel", "Lillian", "Christopher", "Lily", "John", "Brooklyn", "Dylan", "Samantha", "Isaac", "Layla", "Ryan", "Zoe", "Nathan", "Audrey", "Carter", "Leah", "Caleb", "Allison", "Luke", "Anna", "Christian", "Aaliyah", "Hunter", "Savannah", "Henry", "Gabriella", "Owen", "Camila", "Landon", "Aria", "Jack", "Kaylee", "Wyatt", "Scarlett", "Jonathan", "Hailey", "Eli", "Arianna", "Isaiah", "Riley", "Sebastian", "Alexis", "Jaxon", "Nevaeh", "Julian", "Sarah", "Brayden", "Claire", "Gavin", "Sadie", "Levi", "Peyton", "Aaron", "Aubree", "Oliver", "Serenity", "Jordan", "Ariana", "Nicholas", "Genesis", "Evan", "Penelope", "Connor", "Alyssa", "Charles", "Bella", "Jeremiah", "Taylor", "Cameron", "Alexa", "Adrian", "Kylie", "Thomas", "Mackenzie", "Robert", "Caroline", "Tyler", "Kennedy", "Colton", "Autumn", "Austin", "Lucy", "Jace", "Ashley", "Angel", "Madelyn", "Dominic", "Violet", "Josiah", "Stella", "Brandon", "Brianna", "Ayden", "Maya", "Kevin", "Skylar", "Zachary", "Ellie", "Parker", "Julia", "Blake", "Sophie", "Jose", "Katherine", "Chase", "Mila", "Grayson", "Khloe", "Jason", "Paisley", "Ian", "Annabelle", "Bentley", "Alexandra", "Adam", "Nora", "Xavier", "Melanie", "Cooper", "London", "Justin", "Gianna", "Nolan", "Naomi", "Hudson", "Eva", "Easton", "Faith", "Jase", "Madeline", "Carson", "Lauren", "Nathaniel", "Nicole", "Jaxson", "Ruby", "Kayden"}

type addBot struct{}

func (b *addBot) Execute(args string, user types.User) error {
	split := strings.Split(args, " ")
	bots := 10
	ttl := 60

	if len(split) >= 1 {
		n, err := strconv.Atoi(split[0])
		if err == nil {
			bots = n
		}
	}

	if len(split) >= 2 {
		n, err := strconv.Atoi(split[1])
		if err == nil {
			ttl = n
		}
	}

	zone := user.Zone()
	world := zone.World()

	for i := 0; i < bots; i++ {
		go func(num int) {
			name := botNames[rand.Intn(len(botNames))]

			lat := zone.SouthWest().Lat() + (rand.Float64() * (zone.NorthEast().Lat() - zone.SouthWest().Lat()))
			lng := zone.SouthWest().Lng() + (rand.Float64() * (zone.NorthEast().Lng() - zone.SouthWest().Lng()))

			bot, err := world.NewUser(name+strconv.Itoa(num), name, lat, lng)
			if err != nil {
				return
			}

			if _, err := zone.Join(bot); err != nil {
				println(err.Error())
			}

			// Bot event handler
			go func() {
				timer := time.NewTimer(time.Duration(ttl) * time.Second)
				for {
					select {
					case <-timer.C:
						zone.Leave(bot)
						return
					}
				}
			}()
		}(i)
	}

	return nil
}
