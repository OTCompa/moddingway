package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// mapOptions is a helper function that creates a map out of the arguments used in the slash command
func mapOptions(i *discordgo.InteractionCreate) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	return optionMap
}

// Kick attempts to kick the user specified user from the server the command was invoked in.
// Fields:
//
//	user: 	User
//	reason: string
func (d *Discord) Kick(s *discordgo.Session, i *discordgo.InteractionCreate) {
	dmFailed := false
	optionMap := mapOptions(i)

	// Log usage of command
	logMsg, err := d.LogCommand(i.Interaction)
	if err != nil {
		fmt.Printf("Failed to log: %v\n", err)
	}

	userToKick := optionMap["user"].UserValue(nil).ID

	// Check if user exists in guild
	err = d.CheckUserInGuild(i.GuildID, userToKick)
	if err != nil {
		tempstr := fmt.Sprintf("Could not kick user <@%v>", userToKick)
		fmt.Printf("%v: %v\n", tempstr, err)

		err = StartInteraction(s, i.Interaction, tempstr)
		if err != nil {
			fmt.Printf("Unable to send ephemeral message: %v\n", err)
		}

		return
	}

	// DM the user regarding the kick
	channel, err := s.UserChannelCreate(userToKick)
	if err != nil {
		tempstr := fmt.Sprintf("Could not create a DM with user %v", userToKick)
		fmt.Printf("%v: %v\n", tempstr, err)
		dmFailed = true
	} else {
		tempstr := fmt.Sprintf("You are being kicked from %v for the reason:\n%v",
			GuildName,
			optionMap["reason"].StringValue(),
		)

		_, err = s.ChannelMessageSend(channel.ID, tempstr)
		if err != nil {
			tempstr := fmt.Sprintf("Could not send a DM to user %v", userToKick)
			fmt.Printf("%v: %v\n", tempstr, err)
			dmFailed = true
		}
	}

	// Attempt to kick user
	if len(optionMap["reason"].StringValue()) > 0 {
		err = d.Session.GuildMemberDeleteWithReason(i.GuildID, userToKick, optionMap["reason"].StringValue())
	} else {
		err = StartInteraction(s, i.Interaction, "Please provide a reason for the kick.")
		if err != nil {
			fmt.Printf("Unable to send ephemeral message: %v\n", err)
		}

		return
	}

	if err != nil {
		tempstr := fmt.Sprintf("Could not kick user <@%v>", userToKick)
		fmt.Printf("%v: %v\n", tempstr, err)

		err = StartInteraction(s, i.Interaction, tempstr)
		if err != nil {
			fmt.Printf("Unable to send ephemeral message: %v\n", err)
		}
		return
	} else {
		tempstr := fmt.Sprintf("User <@%v> has been kicked.", userToKick)
		fmt.Printf("%v\n", tempstr)

		err = StartInteraction(s, i.Interaction, tempstr)
		if err != nil {
			fmt.Printf("Unable to send ephemeral message: %v\n", err)
		}
	}

	// Inform of failure to DM
	if dmFailed {
		err = ContinueInteraction(s, i.Interaction, "Unable to send DM to user.")
		if err != nil {
			fmt.Printf("Unable to send ephemeral followup message: %v\n", err)
		}
		logMsg.Embeds[0].Description += "\nFailed to notify the user of the kick via DM."
		_, err = d.Session.ChannelMessageEditEmbed(d.LogChannelID, logMsg.ID, logMsg.Embeds[0])
		if err != nil {
			fmt.Printf("Unable to edit log message: %v\n", err)
		}
	}
}

// Mute attempts to mute the user specified user from the server the command was invoked in.
// Fields:
//
//	user: 		User
//	duration:	string
//	reason:		string
func (d *Discord) Mute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return
}

// Unmute attempts to unmute the user specified user from the server the command was invoked in.
// Fields:
//
//	user: 		User
//	reason:		string
func (d *Discord) Unmute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return
}

// Ban attempts to ban the user specified user from the server the command was invoked in.
// Fields:
//
//	user:		User
//	reason:		string
func (d *Discord) Ban(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return
}

// Unban attempts to unban the user specified user from the server the command was invoked in.
// Fields:
//
//	user:		User
//	reason:		string
func (d *Discord) Unban(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return
}

// RemoveNickname attempts to remove the currently set nickname on the specified user
// in the server the command was invoked in.
// Fields:
//
//	user:		User
//	reason:		string
func (d *Discord) RemoveNickname(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return
}

// SetNickname attempts to set the nickname of the specified user in the server
// the command was invoked in.
// Fields:
//
//	user:		User
//	nickname:	string
//	reason:		string
func (d *Discord) SetNickname(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return
}

// Slowmode attempts to set the current channel to slowmode.
// Fields:
//
//	duration:	string
func (d *Discord) Slowmode(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return
}

// SlowmodeOff attempts to remove slowmode from the current channel.
func (d *Discord) SlowmodeOff(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return
}

// Purge attempts to remove the last message-number messages from the specified channel.
// Fields:
//
//	channel:		Channel
//	message-number:	integer
func (d *Discord) Purge(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return
}

// Exile attempts to add the exile role to the user, effectively soft-banning them.
// Fields:
//
//	user:		User
//	reason:		string
func (d *Discord) Exile(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return
}

// Unexile attempts to remove the exile role from the user.
// Fields:
//
//	user:		User
//	reason:		string
func (d *Discord) Unexile(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return
}

// SetModLog sets the specified channel to the moderation log channel
// All logged commands will be logged to this channel.
// Fields:
//
//	channel:	Channel
func (d *Discord) SetModLog(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	channelID := options[0].ChannelValue(nil).ID
	d.LogChannelID = channelID

	tempstr := fmt.Sprintf("Mod log channel set to: <#%v>", channelID)

	err := StartInteraction(s, i.Interaction, tempstr)
	if err != nil {
		fmt.Printf("Unable to send ephemeral message: %v\n", err)
	}
	fmt.Printf("Set the moderation log channel to channel: %v\n", channelID)
}
