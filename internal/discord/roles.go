package discord

import "github.com/bwmarrin/discordgo"

type Role struct {
	Name        string
	Category    string
	DiscordRole *discordgo.Role
}

var (
	RolesToPopulate = []Role{}
)
