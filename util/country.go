package util

import (
	"fmt"
	"math/rand"
	"strings"

	"flag-guessr/data"
	"github.com/disgoorg/disgo/discord"
	"golang.org/x/exp/maps"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func GetCountryCreate(user discord.User, streak int) discord.MessageCreate {
	keys := maps.Keys(data.CountryMap)
	cca := keys[rand.Intn(len(keys))]
	country := data.CountryMap[cca]
	userID := user.ID
	embedBuilder := discord.NewEmbedBuilder()
	embedBuilder.SetTitle("Guess the country!")
	embedBuilder.SetDescriptionf("Game of <@%d>\n\nStreak: **%d**", userID, streak)
	embedBuilder.SetColor(0xFFFFFF)
	embedBuilder.SetImage(country.Flags.Png)
	embedBuilder.SetThumbnail(user.EffectiveAvatarURL())
	embedBuilder.SetFooterText("Country data provided by restcountries.com")
	return discord.NewMessageCreateBuilder().
		SetEmbeds(embedBuilder.Build()).
		AddActionRow(GetGuessButtons(userID, cca, streak, HintType(0))...).
		Build()
}

func GetCountryInfo(country data.Country) string {
	capitols := country.Capitols
	tlds := country.Tlds
	population := fmt.Sprintf("Population: %s\n", FormatPopulation(country))
	capitol := fmt.Sprintf("Capitol(s): **%s**\n", Ternary(len(capitols) == 0, "None", strings.Join(capitols, ", ")))
	tld := fmt.Sprintf("Top Level Domain(s): **%s**\n", Ternary(len(tlds) == 0, "None", strings.Join(tlds, ", ")))
	gMaps := fmt.Sprintf("Google Maps: **<%s>**\n", country.Maps.GoogleMaps)
	return "\n\n" + population + capitol + tld + gMaps
}

func FormatPopulation(country data.Country) string {
	p := message.NewPrinter(language.AmericanEnglish)
	return p.Sprintf("**%d**", country.Population)
}
