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
    commands map[string]func
		greeting map[string]string
<<<<<<< HEAD
		commands map[string]string
=======
>>>>>>> 81e7a89d770dd06b697144b471b3f04b8ab74710
	}
	var r messageEvent

	r.greeting = map[string]string{
		"hello":     "Greetings!",
		"whats up?": "Not much",
		"hi":        "Hello!",
		"hey":       "hey",
	}

<<<<<<< HEAD
	r.commands = map[string]string{
		"begin battle": beginBattle(),
	}
=======
  r.commands = map[string]func{
    "begin combat" : beginCombat,
    }

>>>>>>> 81e7a89d770dd06b697144b471b3f04b8ab74710

	input := strings.ToLower(m.Content)

	if value, ok := r.greeting[input]; ok {
		s.ChannelMessageSend(m.ChannelID, value+" "+m.Author.Username)
	}
	if value, ok := r.commands[input]; ok {
		s.ChannelMessageSend(m.ChannelID, value)
	}
<<<<<<< HEAD
=======
  if value, ok := r.command[input]; ok {
    s.ChannelMessageSend(m.ChannelID, value)
>>>>>>> 81e7a89d770dd06b697144b471b3f04b8ab74710
}

type entity interface {
	attack(target entity, dmg int)
	getCurrentDmgDone() int
	getInitiative() int
}

<<<<<<< HEAD
func beginBattle() string {
	string := "begin battle"
	return string
}
=======
func beginCombat(s *discordgo.Session, m *discordgo.MessageCreate) {
  for time.Sleep(time.second * 180) {
    s.ChannelMessageSend(m.ChannelID, "Starting initative tracker")
    s.channelMessageSend(m.ChannelID, "Please input your initative total, you have 3 minutes to comply")
    
  }
}



>>>>>>> 81e7a89d770dd06b697144b471b3f04b8ab74710
