package main

import (
	"fmt"
	"github.com/thmhoag/cmangos-discord/bot"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	defer func() {
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
		<-sc
		fmt.Println("exiting")
		os.Exit(0)
	}()

	bot.Execute()
}