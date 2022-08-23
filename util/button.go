package util

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

var (
	buttonTemplate     = "%s-%d-%s"
	hintButtonTemplate = buttonTemplate + "-%d"
)

func GetGuessButtons(userID snowflake.ID, cca string, hintType HintType, hintDisabled bool) []discord.InteractiveComponent {
	guessButton := discord.NewPrimaryButton("Guess", fmt.Sprintf(buttonTemplate, Guess, userID, cca)).
		WithEmoji(discord.ComponentEmoji{
			Name: "🍀",
		})
	newCountryButton := discord.NewSecondaryButton("New country", fmt.Sprintf(buttonTemplate, NewCountry, userID, cca)).
		WithEmoji(discord.ComponentEmoji{
			Name: "♻",
		})
	hintButton := discord.NewSecondaryButton("Hint", fmt.Sprintf(hintButtonTemplate, Hint, userID, cca, hintType)).
		WithEmoji(discord.ComponentEmoji{
			Name: "❓",
		}).
		WithDisabled(hintDisabled)
	deleteButton := discord.NewDangerButton("Delete", fmt.Sprintf(buttonTemplate, Delete, userID, cca)).
		WithEmoji(discord.ComponentEmoji{
			Name: "🗑",
		})
	return []discord.InteractiveComponent{guessButton, newCountryButton, hintButton, deleteButton}
}
