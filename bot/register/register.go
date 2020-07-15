package register

import (
	"fmt"
	"github.com/thmhoag/cmangos-discord/pkg/cmangos"
	"github.com/thmhoag/cmangos-discord/pkg/dgmux"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func NewRegisterCmd(ctx Ctx) *dgmux.Command {

	client := ctx.MangosClient()

	cmd := &dgmux.Command{
		Name: 			"register",
		Description: 	"Registers a new account for your discord user",
		IgnoreCase: 	true,
		DmOnly: 		true,
		Handler: 		func(ctx *dgmux.Ctx) {
			if ctx.Event.GuildID == "" {
				ctx.ReplyDm("You may not register in a DM.")
				return
			}

			allowedRole := os.Getenv("REGISTER_CMD_ROLE")
			if allowedRole != "" {
				if ctx.Event.Member == nil || ctx.Event.Member.Roles == nil {
					ctx.ReplyDm("Unable to register accounts at this time.")
					log.Println("Unable to get discord member roles.")
					return
				}

				guildRoles, _ := ctx.Session.GuildRoles(ctx.Event.GuildID)
				guildRoleId := ""
				for _,role := range guildRoles {
					if strings.ToLower(role.Name) == strings.ToLower(allowedRole) {
						guildRoleId = role.ID
					}
				}

				isAllowed := false
				for _,role := range ctx.Event.Member.Roles {
					if guildRoleId == role {
						isAllowed = true
						break
					}
				}

				if !isAllowed {
					ctx.ReplyDm(fmt.Sprintf("You are not permitted to register. You must be a member of the %s role.", allowedRole))
					return
				}
			}

			acctName := ctx.Msg().Author.Username + ctx.Msg().Author.Discriminator
			password := generatePassword()

			resp, err := client.SendExecCmd(&cmangos.ExecCmdRequest{
				Command: fmt.Sprintf("account create %s %s", acctName, password),
			})

			if err != nil {
				// if error is nil, check if the account already exists
				_, err2 := client.SendExecCmd(&cmangos.ExecCmdRequest{
					Command: fmt.Sprintf("account set addon %s 0", acctName),
				})

				if err2 != nil {
					log.Printf("error executing register command: %s", err)
					ctx.ReplyDm("Unable to register accounts at this time.")
				}

				ctx.ReplyDm(fmt.Sprintf("Account has already been registered. Your username is `%s`.", acctName))
				return
			}

			if resp.Body.Fault != nil {
				log.Printf("error from cmangos server: %s", resp.Body.Fault.Faultstring)
				ctx.ReplyDm(fmt.Sprintf("error: %s", resp.Body.Fault.Faultstring))
				return
			}

			result := resp.Body.ExecCmdResponseText.Result
			if strings.Contains(strings.ToLower(result), "already exist") {
				ctx.ReplyDm("Your account has already been registered!")
			} else if !strings.Contains(strings.ToLower(result), "account created") {
				log.Printf("account creation - unexpected result: %s\n", result)
			}

			log.Printf("account created: %s\n", acctName)
			ctx.ReplyDm(generateAccountCreatedReply(ctx.Msg().Author.Username, acctName, password))
		},
	}

	return cmd
}

func generatePassword() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")

	length := 16 // this is the max length the client will allow
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}

	str := b.String()
	return str
}

func generateAccountCreatedReply(username string, acctName string, password string) string {
	return fmt.Sprintf(`Thanks %s, your account was created successfully!

Your credentials are as follows:
Account: ` + "`%s`" + `
Password: ` + "`%s`" + `

After logging in for the first time, please change your password by typing the following command in your chatbox:
` + "`.account password $old_password $new_password $new_password`", username, acctName, password)
}