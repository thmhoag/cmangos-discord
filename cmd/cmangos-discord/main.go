package main

import (
	"github.com/thmhoag/cmangos-discord/bot"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	defer func() {
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
		<-sc
		log.Println("exiting")
		os.Exit(0)
	}()

	bot.Execute()
}