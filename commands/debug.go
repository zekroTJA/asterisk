package commands

import (
	"fmt"
	"reflect"

	"github.com/Lukaesebrot/asterisk/config"
	"github.com/Lukaesebrot/asterisk/embeds"
	"github.com/Lukaesebrot/asterisk/users"
	"github.com/Lukaesebrot/dgc"
	"github.com/containous/yaegi/interp"
	"github.com/containous/yaegi/stdlib"
)

// Debug handles the debug command
func Debug(ctx *dgc.Ctx) {
	// Validate the arguments
	codeblock := ctx.Arguments.AsCodeblock()
	if codeblock == nil {
		ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, embeds.InvalidUsage(ctx.Command.Usage))
		return
	}

	// Create the interpreter
	interpreter := interp.New(interp.Options{})

	// Inject the custom variables
	custom := make(map[string]map[string]reflect.Value)
	custom["asterisk"] = map[string]reflect.Value{
		"ctx":     reflect.ValueOf(ctx),
		"getUser": reflect.ValueOf(users.RetrieveCached),
		"config":  reflect.ValueOf(config.CurrentConfig),
	}
	interpreter.Use(stdlib.Symbols)
	interpreter.Use(custom)
	_, err := interpreter.Eval("import (\n. \"asterisk\"\n\"fmt\"\n\"time\"\n)")
	if err != nil {
		ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, embeds.Error(err.Error()))
		return
	}

	// Evaluate the given string and respond with the result
	result, err := interpreter.Eval(codeblock.Content)
	if err != nil {
		ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, embeds.Error(err.Error()))
		return
	}
	ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, embeds.Success(fmt.Sprintf("%+v", result)))
}
