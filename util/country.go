package util

import (
	"fmt"
	"math/rand"
	"strings"

	"flag_guessr/data"
	"github.com/disgoorg/disgo/discord"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func GetCountryCreate(startData GameStartData) discord.MessageCreate {
	user := startData.User
	userID := user.ID
	difficulty := startData.Difficulty
	minPopulation := startData.MinPopulation
	streak := startData.Streak
	ephemeral := startData.Ephemeral
	countryIndex, country := getRandomCountry(minPopulation)
	embedBuilder := discord.NewEmbedBuilder()
	embedBuilder.SetTitle("Guess the country!")
	embedBuilder.SetDescriptionf("Game of <@%d>\n\nDifficulty: **%s**\nMinimum population: %s\n\nStreak: **%d**", userID, difficulty, formatRawPopulation(minPopulation), streak)
	embedBuilder.SetColor(0xFFFFFF)
	embedBuilder.SetImage(country.Flags.Png)
	embedBuilder.SetThumbnail(user.EffectiveAvatarURL())
	embedBuilder.SetFooterText("Country data provided by restcountries.com")
	return discord.NewMessageCreateBuilder().
		SetEmbeds(embedBuilder.Build()).
		AddActionRow(GetGuessButtons(ButtonStateData{
			UserID:        userID,
			Difficulty:    difficulty,
			MinPopulation: minPopulation,
			SliceIndex:    countryIndex,
			Ephemeral:     ephemeral,
			Streak:        streak,
		})...).
		SetEphemeral(ephemeral).
		Build()
}

func GetCountryInfo(country data.Country) string {
	capitals := country.Capitals
	tlds := country.Tlds
	population := fmt.Sprintf("Population: %s\n", FormatPopulation(country))
	side := fmt.Sprintf("Driving side: **%s**\n", country.Car.Side)
	capital := fmt.Sprintf("Capital(s): **%s**\n", Ternary(len(capitals) == 0, "None", strings.Join(capitals, ", ")))
	tld := fmt.Sprintf("Top Level Domain(s): **%s**\n", Ternary(len(tlds) == 0, "None", strings.Join(tlds, ", ")))
	gMaps := fmt.Sprintf("Google Maps: **<%s>**\n", country.Maps.GoogleMaps)
	return "\n\n" + population + side + capital + tld + gMaps
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
