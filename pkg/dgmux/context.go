package dgmux

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

var (
	dmChannelsCache map[string]*discordgo.Channel
)

// Ctx represents the context for a command event
type Ctx struct {
	Session       *discordgo.Session
	Event         *discordgo.MessageCreate
	Arguments     *Arguments
	CustomObjects *ObjectsMap
	Router        *Router
	Command       *Command
}

// ExecutionHandler represents a handler for a context execution
type ExecutionHandler func(*Ctx)

// ResponseChannelID returns the correct channel to use for reply/response
func (ctx *Ctx) ResponseChannelID() string {

	if ctx.Command == nil || !ctx.Command.DmOnly {
		return ctx.Event.ChannelID
	}

	// this caching method is effectively a memory leak
	// and should be changed to a more robust solution later
	dmChan := dmChannelsCache[ctx.Msg().Author.ID]
	if dmChan == nil {
		userDmChan, err := ctx.Session.UserChannelCreate(ctx.Msg().Author.ID)
		if err != nil {
			fmt.Errorf("error starting dm with user %s\nerror: %s\n", ctx.Msg().Author.ID, err)
			return ctx.Event.ChannelID
		}

		dmChan = userDmChan
	}

	return dmChan.ID
}

// Msg returns the message associated with the discord event
func (ctx *Ctx) Msg() *discordgo.Message {
	if ctx.Event == nil {
		return nil
	}

	return ctx.Event.Message
}

// Reply responds with the given text message
func (ctx *Ctx) Reply(msg string) error {
	channelID := ctx.ResponseChannelID()
	_, err := ctx.Session.ChannelMessageSend(channelID, msg)
	return err
}

// ReplyEmbed responds with the given text and embed message
func (ctx *Ctx) ReplyEmbed(msg string, embed *discordgo.MessageEmbed) error {
	channelID := ctx.ResponseChannelID()
	_, err := ctx.Session.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Content: msg,
		Embed:   embed,
	})
	return err
}

// Reply responds with the given text message
func (ctx *Ctx) ReplyDm(msg string) error {
	dmChan := dmChannelsCache[ctx.Msg().Author.ID]
	if dmChan == nil {
		userDmChan, err := ctx.Session.UserChannelCreate(ctx.Msg().Author.ID)
		if err != nil {
			return fmt.Errorf("error starting dm with user %s\nerror: %s\n", ctx.Msg().Author.ID, err)
		}

		dmChan = userDmChan
	}

	_, err := ctx.Session.ChannelMessageSend(dmChan.ID, msg)
	if err != nil {
		// clear cached channel in case the cache is invalid
		delete(dmChannelsCache, ctx.Msg().Author.ID)
	}

	return err
}

// Reply responds with the given text message
func (ctx *Ctx) ReplyDmEmbed(msg string, embed *discordgo.MessageEmbed) error {
	dmChan := dmChannelsCache[ctx.Msg().Author.ID]
	if dmChan == nil {
		userDmChan, err := ctx.Session.UserChannelCreate(ctx.Msg().Author.ID)
		if err != nil {
			return fmt.Errorf("error starting dm with user %s\nerror: %s\n", ctx.Msg().Author.ID, err)
		}

		dmChan = userDmChan
	}

	_, err := ctx.Session.ChannelMessageSendComplex(dmChan.ID, &discordgo.MessageSend{
		Content: msg,
		Embed:   embed,
	})

	if err != nil {
		// clear cached channel in case the cache is invalid
		delete(dmChannelsCache, ctx.Msg().Author.ID)
	}

	return err
}