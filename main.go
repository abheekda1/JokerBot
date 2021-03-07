package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
)

func init() {

}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + os.Args[1])
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
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
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

	if strings.ToLower(m.Content) == "?why so serious" {
		resp, err := http.Get("http://localhost:3587/jokes/random/general")
		if err != nil {
			fmt.Println("error with request: ", err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("error with request: ", err)
		}

		type Joke struct {
			Setup     string
			Punchline string
		}

		var joke Joke

		json.Unmarshal([]byte(string(body)), &joke)

		embed := &discordgo.MessageEmbed{
			Color:       0x004400,
			Description: joke.Setup + "\n||" + joke.Punchline + "||",
			Timestamp:   time.Now().Format(time.RFC3339),
			Title:       "Let's put a smile on that face...",
		}

		s.ChannelMessageSendEmbed(m.ChannelID, embed)
	}

	if strings.ToLower(m.Content) == "?why so scientific" {
		resp, err := http.Get("http://localhost:3587/jokes/random/science")
		if err != nil {
			fmt.Println("error with request: ", err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("error with request: ", err)
		}

		type Joke struct {
			Subject   string
			Title     string
			Setup     string
			Punchline string
			Oneliner  string
		}

		var joke Joke

		json.Unmarshal([]byte(string(body)), &joke)

		if joke.Oneliner != "" {
			embed := &discordgo.MessageEmbed{
				Color:       0x004400,
				Description: joke.Oneliner,
				Timestamp:   time.Now().Format(time.RFC3339),
				Title:       "Let's put a smile on that face...",
			}
			s.ChannelMessageSendEmbed(m.ChannelID, embed)
		} else {
			embed := &discordgo.MessageEmbed{
				Color:       0x004400,
				Description: joke.Setup + "\n||" + joke.Punchline + "||",
				Timestamp:   time.Now().Format(time.RFC3339),
				Title:       "Let's put a smile on that face...",
			}
			s.ChannelMessageSendEmbed(m.ChannelID, embed)
		}
	}

	if strings.ToLower(m.Content) == "?why so helpful" {
		embed := &discordgo.MessageEmbed{
			Color:       0x004400,
			Description: "`?why so serious`: get a random joke\n`?why so scientific`: get a science joke",
			Title:       "Commands:",
		}
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
	}

	if strings.ToLower(m.Content) == "?why so statistical" {
		embed := &discordgo.MessageEmbed{
			Color:       0x004400,
			Description: fmt.Sprintf("Servers: %v", len(s.State.Guilds)),
			Title:       "Stats",
		}
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
	}
}
