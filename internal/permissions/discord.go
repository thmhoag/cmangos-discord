package permissions

import "github.com/bwmarrin/discordgo"

type DiscordSession interface {
	GuildRoles(guildID string) ([]*discordgo.Role, error)
}