package util

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
)

func SendFollowup(interaction discord.BaseInteraction, client rest.Rest, content string, ephemeral bool, buttons ...discord.InteractiveComponent) error {
	messageBuilder := discord.NewMessageCreateBuilder()
	messageBuilder.SetContent(content)
	messageBuilder.SetEphemeral(ephemeral)
	if len(buttons) != 0 {
		messageBuilder.AddActionRow(buttons...)
	}
	_, err := client.CreateFollowupMessage(interaction.ApplicationID(), interaction.Token(), messageBuilder.Build())
	return err
}

func Ternary[T any](exp bool, ifCond T, elseCond T) T { // https://github.com/aidenwallis/go-utils/blob/main/utils/ternary.go
	if exp {
		return ifCond
	}
	return elseCond
}
