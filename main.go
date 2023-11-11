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

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	type messageListener struct {
		greeting map[string]string
		handlers struct{}
	}

	var r messageListener

	r.greeting = map[string]string{
		"hello":     "Greetings!",
		"whats up?": "Not much",
		"hi":        "Hello!",
		"hey":       "hey",
	}

	input := strings.ToLower(m.Content)

	if value, ok := r.greeting[input]; ok {
		s.ChannelMessageSend(m.ChannelID, value+" "+m.Author.Username)
	}
}

type monster struct {
	name      string
	damage    int
	initative int
}
type player struct {
	name          string
	initative     int
	currentHealth int
	maxHealth     int
	AC            int
}

type entity interface {
	attack(target entity, dmg int)
	getCurrentDmgDone() int
}

func (p *player) attack(target entity, dmg int) {
	switch v := target.(type) {
	case *player:
		v.currentHealth = v.currentHealth - dmg
	case *monster:
		v.damage = v.damage + dmg
	}
}

func (m *monster) attack(target entity, dmg int) {
	switch v := target.(type) {
	case *player:
		v.currentHealth = v.currentHealth - dmg
	case *monster:
		v.damage = v.damage + dmg
	}
}

func (p *player) getCurrentDmgDone() int {
	return 0 // Implement the logic if needed
}

func (m *monster) getCurrentDmgDone() int {
	return m.damage
}
