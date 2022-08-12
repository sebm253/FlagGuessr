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

func GetGuessButtons(userID snowflake.ID, cca string, hintType HintType) []discord.InteractiveComponent {
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
		})
	return []discord.InteractiveComponent{guessButton, newCountryButton, hintButton, GetDeleteButton(userID, cca)}
}

func GetDetailsButton(userID snowflake.ID, cca string) discord.InteractiveComponent {
	return discord.NewSecondaryButton("See country details", fmt.Sprintf(buttonTemplate, Details, userID, cca)).
		WithEmoji(discord.ComponentEmoji{
			Name: "🗺",
		})
}

func GetDeleteButton(userID snowflake.ID, cca string) discord.InteractiveComponent {
	return discord.NewDangerButton("Delete", fmt.Sprintf(buttonTemplate, Delete, userID, cca)).
		WithEmoji(discord.ComponentEmoji{
			Name: "🗑",
		})
}
