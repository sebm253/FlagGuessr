package util

import (
	"fmt"
	"strings"

	"flag-guessr/data"

	"github.com/disgoorg/disgo/discord"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func GetCountryCreate(startData GameStartData, countryData *data.CountryData) discord.MessageCreate {
	user := startData.User
	difficulty := startData.Difficulty
	minPopulation := startData.MinPopulation
	streak := startData.Streak
	ephemeral := startData.Ephemeral
	index, country := countryData.GetRandomCountry(minPopulation)
	embedBuilder := discord.NewEmbedBuilder()
	embedBuilder.SetTitle("Guess the country!")
	embedBuilder.SetDescriptionf("Game of %s\n\nDifficulty: **%s**\nMinimum population: %s\n\nStreak: **%d**", user.Mention(), difficulty, formatRawPopulation(minPopulation), streak)
	embedBuilder.SetColor(0xFFFFFF)
	embedBuilder.SetImage(country.Flags.Png)
	embedBuilder.SetThumbnail(user.EffectiveAvatarURL())
	return discord.NewMessageCreateBuilder().
		SetEmbeds(embedBuilder.Build()).
		AddActionRow(GetGuessButtons(ButtonStateData{
			UserID:        user.ID,
			Difficulty:    difficulty,
			MinPopulation: minPopulation,
			SliceIndex:    index,
			Ephemeral:     ephemeral,
			Streak:        streak,
		})...).
		SetEphemeral(ephemeral).
		Build()
}

func GetCountryInfo(country *data.Country) string {
	capitals := country.Capitals
	tlds := country.Tlds
	population := fmt.Sprintf("Population: %s\n", FormatPopulation(country))
	side := fmt.Sprintf("Driving side: **%s**\n", country.Car.Side)
	capital := fmt.Sprintf("Capital(s): **%s**\n", Ternary(len(capitals) == 0, "None", strings.Join(capitals, ", ")))
	tld := fmt.Sprintf("Top Level Domain(s): **%s**\n", Ternary(len(tlds) == 0, "None", strings.Join(tlds, ", ")))
	gMaps := fmt.Sprintf("Google Maps: **<%s>**\n", country.Maps.GoogleMaps)
	return "\n\n" + population + side + capital + tld + gMaps
}

func FormatPopulation(country *data.Country) string {
	return formatRawPopulation(country.Population)
}

func formatRawPopulation(population int) string {
	p := message.NewPrinter(language.AmericanEnglish)
	return p.Sprintf("**%d**", population)
}
