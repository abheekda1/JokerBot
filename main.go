package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
  "strings"
  "io/ioutil"
  "net/http"
  "time"
  "encoding/json"
  "math/rand"

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
	dg, err := discordgo.New("Bot " + "ODE3NDkwNDgyMDAyNTI2Mjk4.YEKRVw.CFi4lxyf74JcKR2lB-uEAK0tMgk")
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
    resp, err := http.Get("https://official-joke-api.appspot.com/random_joke")
    if err != nil {
      fmt.Println("error with request: ", err)
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
      fmt.Println("error with request: ", err)
    }

    type Joke struct {
      Setup string
      Punchline string
    }

    var joke Joke;

    json.Unmarshal([]byte(string(body)), &joke)

    embed := &discordgo.MessageEmbed{
      Color: 0x004400,
      Description: joke.Setup + "\n||" + joke.Punchline + "||",
      Timestamp: time.Now().Format(time.RFC3339),
      Title: "Let's put a smile on that face...",
    }

    s.ChannelMessageSendEmbed(m.ChannelID, embed)
	}

  if strings.ToLower(m.Content) == "?why so memeless" {
    resp, err := http.Get("https://imgur.com/r/memes/top.json")
    if err != nil {
      fmt.Println("error with request: ", err)
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
      fmt.Println("error with request: ", err)
    }

    //fmt.Println("Body ", string(body))

    type Data struct {
      Hash string
    }

    type Meme struct {
      Data []Data
    }

    var meme Meme

    json.Unmarshal([]byte(body), &meme)
    hash := meme.Data[rand.Intn(60)].Hash

    embed := &discordgo.MessageEmbed{
      Color: 0x004400,
      Timestamp: time.Now().Format(time.RFC3339),
      Title: "Let's put a smile on that face...",
      Image: &discordgo.MessageEmbedImage{
        URL: "https://i.imgur.com/" + hash + ".jpg",
      },
    }

    s.ChannelMessageSendEmbed(m.ChannelID, embed)
	}
}
