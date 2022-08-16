package util

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/log"
)

func SendFollowup(interaction discord.BaseInteraction, client rest.Rest, content string) error {
	_, err := client.CreateFollowupMessage(interaction.ApplicationID(), interaction.Token(),
		discord.NewMessageCreateBuilder().
			SetContent(content).
			SetEphemeral(true).
			Build())
	return err
}

func SendNewCountryMessages(data NewCountryData) {
	client := data.Client
	interaction := data.Interaction
	token := interaction.Token()
	embedBuilder := data.EmbedBuilder
	user := data.User
	userID := user.ID
	err := client.CreateInteractionResponse(interaction.ID(), token, discord.InteractionResponse{
		Type: discord.InteractionResponseTypeUpdateMessage,
		Data: discord.NewMessageUpdateBuilder().
			SetEmbeds(embedBuilder.Build()).
			AddActionRow(GetDetailsButton(userID, data.Cca)).
			Build(),
	})
	if err != nil {
		log.Error("there was an error while updating original message: ", err)
	}
	_, err = client.CreateFollowupMessage(interaction.ApplicationID(), token, GetCountryCreate(user, HintType(0)))
	if err != nil {
		log.Error("there was an error while creating new country message: ", err)
	}
	err = SendFollowup(interaction, client, data.FollowupContent)
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
