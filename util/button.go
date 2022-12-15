package util

import (
	"encoding/json"
	"github.com/disgoorg/disgo/discord"
)

func GetGuessButtons(stateData ButtonStateData) []discord.InteractiveComponent {
	guessButton := discord.NewPrimaryButton("Submit guess", marshalStateData(stateData, ActionTypeGuess)).
		WithEmoji(discord.ComponentEmoji{
			Name: "🍀",
		})
	newCountryButton := discord.NewSecondaryButton("New country", marshalStateData(stateData, ActionTypeNewCountry)).
		WithEmoji(discord.ComponentEmoji{
			Name: "♻",
		})
	hintButton := discord.NewSecondaryButton("Hint", marshalStateData(stateData, ActionTypeHint)).
		WithEmoji(discord.ComponentEmoji{
			Name: "❓",
		}).
		WithDisabled(stateData.HintType == HintTypeUnknown)
	return []discord.InteractiveComponent{guessButton, newCountryButton, hintButton}
}

func marshalStateData(stateData ButtonStateData, actionType ActionType) string {
	stateData.ActionType = actionType
	data, _ := json.Marshal(stateData)
	return string(data)
}
