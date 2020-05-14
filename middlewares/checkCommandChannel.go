package middlewares

import (
	"github.com/Lukaesebrot/asterisk/guildconfig"
	"github.com/Lukaesebrot/asterisk/utils"
	"github.com/Lukaesebrot/dgc"
	"github.com/bwmarrin/discordgo"
)

// CheckCommandChannel checks if the current channel is a valid command channel
func CheckCommandChannel(ctx *dgc.Ctx) bool {
	guildConfig := ctx.CustomObjects["guildConfig"].(*guildconfig.GuildConfig)
	if utils.IsBotAdmin(ctx.Event.Author.ID) {
		return true
	}
	if guildConfig.ChannelRestriction && !utils.StringArrayContains(guildConfig.CommandChannels, ctx.Event.ChannelID) {
		isAdmin, _ := utils.HasPermission(ctx.Session, ctx.Event.GuildID, ctx.Event.Author.ID, discordgo.PermissionAdministrator)
		return isAdmin
	}
	return true
}
