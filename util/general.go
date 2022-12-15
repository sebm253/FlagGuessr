package util

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/log"
)

func SendNewCountryMessages(data NewCountryData) {
	client := data.Client
	interaction := data.Interaction
	token := interaction.Token()
	user := data.User
	err := client.CreateInteractionResponse(interaction.ID(), token, discord.InteractionResponse{
		Type: discord.InteractionResponseTypeUpdateMessage,
		Data: GetCountryCreate(user, data.Streak),
	})
	if err != nil {
		log.Error("there was an error while creating new country message: ", err)
	}
	_, err = client.CreateFollowupMessage(interaction.ApplicationID(), token, discord.NewMessageCreateBuilder().
		SetContent(data.FollowupContent).
		AddActionRow(discord.NewSecondaryButton("See country details", fmt.Sprintf(buttonTemplate, Details, user.ID, data.Cca, 0)).
			WithEmoji(discord.ComponentEmoji{
				Name: "🗺",
			})).
		SetEphemeral(true).
		Build())
	if err != nil {
		log.Error("there was an error while creating new country info message: ", err)
	}
}

func Ternary[T any](exp bool, ifCond T, elseCond T) T { // https://github.com/aidenwallis/go-utils/blob/main/utils/ternary.go
	if exp {
		return ifCond
	}
	return elseCond
}
