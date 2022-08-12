package util

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/log"
)

func SendFollowup(client rest.Rest, interaction discord.BaseInteraction, content string, buttons ...discord.InteractiveComponent) {
	_, err := client.CreateFollowupMessage(interaction.ApplicationID(), interaction.Token(), discord.NewMessageCreateBuilder().
		SetContent(content).
		AddActionRow(buttons...).
		Build())
	if err != nil {
		log.Errorf("there was an error while sending followup: ", err)
	}
}

func Ternary[T any](exp bool, ifCond T, elseCond T) T { // https://github.com/aidenwallis/go-utils/blob/main/utils/ternary.go
	if exp {
		return ifCond
	}
	return elseCond
}
