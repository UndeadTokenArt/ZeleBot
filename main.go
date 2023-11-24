package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var Token string

func init() {
	Token = os.Getenv("token")
}

func main() {
	if Token == "" {
		fmt.Println("Bot Token not found in the environment variable 'token'.")
		return
	}
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	mInput := strings.Fields(strings.ToLower(m.Content))

	if len(mInput) < 1 {
		// Ignore messages with no content or command
		return
	}

	command := mInput[0]

	switch {
	case command == "!addcity" || strings.HasPrefix(command, "!addcity"):
		handleAddCity(s, m, mInput)
	case command == "!getcities" || strings.HasPrefix(command, "!getcities"):
		handleGetCities(s, m)

	default:
		s.ChannelMessageSend(m.ChannelID, "that's not a command I know")
	}
}

type entity interface {
	attack(target entity, dmg int)
	getCurrentDmgDone() int
	getInitiative() int
}

type Region struct {
	Cities
}

type Cities struct {
	name        []string
	Factions    []string
	Shops       []string
	Civic       []string
	Resedential []string
	Workshops   []string
}

var rulian Region

func handleAddCity(s *discordgo.Session, m *discordgo.MessageCreate, mInput []string) {
	if len(mInput) < 2 {
		// Insufficient arguments for !addCity
		s.ChannelMessageSend(m.ChannelID, "Insufficient arguments for !addCity")
		return
	}
	cityName := strings.Join(mInput[1:], " ")
	rulian.Cities.name = append(rulian.Cities.name, cityName)
	s.ChannelMessageSend(m.ChannelID, "City Added!")

}

func handleGetCities(s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(rulian.Cities.name) < 1 {
		s.ChannelMessageSend(m.ChannelID, "No cities listed")
	}
	cityList := strings.Join(rulian.Cities.name, "\n")
	s.ChannelMessageSend(m.ChannelID, cityList)
}
