package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/naurffxiv/moddingway/internal/discord"
	"github.com/naurffxiv/moddingway/internal/enum"
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
		fmt.Println("DEBUG")
		guildID := eg.GetEnv("GUILD_ID")
		modLoggingChannelID := eg.GetEnv("MOD_LOGGING_CHANNEL_ID")
		d.Token = discordToken
		d.GuildID = guildID
		d.ModLoggingChannelID = modLoggingChannelID
	} else {
		d.InitWithDefaults(discordToken)
	}

	if !eg.Ok {
		tempstr := fmt.Sprintf("You must supply a %s to start!", eg.EnvName)
		panic(tempstr)
	}

	// db
	var DATABASE_URL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbname)
	conn, err := pgxpool.New(context.Background(), DATABASE_URL)
	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()
	d.Conn = conn

	// discord
	fmt.Printf("Starting Discord...\n")
	err = d.Start()
	if err != nil {
		panic(fmt.Errorf("Could not instantiate Discord: %w", err))
	}
	defer d.Session.Close()
	start(d)

	d.Ready.Wait()
	fmt.Println("Worker is ready. Press CTRL+C to exit.")

	// scheduler
	ticker := time.NewTicker(time.Second * 10)
	done := make(chan bool)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer stop()
	go scheduler(done, ticker, d)
	<-ctx.Done()
	stop()
	done <- true
}

// start adds all the commands and connects the bot to Discord.
// Listens for CTRL+C then terminates the connection.
func start(d *discord.Discord) {
	d.Ready.Add(1)
	// I have 0 idea why AddCommands() is necessary for MapExistingRoles() to work
	d.Session.AddHandler(d.DiscordReady)
	err := d.Session.Open()
	if err != nil {
		panic(fmt.Errorf("Could not open Discord session: %f", err))
	}
}

func scheduler(done chan bool, ticker *time.Ticker, d *discord.Discord) {
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			testQuery(d)
		}
	}
}

func testQuery(d *discord.Discord) {
	type queryRow struct {
		exileID        int
		dbUserID       string
		exileStatus    enum.ExileStatus
		discordUserID  string
		discordGuildID string
	}

	var rowSlice []queryRow
	fmt.Printf("tick\n")
	query := `SELECT e.exileID, e.userID, e.exileStatus, u.discordUserID, u.discordGuildID
				FROM exiles e
				JOIN users u ON e.userID = u.userID
				WHERE e.exileStatus = 2 AND e.endTimestamp < $1;`

	rows, err := d.Conn.Query(context.Background(), query, time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		currRow := queryRow{}
		err := rows.Scan(&currRow.exileID, &currRow.dbUserID, &currRow.exileStatus, &currRow.discordUserID, &currRow.discordGuildID)
		if err != nil {
			fmt.Printf("Error db: %v", err)
			continue
		}
		rowSlice = append(rowSlice, currRow)
		fmt.Println(currRow)
		{
			roleIDToRemove := d.Roles[currRow.discordGuildID]["Exiled"].ID
			roleIDToAdd := d.Roles[currRow.discordGuildID]["Verified"].ID
			// Attempt to remove role first
			err := d.Session.GuildMemberRoleRemove(currRow.discordGuildID, currRow.discordUserID, roleIDToRemove)
			if err != nil {
				// Abort entire process if role removal fails
				tempstr := fmt.Sprintf("Could not remove the role <@&%v> from user <@%v>", roleIDToAdd, currRow.discordUserID)
				fmt.Printf("%v: %v\n", tempstr, err)
				continue
			} else {
				// Otherwise add role
				err = d.Session.GuildMemberRoleAdd(currRow.discordGuildID, currRow.discordUserID, roleIDToAdd)
				if err != nil {
					tempstr := fmt.Sprintf("Could not give user <@%v> role <@&%v>", currRow.discordUserID, roleIDToAdd)
					fmt.Printf("%v: %v\n", tempstr, err)
					continue
				}
			}
			query := `UPDATE exiles
						SET exileStatus = $1
						WHERE exileID = $2`
			_, err = d.Conn.Exec(context.Background(), query, enum.Unexiled, currRow.exileID)
			if err != nil {
				fmt.Printf("Error db: %v", err)
				continue
			}
			fmt.Printf("Automatically unexiled user <@%v>.\n", currRow.discordUserID)
		}
	}
}
