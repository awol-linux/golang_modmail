package listeners

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/awol/golang_modmail/database"
	"github.com/bwmarrin/discordgo"
)

const (
	guild string = "827340883480412173"
)

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	db, err := database.GetDB()
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
	channel, err := bot.Channel(message.ChannelID)
	if err != nil {
		log.Printf("Failed to get Channel: %v\n", err)
		return nil, err
	}
	var sendto *discordgo.Channel
	if channel.Type == discordgo.ChannelTypeDM {
		sendto, err = bot.Channel(strconv.Itoa(int(ticket.ChannelID.Int64)))
		if err != nil {
			log.Printf("Error finding channel: %v\n", err)
		}
	} else {
		requester := strconv.Itoa(int(ticket.Requester))
		sendto, err = bot.UserChannelCreate(requester)
		if err != nil {
			log.Printf("Error creating DM with %s: %v", requester, err)
		}
	}
	data := database.AddMessageParams{
		Sender:      Sender,
		TicketID:    int64(ticket.ID),
		MessageText: MessageText,
		MessageID:   MessageID,
		ChannelID:   ChannelID,
	}
	err = db.AddMessage(context.Background(), data)
	if err != nil {
		log.Printf("Failed to add to database: %v\n", err)
	}
	ret, err := bot.ChannelMessageSend(sendto.ID, MessageText)
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
