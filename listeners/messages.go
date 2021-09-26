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

type Listeners struct {
	DB *database.Queries
}

func (l *Listeners) MessageDelete(bot *discordgo.Session, message *discordgo.MessageDelete) {
	messageId, err := strconv.ParseInt(message.ID, 10, 64)
	if err != nil {
		log.Printf("Failed convert messageID to bot: %v", err)
		return
	}
	msg, err := l.DB.GetMessage(context.Background(), messageId)
	if err != nil {
		log.Printf("Failed to get message from DB: %v", err)
		return
	}
	if messageId == msg.SendtoMessageID {
		return
	}
	deleteData := database.DeleteMessageParams{
		Deleted:   true,
		MessageID: messageId,
	}
	err = l.DB.DeleteMessage(context.Background(), deleteData)
	if err != nil {
		log.Printf("Failed to delete message from DB: %v", err)
		return
	}
	channel, sendto, _ := l.getChannels(bot, msg, message.ID)

	err = bot.ChannelMessageDelete(channel, sendto)
	if err != nil {
		log.Printf("Failed to delete message: %v", err)
		return
	}
}

func (l *Listeners) MessageCreate(s *discordgo.Session, message *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if message.Author.ID == s.State.User.ID {
		return
	}
	Sender, _ := strconv.ParseInt(message.Author.ID, 10, 64)
	Requester := Sender

	ticket, err := l.DB.GetOpenTicket(context.Background(), Requester)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Print("No tickets found, Creating new ticket\n")
			ticket, err = l.newTicket(
				s, Requester,
			)
			if err != nil {
				log.Printf("Failed to create ticket with err %v", err)
				return
			}
		} else {
			log.Printf("failed to get ticket %v", err)
			return
		}
	}
	forward, err := l.forwardMessage(s, message, ticket)
	if err != nil {
		log.Printf("Failed to forward message: %v", err)
		return
	}
	ogmsgid, _ := strconv.ParseInt(message.ID, 10, 64)
	frwrdid, _ := strconv.ParseInt(forward.ID, 10, 64)
	frwrdchnlid, _ := strconv.ParseInt(forward.ChannelID, 10, 64)
	data := database.InsertForwardParams{
		SendtoMessageID: frwrdid,
		SendtoChannelID: frwrdchnlid,
	}
	err = l.DB.InsertForward(context.Background(), data)
	if err != nil {
		log.Printf("Failed to insert forwarded message: %v", err)
		return
	}
	linkForward := database.LinkForwardParams{
		Forwarded: sql.NullInt64{
			Int64: frwrdid,
			Valid: true,
		},
		MessageID: ogmsgid,
	}
	err = l.DB.LinkForward(context.Background(), linkForward)
	if err != nil {
		log.Printf("Failed to link forwarded message: %v", err)
	}
}

func (l *Listeners) forwardMessage(bot *discordgo.Session, message *discordgo.MessageCreate, ticket database.GetOpenTicketRow) (*discordgo.Message, error) {
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
		sendto, err = bot.Channel(strconv.Itoa(int(ticket.TicketChannelID.Int64)))
		if err != nil {
			log.Printf("Error finding channel: %v\n", err)
			return nil, err
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
	err = l.DB.AddMessage(context.Background(), data)
	if err != nil {
		log.Printf("Failed to add to database: %v\n", err)
		return nil, err
	}
	ret, err := bot.ChannelMessageSend(sendto.ID, MessageText)
	return ret, err
}

func (l *Listeners) newTicket(bot *discordgo.Session, Requester int64) (database.GetOpenTicketRow, error) {
	err := l.DB.AddTicket(
		context.Background(), Requester,
	)
	if err != nil {
		log.Printf("Failed to add ticket to database: %v\n", err)
	}
	ticket, err := l.DB.GetOpenTicket(context.Background(), Requester)
	if err != nil {
		log.Printf("Failed to Get ticket from DB: %v\n", err)
	}
	log.Printf("ticket is %v", ticket.TicketChannelID.Int64)
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
		return database.GetOpenTicketRow{}, err
	}
	ChannelID, _ := strconv.ParseInt(channel.ID, 10, 64)
	data := database.InsertChannelParams{
		TicketChannelID: sql.NullInt64{
			Int64: ChannelID,
			Valid: true,
		},
		ID: ticket.ID,
	}
	l.DB.InsertChannel(context.Background(), data)
	ticket, err = l.DB.GetOpenTicket(context.Background(), Requester)
	fmt.Printf("Ticket is: %v\n", ticket)
	return ticket, err
}
