package util

import (
	"encoding/json"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/log"
)

func SendGameUpdates(data NewCountryData) {
	client := data.Client
	interaction := data.Interaction
	token := interaction.Token()
	user := data.User
	err := client.CreateInteractionResponse(interaction.ID(), token, discord.InteractionResponse{
		Type: discord.InteractionResponseTypeUpdateMessage,
		Data: GetCountryCreate(user, data.Difficulty, data.Streak),
	})
	if err != nil {
		log.Error("there was an error while creating new country message: ", err)
	}
	stateData, _ := json.Marshal(&ButtonStateData{
		UserID:     user.ID,
		Cca:        data.Cca,
		ActionType: ActionTypeDetails,
	})
	_, err = client.CreateFollowupMessage(interaction.ApplicationID(), token, discord.NewMessageCreateBuilder().
		SetContent(data.FollowupContent).
		AddActionRow(discord.NewSecondaryButton("See country details", string(stateData)).
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
