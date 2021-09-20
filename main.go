package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/awol/golang_modmail/database"
	"github.com/bwmarrin/discordgo"
	_ "github.com/lib/pq"
)

var Token string

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

const (
	host     string = "db"
	port     int    = 5432
	user     string = "khong"
	password string = "khongpass"
	dbName   string = "khong"
	guild    string = "827340883480412173"
)

var sourceName = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)

func main() {
	fmt.Printf("sourcename is: %s\n", sourceName)
	/*
		db, err := sql.Open("postgres", sourceName)
		if err != nil {
			log.Fatalf("failed to open connection to DB: %v", err)
		}

		data := database.AddMessageParams{
			Sender:      1,
			TicketID:    1234,
			MessageText: "tesrst",
			MessageID:   12344,
			ChannelID:   12345,
		}
		database := database.New(db)
		err = database.AddMessage(context.Background(), data)
		if err != nil {
			log.Printf("Failed to add to database: %v,", err)
		}
		res, err := database.GetMessages(
			context.Background(), 12344,
		)
		if err != nil {
			log.Fatalf("failed to retrieve data from db: %v", err)
		}
		db.Close()
		fmt.Printf("database responeded: %v\n", res)
	*/
	runBot()

}

func runBot() {
	fmt.Printf("token is: %v\n", Token)
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatalf("failed to start bot with err: %v", err)
	}
	dg.AddHandler(messageCreate)

	dg.Identify.Intents = discordgo.IntentsDirectMessages
	err = dg.Open()
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
	dg.Close()

}
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	db, err := sql.Open("postgres", sourceName)
	if err != nil {
		log.Fatalf("failed to open connection to DB: %v", err)
	}
	Sender, _ := strconv.ParseInt(m.Author.ID, 10, 64)
	Requester := Sender

	session := database.New(db)
	ticket, err := session.GetOpenTicket(context.Background(), Requester)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Print("No tickets found, Creating new ticket\n")
			ticket, _ = NewTicket(
				session, s, Requester,
			)
		} else {
			log.Printf("failed to get ticket %v", err)
			return
		}
	}
	ForwardMessage(session, s, m, ticket)
}

func ForwardMessage(db *database.Queries, bot *discordgo.Session, message *discordgo.MessageCreate, ticket database.GetOpenTicketRow) (*discordgo.Message, error) {
	Sender, _ := strconv.ParseInt(message.Author.ID, 10, 64)
	MessageText := message.Content
	MessageID, _ := strconv.ParseInt(message.ID, 10, 64)
	ChannelID, _ := strconv.ParseInt(message.ChannelID, 10, 64)

	data := database.AddMessageParams{
		Sender:      Sender,
		TicketID:    int64(ticket.ID),
		MessageText: MessageText,
		MessageID:   MessageID,
		ChannelID:   ChannelID,
	}
	ticketChannel := strconv.Itoa(int(ticket.ChannelID.Int64))
	err := db.AddMessage(context.Background(), data)
	if err != nil {
		log.Printf("Failed to add to database: %v\n", err)
	}
	ret, err := bot.ChannelMessageSend(ticketChannel, MessageText)
	return ret, err
}
func NewTicket(db *database.Queries, bot *discordgo.Session, Requester int64) (database.GetOpenTicketRow, error) {
	err := db.AddTicket(
		context.Background(), Requester,
	)
	if err != nil {
		log.Printf("Failed to add ticket to database: %v\n", err)
	}
	ticket, err := db.GetOpenTicket(context.Background(), Requester)
	if err != nil {
		log.Printf("Failed to Get ticket from DB: %v\n", err)
	}
	channel_data := discordgo.GuildChannelCreateData{
		Name: fmt.Sprintf(
			"ticket_%s", strconv.Itoa(int(ticket.ID)),
		),
		Type:     discordgo.ChannelTypeGuildText,
		ParentID: "889304423455670282",
		Position: 3,
	}
	channel, err := bot.GuildChannelCreateComplex(
		guild,
		channel_data,
	)
	if err != nil {
		log.Printf("Failed to create channel: %v\n", err)
	}
	ChannelID, _ := strconv.ParseInt(channel.ID, 10, 64)
	data := database.InsertChannelParams{
		ChannelID: sql.NullInt64{
			Int64: ChannelID,
			Valid: true,
		},
		ID: ticket.ID,
	}
	db.InsertChannel(context.Background(), data)
	ticket, err = db.GetOpenTicket(context.Background(), Requester)
	fmt.Printf("Ticket is: %v\n", ticket)
	return ticket, err

}
