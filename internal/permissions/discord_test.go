package permissions_test

import "github.com/bwmarrin/discordgo"

type mockGuildRolesFunc func(guildID string) ([]*discordgo.Role, error)
type mockDiscordSession struct {
	execGuildRoles mockGuildRolesFunc
}

func (mds *mockDiscordSession) GuildRoles(guildID string) ([]*discordgo.Role, error) {
	return mds.execGuildRoles(guildID)
}