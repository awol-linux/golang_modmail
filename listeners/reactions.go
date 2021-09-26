package listeners

import (
	"context"
	"fmt"
	"log"
	"strconv"

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
	var channel int64
	var reactToMessageID int64
	if message.Author.ID == bot.State.User.ID {
		channel = msg.ChannelID
		reactToMessageID = msg.MessageID
	} else {
		channel = msg.SendtoChannelID
		reactToMessageID = msg.SendtoMessageID
	}
	err = bot.MessageReactionAdd(fmt.Sprint(channel), fmt.Sprint(reactToMessageID), r.Emoji.Name)
	if err != nil {
		log.Printf("Failed to react to message: %v", err)
		return
	}
}
