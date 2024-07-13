package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/naurffxiv/moddingway/internal/discord"
)

func main() {
	discordToken, ok := os.LookupEnv("DISCORD_TOKEN")
	if !ok {
		panic("You must supply a DISCORD_TOKEN to start!")
	}
	discordToken = strings.TrimSpace(discordToken)

	d := &discord.Discord{}

	debug, _ := os.LookupEnv("DEBUG")
	debug = strings.ToLower(debug)

	if debug == "true" {
		guildID, ok := os.LookupEnv("GUILD_ID")
		if !ok {
			panic("You must supply a GUILD_ID to start!")
		}

		modLoggingChannelID, ok := os.LookupEnv("MOD_LOGGING_CHANNEL_ID")
		if !ok {
			panic("You must supply a MOD_LOGGING_CHANNEL_ID to start!")
		}

		d.InitDebug(discordToken, guildID, modLoggingChannelID)
	} else {
		d.Init(discordToken)
	}

	fmt.Printf("Starting Discord...\n")
	err := d.Start()
	if err != nil {
		panic(fmt.Errorf("Could not instantiate Discord: %w", err))
	}
	defer d.Session.Close()
	start(d)
}

// start adds all the commands and connects the bot to Discord.
// Listens for CTRL+C then terminates the connection.
func start(d *discord.Discord) {
	d.Ready.Add(1)
	d.Session.AddHandler(d.DiscordReady)
	err := d.Session.Open()
	if err != nil {
		panic(fmt.Errorf("Could not open Discord session: %f", err))
	}

	d.Ready.Wait()
	d.Session.AddHandler(d.InteractionCreate)
	fmt.Println("Moddingway is ready. Press CTRL+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	d.Session.Close()
}
