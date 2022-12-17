package util

import (
	"fmt"
	"math/rand"
	"strings"

	"flag-guessr/data"
	"github.com/disgoorg/disgo/discord"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func GetCountryCreate(startData GameStartData) discord.MessageCreate {
	user := startData.User
	userID := user.ID
	streak := startData.Streak
	minPopulation := startData.MinPopulation
	countryIndex, country := getRandomCountry(minPopulation)
	embedBuilder := discord.NewEmbedBuilder()
	embedBuilder.SetTitle("Guess the country!")
	embedBuilder.SetDescriptionf("Game of <@%d>\n\nMinimum population: %s\nStreak: **%d**", userID, formatRawPopulation(minPopulation), streak)
	embedBuilder.SetColor(0xFFFFFF)
	embedBuilder.SetImage(country.Flags.Png)
	embedBuilder.SetThumbnail(user.EffectiveAvatarURL())
	embedBuilder.SetFooterText("Country data provided by restcountries.com")
	return discord.NewMessageCreateBuilder().
		SetEmbeds(embedBuilder.Build()).
		AddActionRow(GetGuessButtons(ButtonStateData{
			UserID:        userID,
			Difficulty:    startData.Difficulty,
			MinPopulation: minPopulation,
			SliceIndex:    countryIndex,
			Ephemeral:     startData.Ephemeral,
			Streak:        streak,
		})...).
		SetEphemeral(startData.Ephemeral).
		Build()
}

func GetCountryInfo(country data.Country) string {
	capitals := country.Capitals
	tlds := country.Tlds
	population := fmt.Sprintf("Population: %s\n", FormatPopulation(country))
	capital := fmt.Sprintf("Capital(s): **%s**\n", Ternary(len(capitals) == 0, "None", strings.Join(capitals, ", ")))
	tld := fmt.Sprintf("Top Level Domain(s): **%s**\n", Ternary(len(tlds) == 0, "None", strings.Join(tlds, ", ")))
	gMaps := fmt.Sprintf("Google Maps: **<%s>**\n", country.Maps.GoogleMaps)
	return "\n\n" + population + capital + tld + gMaps
}

func FormatPopulation(country data.Country) string {
	return formatRawPopulation(country.Population)
}

func getRandomCountry(minPopulation int) (int, data.Country) {
	i := rand.Intn(data.IndexBoundaries[minPopulation])
	return i, data.CountrySlice[i]
}

func formatRawPopulation(population int) string {
	p := message.NewPrinter(language.AmericanEnglish)
	return p.Sprintf("**%d**", population)
}
