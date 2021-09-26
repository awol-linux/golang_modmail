package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/awol/golang_modmail/database"
	"github.com/awol/golang_modmail/listeners"
	"github.com/bwmarrin/discordgo"
	_ "github.com/lib/pq"
)

var Token string

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	bot, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatalf("failed to start bot with err: %v", err)
	}
	runBot(bot)
}

func runBot(bot *discordgo.Session) {
	fmt.Printf("token is: %v\n", Token)
	bot.Identify.Intents = discordgo.IntentsAll
	db, err := database.GetDB()
	if err != nil {
		log.Printf("failed to open connection to DB: %v", err)
	}
	session := database.New(db)
	defer db.Close()
	messages := listeners.Listeners{DB: session}
	bot.AddHandler(messages.MessageCreate)
	bot.AddHandler(messages.MessageReact)

	err = bot.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	bot.Close()
}
