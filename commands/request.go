package commands

import (
	"github.com/Lukaesebrot/asterisk/config"
	"github.com/Lukaesebrot/asterisk/embeds"
	"github.com/Lukaesebrot/asterisk/users"
	"github.com/Lukaesebrot/dgc"
	"github.com/bwmarrin/discordgo"
)

// Request handles the request command
func Request(ctx *dgc.Ctx) {
	// Validate the input
	if ctx.Arguments.Amount() == 0 {
		ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, embeds.InvalidUsage("You need to specify a feature you want to request."))
		return
	}

	// Send the feature request to the feature request channel and add the delete emote
	message, err := ctx.Session.ChannelMessageSendEmbed(config.CurrentConfig.FeatureRequestChannel, embeds.FeatureRequest(ctx))
	if err != nil {
		ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, embeds.Error("Your feature request couldn't be submitted. Please try again later."))
		return
	}
	ctx.Session.MessageReactionAdd(config.CurrentConfig.FeatureRequestChannel, message.ID, "✅")

	// Confirm the creation of the feature request
	ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, embeds.Success("Your feature request got submitted."))
}

// RequestReactionListener has to be registered to enable the tick reaction on feature requests
func RequestReactionListener(session *discordgo.Session, event *discordgo.MessageReactionAdd) {
	// Check if the channel is the feature request channel
	if event.ChannelID != config.CurrentConfig.FeatureRequestChannel {
		return
	}

	// Check if the user is a bot admin
	user, err := users.RetrieveCached(event.UserID)
	if err != nil || !user.HasPermission(users.PermissionAdministrator) {
		return
	}

	// Check of the reaction is the tick reaction
	if event.Emoji.Name != "✅" {
		return
	}

	// Delete the message
	session.ChannelMessageDelete(event.ChannelID, event.MessageID)
}
