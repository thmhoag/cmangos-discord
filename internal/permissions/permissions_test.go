package permissions_test

import (
	"github.com/bwmarrin/discordgo"
	"github.com/segmentio/ksuid"
	perm "github.com/thmhoag/cmangos-discord/internal/permissions"
	"testing"
)

type hasAccessTestArgs struct {
	event *discordgo.MessageCreate
	cmd   string
}

type hasAccessTestCase struct {
	opts    *perm.CmdAuthorityOpts
	args    *hasAccessTestArgs
	want    bool
	wantErr bool
}

func TestHasAccess(t *testing.T) {

	AssertTestCase := func(tc *hasAccessTestCase, t *testing.T) {
		ca := perm.NewCmdAuthority(tc.opts)
		got, err := ca.HasAccess(tc.args.event, tc.args.cmd)
		if (err != nil) != tc.wantErr {
			t.Errorf("HasAccess() error = %v, wantErr %v", err, tc.wantErr)
			return
		}
		if got != tc.want {
			t.Errorf("HasAccess() got = %v, want %v", got, tc.want)
		}
	}

	AssertReturnsDefaultDeny := func(tc *hasAccessTestCase, t *testing.T) {
		t.Run("should return false if denyByDefault is true", func(t *testing.T) {
			tc.opts.DenyByDefault = true
			tc.want = false
			AssertTestCase(tc, t)
		})

		t.Run("should return true if denyByDefault is false", func(t *testing.T) {
			tc.opts.DenyByDefault = false
			tc.want = true
			AssertTestCase(tc, t)
		})
	}

	t.Run("when config is empty", func(t *testing.T) {
		tc := newHasAccessTestCase()
		AssertReturnsDefaultDeny(tc, t)
	})

	t.Run("when command does not exist in config", func(t *testing.T) {
		tc := newHasAccessTestCase()
		tc.opts.Config.Cmds["mycommand"] = &perm.CmdPermission{
			AllowRoles: []string{"myrole"},
			DenyRoles:  []string{"nonehere"},
		}

		tc.args.cmd = "mycommand"

		AssertReturnsDefaultDeny(tc, t)
	})

	t.Run("when event is a DM and command is found in config", func(t *testing.T) {
		tc := newHasAccessTestCase()
		tc.opts.Config.Cmds["mycommand"] = &perm.CmdPermission{
			AllowRoles: []string{"myrole"},
			DenyRoles:  []string{"nonehere"},
		}

		tc.args.cmd = "mycommand"
		tc.args.event.GuildID = ""
		AssertReturnsDefaultDeny(tc, t)
	})

	t.Run("when event author is found in denyUsers", func(t *testing.T) {
		tc := newHasAccessTestCase()
		testCmdName := "testcmd"
		tc.args.cmd = testCmdName

		denyUserID := ksuid.New().String()
		tc.opts.Config.Cmds[testCmdName] = &perm.CmdPermission{
			DenyUsers: []string{denyUserID},
		}

		tc.args.event.Author.ID = denyUserID
		tc.want = false

		t.Run("access should be denied even if author is also in allowUsers", func(t *testing.T) {
			tc.opts.Config.Cmds[testCmdName].AllowUsers = []string{denyUserID}
			AssertTestCase(tc, t)
		})

		t.Run("access should be denied even if author is also in an allowed role", func(t *testing.T) {
			roleID := ksuid.New().String()
			roleName := "thisismytestrole"
			tc.opts.Session = newSessionWithRoles(&discordgo.Role{
				ID: roleID,
				Name: roleName,
			})

			tc.opts.Config.Cmds[testCmdName].AllowRoles = []string{roleName}
			tc.args.event.Member.Roles = []string{roleID}
			AssertTestCase(tc, t)
		})
	})

	t.Run("when event author is found in allowUsers and defaultDeny is true, access should be allowed", func(t *testing.T) {
		tc := newHasAccessTestCase()
		testCmdName := "testcmd"
		tc.args.cmd = testCmdName

		allowUserID := ksuid.New().String()
		tc.opts.Config.Cmds[testCmdName] = &perm.CmdPermission{
			AllowUsers: []string{allowUserID},
		}

		tc.opts.DenyByDefault = true
		tc.args.event.Author.ID = allowUserID
		tc.want = true

		AssertTestCase(tc, t)
	})

	t.Run("when event author is found in denyUsers and defaultDeny is false, access should be denied", func(t *testing.T) {
		tc := newHasAccessTestCase()
		testCmdName := "testcmd"
		tc.args.cmd = testCmdName

		denyUserID := ksuid.New().String()
		tc.opts.Config.Cmds[testCmdName] = &perm.CmdPermission{
			DenyUsers: []string{denyUserID},
		}

		tc.opts.DenyByDefault = false
		tc.args.event.Author.ID = denyUserID
		tc.want = false

		AssertTestCase(tc, t)
	})

	t.Run("when event author has an allowed role and defaultDeny is true, access should be allowed", func(t *testing.T) {
		tc := newHasAccessTestCase()
		testCmdName := "testcmd"
		tc.args.cmd = testCmdName

		roleID := ksuid.New().String()
		roleName := "thisismytestrole"
		tc.opts.Session = newSessionWithRoles(&discordgo.Role{
			ID: roleID,
			Name: roleName,
		})

		tc.opts.Config.Cmds[testCmdName] = &perm.CmdPermission{
			AllowRoles: []string{roleName},
		}

		tc.args.event.Member.Roles = []string{roleID}
		tc.opts.DenyByDefault = true
		tc.want = true

		AssertTestCase(tc, t)
	})

	t.Run("when event author has a denied role and defaultDeny is false, access should be denied", func(t *testing.T) {
		tc := newHasAccessTestCase()
		testCmdName := "testcmd"
		tc.args.cmd = testCmdName

		roleID := ksuid.New().String()
		roleName := "thisismytestrole"
		tc.opts.Session = newSessionWithRoles(&discordgo.Role{
			ID: roleID,
			Name: roleName,
		})

		tc.opts.Config.Cmds[testCmdName] = &perm.CmdPermission{
			DenyRoles: []string{roleName},
		}

		tc.args.event.Member.Roles = []string{roleID}
		tc.opts.DenyByDefault = false
		tc.want = false

		AssertTestCase(tc, t)
	})
}

