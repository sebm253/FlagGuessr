package util

import (
	"fmt"
	"math/rand"
	"reflect"
	"strings"

	"flag-guessr/data"
	"github.com/disgoorg/disgo/discord"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func getCountry(user discord.User, hintType HintType) (discord.Embed, []discord.InteractiveComponent) {
	keys := reflect.ValueOf(data.CountryMap).MapKeys()
	cca := keys[rand.Intn(len(keys))].String()
	country := data.CountryMap[cca]
	userID := user.ID
	embedBuilder := discord.NewEmbedBuilder()
	embedBuilder.SetTitle("Guess the country!")
	embedBuilder.SetDescriptionf("Game of <@%d>", userID)
	embedBuilder.SetColor(0xFFFFFF)
	embedBuilder.SetImage(country.Flags.Png)
	embedBuilder.SetThumbnail(user.EffectiveAvatarURL())
	embedBuilder.SetFooterText("Country data provided by restcountries.com")
	return embedBuilder.Build(), GetGuessButtons(userID, cca, hintType, false)
}

func GetCountryCreate(user discord.User, hintType HintType) discord.MessageCreate {
	embed, buttons := getCountry(user, hintType)
	return discord.NewMessageCreateBuilder().
		SetEmbeds(embed).
		AddActionRow(buttons...).
		Build()
}

func GetCountryUpdate(user discord.User, hintType HintType) discord.MessageUpdate {
	embed, buttons := getCountry(user, hintType)
	return discord.NewMessageUpdateBuilder().
		SetEmbeds(embed).
		AddActionRow(buttons...).
		Build()
}

func GetCountryInfo(country data.Country) string {
	capitals := country.Capitals
	population := fmt.Sprintf("Population: %s\n", FormatPopulation(country))
	capital := fmt.Sprintf("Capital(s): **%s**\n", Ternary(len(capitals) == 0, "None", strings.Join(capitals, ", ")))
	tld := fmt.Sprintf("Top Level Domain(s): **%s**\n", strings.Join(country.Tlds, ", "))
	maps := fmt.Sprintf("Google Maps: **<%s>**\n", country.Maps.GoogleMaps)
	return "\n\n" + population + capital + tld + maps
}

func FormatPopulation(country data.Country) string {
	p := message.NewPrinter(language.AmericanEnglish)
	return p.Sprintf("**%d**", country.Population)
}
