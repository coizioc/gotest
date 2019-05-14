package main

import "github.com/bwmarrin/discordgo"

// Ctx contains information related to a command's invocation context.
type Ctx struct {
	Session *discordgo.Session
	Message *discordgo.MessageCreate
	Guild   *discordgo.Guild
	Channel *discordgo.Channel
	Member  *discordgo.Member
}

// GetCtx creates Ctx with a session and message as parameters.
func GetCtx(s *discordgo.Session, m *discordgo.MessageCreate) (*Ctx, error) {
	var err error
	ctx := new(Ctx)
	ctx.Session = s
	ctx.Message = m
	ctx.Channel, err = s.Channel(m.ChannelID)
	if err != nil {
		return nil, err
	}

	ctx.Guild, err = s.Guild(m.GuildID)
	if err != nil {
		return nil, err
	}

	for _, member := range ctx.Guild.Members {
		if member.User.ID == m.Author.ID {
			ctx.Member = member
			break
		}
	}

	return ctx, nil
}

// Send sends a message to the channel in the Ctx object.
func (ctx *Ctx) Send(content string) (*discordgo.Message, error) {
	message, err := ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, content)
	return message, err
}

// Edit edits a message in the channel in the Ctx object.
func (ctx *Ctx) Edit(mid string, content string) (*discordgo.Message, error) {
	message, err := ctx.Session.ChannelMessageEdit(ctx.Message.ChannelID, mid, content)
	return message, err
}
