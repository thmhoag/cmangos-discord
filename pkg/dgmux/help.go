package dgmux

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

// RegisterDefaultHelpCommand registers the default help command
func (router *Router) RegisterDefaultHelpCommand(session *discordgo.Session, rateLimiter RateLimiter) {
	// Initialize the helo messages storage
	router.InitializeStorage("dgc_helpMessages")

	// Register the default help command
	router.RegisterCmd(&Command{
		Name:        "help",
		Description: "Lists all the available commands or displays some information about a specific command",
		Usage:       "help [command name]",
		Example:     "help yourCommand",
		IgnoreCase:  true,
		DmOnly: 	 true,
		RateLimiter: rateLimiter,
		Handler:     generalHelpCommand,
	})
}

// generalHelpCommand handles the general help command
func generalHelpCommand(ctx *Ctx) {
	// Check if the user provided an argument
	if ctx.Arguments.Amount() > 0 {
		specificHelpCommand(ctx)
		return
	}

	// Define useful variables
	channelID := ctx.ResponseChannelID()

	// Send the general help embed
	embed := renderDefaultGeneralHelpEmbed(ctx.Router)
	_, err := ctx.Session.ChannelMessageSendEmbed(channelID, embed)
	if err != nil {
		log.Printf("unable to send channel message: %s\n", err)
	}
}

// specificHelpCommand handles the specific help command
func specificHelpCommand(ctx *Ctx) {
	// Define the command names
	commandNames := strings.Split(ctx.Arguments.Raw(), " ")

	// Define the command
	var command *Command
	for index, commandName := range commandNames {
		if index == 0 {
			command = ctx.Router.GetCmd(commandName)
			continue
		}
		command = command.GetSubCmd(commandName)
	}

	// Send the help embed
	ctx.Session.ChannelMessageSendEmbed(ctx.ResponseChannelID(), renderDefaultSpecificHelpEmbed(ctx, command))
}

// renderDefaultGeneralHelpEmbed renders the general help embed on the given page
func renderDefaultGeneralHelpEmbed(router *Router) *discordgo.MessageEmbed {
	// Define useful variables
	commands := router.Commands
	prefix := router.Prefixes[0]

	// Prepare the fields for the embed
	var fields []*discordgo.MessageEmbedField
	for _, command := range commands {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   prefix + command.Name,
			Value:  "`" + command.Description + "`",
			Inline: false,
		})
	}

	// Return the embed and the new page
	return &discordgo.MessageEmbed{
		Type:        "rich",
		Title:       "Commands",
		Description: "Type `" + prefix + "help <command name>` to find out more about a specific command.",
		Color:       0xffff00,
		Fields:      fields,
	}
}

// renderDefaultSpecificHelpEmbed renders the specific help embed of the given command
func renderDefaultSpecificHelpEmbed(ctx *Ctx, command *Command) *discordgo.MessageEmbed {
	// Define useful variables
	prefix := ctx.Router.Prefixes[0]

	// Check if the command is invalid
	if command == nil {
		return &discordgo.MessageEmbed{
			Type:      "rich",
			Title:     "Error",
			Color:     0xff0000,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Message",
					Value:  "```The given command doesn't exist. Type `" + prefix + "help` for a list of available commands.```",
					Inline: false,
				},
			},
		}
	}

	var fields []*discordgo.MessageEmbedField
	if command.Usage != "" {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Usage",
			Value:  command.Usage,
			Inline: false,
		})
	}

	if command.Example != "" {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Example",
			Value:  command.Example,
			Inline: false,
		})
	}

	// Define the aliases string
	if len(command.Aliases) > 0 {
		aliases := "`" + strings.Join(command.Aliases, "`, `") + "`"
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Aliases",
			Value:  aliases,
			Inline: false,
		})
	}

	// Define the sub commands string
	if len(command.SubCommands) > 0 {
		subCommandNames := make([]string, len(command.SubCommands))
		for index, subCommand := range command.SubCommands {
			subCommandNames[index] = subCommand.Name
		}
		subCommands := "`" + strings.Join(subCommandNames, "`, `") + "`"
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Sub Commands",
			Value:  subCommands,
			Inline: false,
		})
	}

	result := &discordgo.MessageEmbed{
		Type:        "rich",
		Title:       prefix + command.Name,
		Color:       0xffff00,
		Fields:      fields,
	}

	if command.Description != "" {
		result.Description = command.Description
	}

	return result
}
