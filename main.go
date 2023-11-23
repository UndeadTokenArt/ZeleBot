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

func (r *Region) Populate() {
	r.Cities.name = []string{
		"Eldoria",
		"Silverhaven",
		"Frostfall",
		"Emberkeep",
		"Celestial Springs",
		"Stormwatch",
		"Shadowmere",
		"Crystalis",
		"Starlight Citadel",
		"Ironspire",
		"Moonshadow",
		"Verdant Vale",
		"Thunderpeak",
		"Obsidian Reach",
		"Solaris Sanctum",
		"Serpent's Haven",
		"Mistwood",
		"Phoenix Landing",
		"Astral Haven",
		"Stormsgate",
	}
	r.Cities.Factions = []string{
		"Order of the Silver Shield",
		"Crimson Brotherhood",
		"Arcane Consortium",
		"Serpent's Embrace",
		"The Moonlit Syndicate",
		"Ironclad Alliance",
		"Eternal Enclave",
		"Shadowborn Covenant",
		"Celestial Council",
		"Frostborn Legion",
		"Starweaver Guild",
		"Emberforge Clan",
		"Whispering Shadows",
		"Stormbringer Clan",
		"Phoenix Ascendancy",
		"Obsidian Union",
		"Verdant Pact",
		"Mistwalkers",
		"Astral Conclave",
		"Thunderlords",
	}
	r.Cities.Shops = []string{
		"Enchanted Emporium",
		"Mystic Wares",
		"Dragon's Hoard Outfitters",
		"Wizard's Wonders",
		"Elven Elegance",
		"Dwarven Forge Supplies",
		"Celestial Trinkets",
		"Sorcerer's Stockpile",
		"Alchemy Alley",
		"Adventurer's Attire",
		"Treasure Trove Traders",
		"Goblin Goods",
		"Fey Folk Finds",
		"Necromancer's Necessities",
		"Phoenix Feathers and More",
		"Clockwork Curiosities",
		"Bard's Melodies and Merchandise",
		"Rogue's Rarities",
		"Warlock's Weaves",
		"Druid's Delights",
	}
	r.Cities.Civic = []string{
		"Royal Citadel",
		"Mage's Tower",
		"Elven Arboretum",
		"Dwarven Hall of Ancestors",
		"Celestial Observatory",
		"Sorcerer's Library",
		"Council Hall",
		"Grand Bazaar",
		"Enchanted Gardens",
		"Goblin Market",
		"Dragon Roost",
		"Feywild Plaza",
		"Clockwork Workshop",
		"Bardic Amphitheater",
		"Thieves' Den",
		"Warlock's Spire",
		"Druidic Sanctuary",
		"Phoenix Memorial Park",
		"Necropolis",
		"Titan's Arena",
	}
	r.Cities.Resedential = []string{
		"Elven Treehouse",
		"Dwarven Stone Cottage",
		"Mage's Floating Manor",
		"Feywild Enclave",
		"Dragon-Rider's Roost",
		"Gnome Burrow",
		"Hobbit Hillside Hovel",
		"Mermaid's Grotto",
		"Celestial Sky Tower",
		"Steampunk Loft",
		"Wizard's Study Residence",
		"Phoenix Perch Apartments",
		"Goblin Alley Flats",
		"Sorcerer's Spiral Spire",
		"Druidic Tree Canopy Homes",
		"Clockwork Condos",
		"Necromancer's Crypt Apartments",
		"Aquatic Abyss Apartments",
		"Urban Mage Penthouse",
		"Nomadic Wanderer's Wagon",
	}
	r.Cities.Workshops = []string{
		"Forgemaster's Foundry",
		"Enchanting Emporium",
		"Weaver's Workshop",
		"Lumberjack's Lodge",
		"Alchemist's Apothecary",
		"Clockmaker's Studio",
		"Blacksmith's Forge",
		"Golem Workshop",
		"Artificer's Atelier",
		"Glassblower's Gallery",
		"Tinkerer's Tavern",
		"Tailor's Tailoring",
		"Herbalist's Haven",
		"Crafter's Corner",
		"Gemcutter's Gemworks",
		"Scribe's Scriptorium",
		"Sculptor's Studio",
		"Potion Brewery",
		"Magitech Manufacturing",
		"Metalworks Manor",
	}
}
