package features

import (
	"github.com/Lukaesebrot/asterisk/embeds"
	"github.com/Lukaesebrot/asterisk/services/mathjs"
	"github.com/Lukaesebrot/dgc"
)

// initializeMathFeature initializes the math feature
func initializeMathFeature(router *dgc.Router, rateLimiter dgc.RateLimiter) {
	// Register the 'math' command
	router.RegisterCmd(&dgc.Command{
		Name:        "math",
		Description: "Evaluates the given mathematical expression",
		Usage:       "math <codeblock>",
		Example:     "math `1+4^6`\n",
		IgnoreCase:  true,
		RateLimiter: rateLimiter,
		Handler:     mathCommand,
	})
}

// mathCommand handles the 'math' command
func mathCommand(ctx *dgc.Ctx) {
	// Validate the arguments
	codeblock := ctx.Arguments.AsCodeblock()
	if codeblock == nil {
		ctx.RespondEmbed(embeds.InvalidUsage(ctx.Command.Usage))
		return
	}

	// Check the rate limiter
	if !ctx.Command.NotifyRateLimiter(ctx) {
		return
	}

	// Respond with a loading embed
	msg, _ := ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, embeds.Loading())

	// Evaluate the expression and respond with the result
	result, err := mathjs.Evaluate(codeblock.Content)
	if err != nil {
		ctx.Session.ChannelMessageEditEmbed(ctx.Event.ChannelID, msg.ID, embeds.Error(err.Error()))
		return
	}
	ctx.Session.ChannelMessageEditEmbed(ctx.Event.ChannelID, msg.ID, embeds.Success(result))
}
