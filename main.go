package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/naurffxiv/moddingway/internal/discord"
	"github.com/naurffxiv/moddingway/internal/util"
)

func main() {
	eg := &util.EnvGetter{
		Ok: true,
	}

	host := eg.GetEnv("POSTGRES_HOST")
	port := eg.GetEnv("POSTGRES_PORT")
	user := eg.GetEnv("POSTGRES_USER")
	password := eg.GetEnv("POSTGRES_PASSWORD")
	dbname := eg.GetEnv("POSTGRES_DB")

	discordToken := eg.GetEnv("DISCORD_TOKEN")
	discordToken = strings.TrimSpace(discordToken)

	debug := eg.GetEnv("DEBUG")
	debug = strings.ToLower(debug)

	var d = &discord.Discord{}

	if debug == "true" {
		guildID := eg.GetEnv("GUILD_ID")
		modLoggingChannelID := eg.GetEnv("MOD_LOGGING_CHANNEL_ID")
		d.Token = discordToken
		d.GuildID = guildID
		d.ModLoggingChannelID = modLoggingChannelID
	} else {
		// InitWithDefaults only sets default values and does not
		// do anything else, so error checking can come after
		d.InitWithDefaults(discordToken)
	}

	if !eg.Ok {
		tempstr := fmt.Sprintf("You must supply a %s to start!", eg.EnvName)
		panic(tempstr)
	}

	fmt.Printf("Connecting to db...\n")
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbname)
	d.ConnectDatabase(dbUrl)
	defer d.Conn.Close()

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
}
