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

	type messageEvent struct {
		greeting map[string]string
    commands map[string]string
	}
	var r messageEvent

	r.greeting = map[string]string{
		"hello":     "Greetings!",
		"whats up?": "Not much",
		"hi":        "Hello!",
		"hey":       "hey",
	}

	r.commands = map[string]string{
		"begin battle": beginBattle(),
	}


	input := strings.ToLower(m.Content)

	if value, ok := r.greeting[input]; ok {
		s.ChannelMessageSend(m.ChannelID, value+" "+m.Author.Username)
	}
	if value, ok := r.commands[input]; ok {
		s.ChannelMessageSend(m.ChannelID, value)
	}
}

type entity interface {
	attack(target entity, dmg int)
	getCurrentDmgDone() int
	getInitiative() int
}

func beginBattle() string {
	// prompt the players for their initative total
  s.ChannelMessageSend(m.ChannelID, "What is the initiative total?")
  // warn them of timer 3min
  s.ChannelMessageSend(m.ChannelID, "Timer: 3 minutes")
  // ask the dm for number of monsters to add to the tracker
  s.ChannelMessageSend(m.ChannelID, "How many monsters do you want to add?")
  // add the monsters to the tracker
  
  // give the monsters a color or unique name
  
  // add the monsters to the tracker
  // collect the player name from their meta data and match it with thier input
  // sort the order by value
  // begin the tracker
  
	string := "begin battle"
	return string
}
