package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
	"github.com/thmhoag/cmangos-discord/bot/register"
	"github.com/thmhoag/cmangos-discord/pkg/cmangos"
	"github.com/thmhoag/cmangos-discord/pkg/dgmux"
	"os"
	"strings"
)

func Execute() {
	// Open a simple Discord session
	token := os.Getenv("TOKEN")
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	err = session.Open()
	if err != nil {
		panic(err)
	}

	ctx := &globalCtx{
		mangosClient: createMangosClient(),
	}

	router := dgmux.Create(&dgmux.Router{
		Prefixes: []string{"!"},
		IgnorePrefixCase: true,
		BotsAllowed: false,
		Commands: []*dgmux.Command{},
		Middlewares: []dgmux.Middleware{},
		PingHandler: func(ctx *dgmux.Ctx) {
			ctx.Reply("Pong!")
		},
	})

	router.RegisterDefaultHelpCommand(session, nil)

	router.RegisterCmd(register.NewRegisterCmd(ctx))

	router.Initialize(session)
}

func newConfig(appName string) *viper.Viper {
	configPath := getConfigDirPath(appName)
	configName := "config.yaml"
	createConfigIfNotExists(configPath, configName)

	v := viper.New()
	v.AddConfigPath(configPath)
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.SetEnvPrefix(strings.ToUpper(appName))
	v.AutomaticEnv()
	v.ReadInConfig()

	return v
}

func getConfigDirPath(appName string) string {
	configPath := os.Getenv(strings.ToUpper(appName) + "_CONFIG")
	if configPath == "" {
		configPath = "/config/cmangos-discord"
	}

	return configPath
}

func createConfigIfNotExists(folderPath string, configFileName string) {
	fullPath := fmt.Sprintf("%s/%s", folderPath, configFileName)
	if !fileExists(fullPath) {
		os.MkdirAll(folderPath, os.ModePerm)
		os.OpenFile(fullPath, os.O_RDONLY|os.O_CREATE, 0666)
	}
}

func getWorkingDir() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return dir
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func createMangosClient() *cmangos.SoapClient {
	mangosUser := os.Getenv("MANGOS_USER")
	mangosPass := os.Getenv("MANGOS_PASS")
	mangosAddress := os.Getenv("MANGOS_ADDRESS")

	client, err := cmangos.NewClient(&cmangos.SoapClientOpts{
		Username: 	mangosUser,
		Password: 	mangosPass,
		Address:	mangosAddress,
	})

	if err != nil {
		panic(err)
	}

	return client
}