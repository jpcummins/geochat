package commands

import (
	"github.com/jpcummins/geochat/app/types"
	"math/rand"
	"strconv"
	"strings"
)

var botNames = []string{"Sophia", "Noah", "Emma", "Liam", "Olivia", "Jacob", "Isabella", "Mason", "Ava", "William", "Mia", "Ethan", "Emily", "Michael", "Abigail", "Alexander", "Madison", "Jayden", "Elizabeth", "Daniel", "Charlotte", "Elijah", "Avery", "Aiden", "Sofia", "James", "Chloe", "Benjamin", "Ella", "Matthew", "Harper", "Jackson", "Amelia", "Logan", "Aubrey", "David", "Addison", "Anthony", "Evelyn", "Joseph", "Natalie", "Joshua", "Grace", "Andrew", "Hannah", "Lucas", "Zoey", "Gabriel", "Victoria", "Samuel", "Lillian", "Christopher", "Lily", "John", "Brooklyn", "Dylan", "Samantha", "Isaac", "Layla", "Ryan", "Zoe", "Nathan", "Audrey", "Carter", "Leah", "Caleb", "Allison", "Luke", "Anna", "Christian", "Aaliyah", "Hunter", "Savannah", "Henry", "Gabriella", "Owen", "Camila", "Landon", "Aria", "Jack", "Kaylee", "Wyatt", "Scarlett", "Jonathan", "Hailey", "Eli", "Arianna", "Isaiah", "Riley", "Sebastian", "Alexis", "Jaxon", "Nevaeh", "Julian", "Sarah", "Brayden", "Claire", "Gavin", "Sadie", "Levi", "Peyton", "Aaron", "Aubree", "Oliver", "Serenity", "Jordan", "Ariana", "Nicholas", "Genesis", "Evan", "Penelope", "Connor", "Alyssa", "Charles", "Bella", "Jeremiah", "Taylor", "Cameron", "Alexa", "Adrian", "Kylie", "Thomas", "Mackenzie", "Robert", "Caroline", "Tyler", "Kennedy", "Colton", "Autumn", "Austin", "Lucy", "Jace", "Ashley", "Angel", "Madelyn", "Dominic", "Violet", "Josiah", "Stella", "Brandon", "Brianna", "Ayden", "Maya", "Kevin", "Skylar", "Zachary", "Ellie", "Parker", "Julia", "Blake", "Sophie", "Jose", "Katherine", "Chase", "Mila", "Grayson", "Khloe", "Jason", "Paisley", "Ian", "Annabelle", "Bentley", "Alexandra", "Adam", "Nora", "Xavier", "Melanie", "Cooper", "London", "Justin", "Gianna", "Nolan", "Naomi", "Hudson", "Eva", "Easton", "Faith", "Jase", "Madeline", "Carson", "Lauren", "Nathaniel", "Nicole", "Jaxson", "Ruby", "Kayden"}

type addBot struct{}

func (b *addBot) Execute(args string, user types.User, world types.World) error {
	split := strings.Split(args, " ")
	bots := 10

	if len(split) >= 1 {
		n, err := strconv.Atoi(split[0])
		if err == nil {
			bots = n
		}
	}

	zoneID := user.ZoneID()
	zone, err := world.Zones().Zone(zoneID)
	if err != nil {
		return err
	}

	for i := 0; i < bots; i++ {
		name := botNames[rand.Intn(len(botNames))]
		lat := zone.SouthWest().Lat() + (rand.Float64() * (zone.NorthEast().Lat() - zone.SouthWest().Lat()))
		lng := zone.SouthWest().Lng() + (rand.Float64() * (zone.NorthEast().Lng() - zone.SouthWest().Lng()))
		bot, err := world.NewUser(strconv.Itoa(i) + strconv.Itoa(rand.Intn(1000)))
		bot.SetFirstName(name)
		bot.SetLastName("Bot")
		bot.SetName(name + " Bot")
		bot.SetLocation(lat, lng)
		if err != nil {
			return err
		}

		if _, err := world.Join(bot); err != nil {
			return err
		}
	}

	return nil
}
