package util

import (
	"encoding/json"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/log"
)

func SendGameUpdates(data *NewCountryData) {
	client := data.Client
	interaction := data.Interaction
	token := interaction.Token()
	user := interaction.User()
	err := client.CreateInteractionResponse(interaction.ID(), token, discord.InteractionResponse{
		Type: discord.InteractionResponseTypeDeferredUpdateMessage,
	})
	if err != nil {
		log.Error("there was an error while deferring update message: ", err)
		return
	}
	applicationID := interaction.ApplicationID()
	if err := client.DeleteInteractionResponse(applicationID, token); err != nil {
		log.Error("there was an error while deleting original interaction response: ", err)
		return
	}
	_, err = client.CreateFollowupMessage(applicationID, token, GetCountryCreate(&GameStartData{
		User:          &user,
		Difficulty:    data.Difficulty,
		MinPopulation: data.MinPopulation,
		Ephemeral:     data.Ephemeral,
		Streak:        data.Streak,
	}, data.CountryData))
	if err != nil {
		log.Error("there was an error while creating new country message: ", err)
		return
	}
	stateData, _ := json.Marshal(&ButtonStateData{
		UserID:     user.ID,
		SliceIndex: data.SliceIndex,
		ActionType: ActionTypeDetails,
	})
	_, err = client.CreateFollowupMessage(applicationID, token, discord.NewMessageCreateBuilder().
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
