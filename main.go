package main

import (
	"MoneyGoblin4/command_handler"
	"MoneyGoblin4/commands/overview"
	"MoneyGoblin4/commands/reminder"
	"MoneyGoblin4/commands/report"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	// Create a new Discord session using the provided bot token.
	s, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	s.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = s.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	defer s.Close()

	create_commands(s)
	command_handler.Start_Handler(s)
	defer command_handler.Unregister_Commands(s)

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func create_commands(s *discordgo.Session) {
	command_handler.Register_Command(s, "", &overview.Overview_Command{}, &report.Report_Command{}, &reminder.Reminder_Command{})
}