func newTestConfig() *perm.Config {
	return &perm.Config{
		Cmds: map[string]*perm.CmdPermission{},
	}
}

func newTestAuthorityOpts() *perm.CmdAuthorityOpts {
	return &perm.CmdAuthorityOpts{
		DenyByDefault: false,
		CaseSensitive: false,
		Config:        newTestConfig(),
		Session: &mockDiscordSession{
			execGuildRoles: func(guildID string) ([]*discordgo.Role, error) {
				return []*discordgo.Role{}, nil
			},
		},
	}
}

func newSessionWithRoles(roles ...*discordgo.Role) perm.DiscordSession {
	mroles := roles
	return &mockDiscordSession{
		execGuildRoles: func(guildID string) ([]*discordgo.Role, error) {
			return mroles, nil
		},
	}
}

func newTestMessageCreateEvent() *discordgo.MessageCreate {
	guildID := ksuid.New().String()
	user := &discordgo.User{
		ID:            ksuid.New().String(),
		Username:      ksuid.New().String(),
		Discriminator: "0000",
		Bot:           false,
	}

	return &discordgo.MessageCreate{
		Message: &discordgo.Message{
			ID:        ksuid.New().String(),
			ChannelID: ksuid.New().String(),
			GuildID:   guildID,
			Content:   "This is a message",
			Author:    user,
			Member: &discordgo.Member{
				GuildID: guildID,
				User:    user,
				Roles: []string{
					ksuid.New().String(),
					ksuid.New().String(),
				},
			},
		},
	}
}

func newHasAccessTestCase() *hasAccessTestCase {
	opts := newTestAuthorityOpts()
	opts.Config = newTestConfig()

	event := newTestMessageCreateEvent()
	return &hasAccessTestCase{
		opts: opts,
		args: &hasAccessTestArgs{
			event: event,
			cmd:   "testcmd",
		},
		want:    false,
		wantErr: false,
	}
}
