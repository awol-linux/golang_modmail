package listeners

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/awol/golang_modmail/database"
	"github.com/bwmarrin/discordgo"
)

func (l *Listeners) MessageReact(bot *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.UserID == bot.State.User.ID {
		return
	}
	msgid, _ := strconv.ParseInt(r.MessageID, 10, 64)
	msg, err := l.DB.GetMessage(context.Background(), msgid)
	if err != nil {
		log.Printf("Failed to get message from DB: %v", err)
		return
	}
	message, err := bot.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		log.Printf("Error finding message: %v\n", err)
		return
	}
	channel, sendto, _ := l.getChannels(bot, msg, message.ID)
	err = bot.MessageReactionRemove(channel, sendto, r.Emoji.Name, r.Emoji.Name)
	if err != nil {
		log.Printf("Failed to react to message: %v", err)
		log.Printf("Messge is: %v, channel is %s", channel, sendto)

		return
	}
}

func (l *Listeners) UnMessageReact(bot *discordgo.Session, r *discordgo.MessageReactionRemove) {
	if r.UserID == bot.State.User.ID {
		return
	}
	msgid, _ := strconv.ParseInt(r.MessageID, 10, 64)
	msg, err := l.DB.GetMessage(context.Background(), msgid)
	if err != nil {
		log.Printf("Failed to get message from DB: %v", err)
		return
	}
	message, err := bot.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		log.Printf("Error finding message: %v\n", err)
		return
	}
	channel, sendto, _ := l.getChannels(bot, msg, message.ID)
	err = bot.MessageReactionRemove(channel, sendto, r.Emoji.Name, bot.State.User.ID)
	if err != nil {
		log.Printf("Failed to unreact to message: %v", err)
		return
	}
}
func (l *Listeners) getChannels(bot *discordgo.Session, msg database.GetMessageRow, messageId string) (string, string, error) {
	var channel int64
	var reactToMessageID int64
	msgId, _ := strconv.ParseInt(messageId, 10, 64)
	if msg.MessageID != msgId {
		channel = msg.ChannelID
		reactToMessageID = msg.MessageID
	} else {
		channel = msg.SendtoChannelID
		reactToMessageID = msg.SendtoMessageID
	}
	return fmt.Sprint(channel), fmt.Sprint(reactToMessageID), nil
}
