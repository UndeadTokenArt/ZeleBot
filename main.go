package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var Token string
var numberOfPlayers = 2

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

	if strings.ToLower(m.Content) == "!beginbattle" {
		beginBattle(s, m.ChannelID, m.Reference().MessageID, &InitiativeTracker{})
	}
}

type entity interface {
	attack(target entity, dmg int)
	getCurrentDmgDone() int
	getInitiative() int
}

func beginBattle(s *discordgo.Session, channelID string, messageID string, it *InitiativeTracker) {
	it.ClearEntities()

	s.ChannelMessageSend(channelID, "Let's begin the battle! Please provide your initiative total. Type 'cancel' at any time to cancel the battle.")

	// Create a map to store user input
	userInput := make(map[string]int)

	// Function to check if all players have provided initiative
	allPlayersReady := func() bool {
		return len(userInput) == numberOfPlayers
	}

	// Function to display the current initiative status
	displayInitiativeStatus := func() {
		var statusMessage strings.Builder
		statusMessage.WriteString("Initiative status:\n")
		for player, initiative := range userInput {
			statusMessage.WriteString(fmt.Sprintf("%s: Initiative %d\n", player, initiative))
		}
		s.ChannelMessageSend(channelID, statusMessage.String())
	}

	// Function to cancel the battle
	cancelBattle := func() {
		s.ChannelMessageSend(channelID, "Battle canceled.")
	}

	for {
		// Collect messages
		msg, err := s.State.Message(channelID, messageID)
		if err != nil {
			fmt.Println("Error creating message:", err)
			cancelBattle()
			return
		}

		response, err := s.State.Message(channelID, msg.ID)
		if err != nil {
			fmt.Println("Error waiting for message:", err)
			cancelBattle()
			return
		}

		if response.Author.ID == s.State.User.ID {
			// Ignore messages from the bot itself
			continue
		}

		// Check for cancel command
		if strings.ToLower(response.Content) == "cancel" {
			cancelBattle()
			return
		}

		// Parse user input as initiative
		initiative, err := strconv.Atoi(response.Content)
		if err != nil {
			s.ChannelMessageSend(channelID, "Invalid initiative input. Please enter a number.")
			continue
		}

		// Store user input
		userInput[response.Author.Username] = initiative

		// Display current initiative status
		displayInitiativeStatus()

		// Check if all players have provided initiative
		if allPlayersReady() {
			it.SortEntities()
			displayTurnOrder(s, channelID, it.entities)
			return
		}
	}
}

func displayTurnOrder(s *discordgo.Session, channelID string, entities []entity) {
	s.ChannelMessageSend(channelID, "Initiative order:")
	for _, e := range entities {
		s.ChannelMessageSend(channelID, fmt.Sprintf("%s: Initiative %d", e.(*player).name, e.getInitiative()))
	}
}
