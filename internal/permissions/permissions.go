package permissions

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

type cmdAuthority struct {
	denyByDefault bool
	caseSensitive bool
	config        *Config
	session       DiscordSession
}

type CmdAuthority interface {
	HasAccess(*discordgo.MessageCreate, string) (bool, error)
}

type CmdAuthorityOpts struct {
	DenyByDefault bool
	CaseSensitive bool
	Config        *Config
	Session       DiscordSession
}

func NewCmdAuthority(opts *CmdAuthorityOpts) CmdAuthority {
	return &cmdAuthority{
		denyByDefault: opts.DenyByDefault,
		caseSensitive: opts.CaseSensitive,
		config:        opts.Config,
		session:       opts.Session,
	}
}

func (ca *cmdAuthority) HasAccess(event *discordgo.MessageCreate, cmd string) (bool, error) {

	userID := event.Author.ID
	cmdPerms := ca.config.Cmds[cmd]
	if cmdPerms == nil {
		return !ca.denyByDefault, nil

	}

	if stringArrayContains(cmdPerms.DenyUsers, userID, ca.caseSensitive) {
		return false, nil
	}

	if stringArrayContains(cmdPerms.AllowUsers, userID, ca.caseSensitive) {
		return true, nil
	}

	// check if it's a dm
	if event.GuildID == "" {
		return !ca.denyByDefault, nil
	}

	guildRoles, err := ca.session.GuildRoles(event.GuildID)
	if err != nil {
		return false, err
	}

	denyRoleIDs := map[string]bool{}
	for _, role := range cmdPerms.DenyRoles {
		denyRole := ca.guildRoleByName(guildRoles, role)
		if denyRole == nil {
			continue
		}

		denyRoleIDs[denyRole.ID] = true
	}

	for _, roleID := range event.Member.Roles {
		if denyRoleIDs[roleID] {
			return false, nil
		}
	}

	allowRoleIDs := map[string]bool{}
	for _, role := range cmdPerms.AllowRoles {
		allowRole := ca.guildRoleByName(guildRoles, role)
		if allowRole == nil {
			continue
		}

		allowRoleIDs[allowRole.ID] = true
	}

	for _, roleID := range event.Member.Roles {
		if allowRoleIDs[roleID] {
			return true, nil
		}
	}

	return !ca.denyByDefault, nil
}

func (ca *cmdAuthority) guildRolesByName(guildID string) (map[string]*discordgo.Role, error) {
	if guildID == "" {
		return nil, nil
	}

	guildRoles, err := ca.session.GuildRoles(guildID)
	if err != nil {
		return nil, err
	}

	roleMap := map[string]*discordgo.Role{}
	for i := range guildRoles {
		roleToMap := guildRoles[i]
		roleMap[roleToMap.ID] = roleToMap
	}

	return roleMap, nil
}

func (ca *cmdAuthority) guildRoleByName(guildRoles []*discordgo.Role, desired string) *discordgo.Role {
	if desired == "" || len(guildRoles) < 1 {
		return nil
	}

	for i := range guildRoles {
		role := guildRoles[i]
		if ca.caseSensitive {
			if role.Name == desired {
				return role
			}

			continue
		}

		if strings.EqualFold(role.Name, desired) {
			return role
		}
	}

	return nil
}

func stringArrayContains(strarray []string, desired string, caseSensitive bool) bool {
	for _, str := range strarray {
		if caseSensitive {
			if str == desired {
				return true
			}

			continue
		}

		if strings.EqualFold(str, desired) {
			return true
		}
	}

	return false
}
