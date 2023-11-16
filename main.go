package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
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

func beginBattle(s *discordgo.Session, channelID string) {
	initiativeTracker.clearEntities()

	// Ask each player for initiative
	s.ChannelMessageSend(channelID, "Let's begin the battle! Please provide your initiative total.")

	for {
		select {
		case <-time.After(30 * time.Second):
			// Timeout after 30 seconds
			s.ChannelMessageSend(channelID, "Initiative input timeout. Battle canceled.")
			return
		default:
			msg, err := s.ChannelMessageSend(channelID, "What is your initiative total?")
			if err != nil {
				fmt.Println("Error sending message:", err)
				return
			}

			response, err := s.ChannelMessageWaitForSingleChoice(channelID, time.Second*30, msg.ID, "Please enter your initiative:", []string{"cancel"})
			if err != nil {
				fmt.Println("Error waiting for message:", err)
				return
			}

			if response == "cancel" {
				s.ChannelMessageSend(channelID, "Battle canceled.")
				return
			}

			initiative, err := strconv.Atoi(response)
			if err != nil {
				s.ChannelMessageSend(channelID, "Invalid initiative input. Please enter a number.")
				continue
			}

			player := &player{
				name:       response,
				initiative: initiative,
				// Other player fields...
			}

			initiativeTracker.addEntity(player)

			s.ChannelMessageSend(channelID, fmt.Sprintf("%s, your initiative is %d.", player.name, player.initiative))

			// Check if all players have provided initiative
			if len(initiativeTracker.entities) == numberOfPlayers {
				initiativeTracker.sortEntities()
				displayTurnOrder(s, channelID, initiativeTracker.entities)
				return
			}
		}
	}
}

func displayTurnOrder(s *discordgo.Session, channelID string, entities []entity) {
	s.ChannelMessageSend(channelID, "Initiative order:")
	for _, e := range entities {
		s.ChannelMessageSend(channelID, fmt.Sprintf("%s: Initiative %d", e.(*player).name, e.getInitiative()))
	}
}