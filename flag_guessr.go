package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"flag-guessr/data"
	"flag-guessr/util"
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/log"
)

func main() {
	data.PopulateCountryMap()

	log.SetLevel(log.LevelInfo)
	log.Info("starting the bot...")
	log.Info("disgo version: ", disgo.Version)

	client, err := disgo.New(os.Getenv("FLAG_GUESSR_TOKEN"),
		bot.WithGatewayConfigOpts(gateway.WithIntents(gateway.IntentsNone),
			gateway.WithPresenceOpts(gateway.WithWatchingActivity("your guesses"))),
		bot.WithCacheConfigOpts(cache.WithCacheFlags(cache.FlagsNone)),
		bot.WithEventListeners(&events.ListenerAdapter{
			OnApplicationCommandInteraction: onCommand,
			OnComponentInteraction:          onButton,
			OnModalSubmit:                   onModal,
		}))
	if err != nil {
		log.Fatal("error while building disgo instance: ", err)
	}

	defer client.Close(context.TODO())

	if client.OpenGateway(context.TODO()) != nil {
		log.Fatalf("error while connecting to the gateway: %s", err)
	}

	log.Infof("flag guessr bot is now running.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}

func onCommand(event *events.ApplicationCommandInteractionCreate) {
	if event.Data.CommandName() == "flag" {
		_ = event.CreateMessage(util.GetCountryCreate(event.User(), 0))
	}
}

func onButton(event *events.ComponentInteractionCreate) {
	id := event.Data.CustomID()
	split := strings.Split(id, "-")
	action := util.Action(split[0])
	cca := split[2]
	country := data.CountryMap[cca]
	name := country.Name.Common
	messageBuilder := discord.NewMessageCreateBuilder()
	if action == util.Details {
		err := event.CreateMessage(messageBuilder.
			SetContentf("Viewing details for **%s** %s %s", name, country.Flag, util.GetCountryInfo(country)).
			SetEphemeral(true).
			Build())
		if err != nil {
			log.Error("there was an error while creating details message: ", err)
		}
		return
	}
	user := event.User()
	userID := user.ID
	if split[1] != userID.String() {
		err := event.CreateMessage(messageBuilder.
			SetContent("You can't interact with games of other users! Launch your own game by using </flag:1007718785345667284>.").
			SetEphemeral(true).
			Build())
		if err != nil {
			log.Error("there was an error while creating error message: ", err)
		}
		return
	}
	streak, _ := strconv.Atoi(split[3])
	client := event.Client().Rest()
	switch action {
	case util.Guess:
		err := event.CreateModal(discord.NewModalCreateBuilder().
			SetCustomID(fmt.Sprintf("%s-%d", cca, streak)).
			SetTitle("Guess the country!").
			AddActionRow(discord.NewShortTextInput("name", "Country name").
				WithPlaceholder("This field is case-insensitive.").
				WithRequired(true)).
			Build())
		if err != nil {
			log.Error("there was an error while creating modal: ", err)
		}
	case util.NewCountry:
		flag := country.Flag
		embedBuilder := discord.NewEmbedBuilder()
		embedBuilder.SetTitle("You skipped a country.")
		embedBuilder.SetDescriptionf("It was **%s**. %s", name, flag)
		embedBuilder.SetColor(0x5386c9)
		util.SendNewCountryMessages(util.NewCountryData{
			Interaction:     event,
			User:            user,
			EmbedBuilder:    *embedBuilder,
			FollowupContent: fmt.Sprintf("The country was **%s**. %s", name, flag),
			Cca:             cca,
			Client:          client,
		})
	case util.Delete:
		err := client.DeleteMessage(event.ChannelID(), event.Message.ID)
		if err != nil {
			log.Error("there was an error while deleting message: ", err)
		}
	case util.Hint:
		i, _ := strconv.Atoi(split[4])
		hintType := util.HintType(i)
		var hint string
		switch hintType {
		case util.Population:
			hint = fmt.Sprintf("The population of this country is %s.", util.FormatPopulation(country))
		case util.Tlds:
			tlds := country.Tlds
			if len(tlds) == 0 {
				hint = "This country has no Top Level Domains."
			} else {
				hint = fmt.Sprintf("The Top Level Domains of this country are **%s**.", strings.Join(tlds, ", "))
			}
		case util.Capitals:
			capitals := country.Capitals
			if len(capitals) == 0 {
				hint = "This country has no capitals."
			} else {
				hint = fmt.Sprintf("The capitals of this country are **%s**.", strings.Join(capitals, ", "))
			}
		}
		err := event.UpdateMessage(discord.NewMessageUpdateBuilder().
			AddActionRow(util.GetGuessButtons(userID, cca, streak, hintType+1)...).
			Build())
		if err != nil {
			log.Error("there was an error while updating message after hint usage: ", err)
		}
		err = util.SendFollowup(event, client, hint)
		if err != nil {
			log.Error("there was an error while creating hint message: ", err)
		}
	}
}

func onModal(event *events.ModalSubmitInteractionCreate) {
	evData := event.Data
	lower := strings.TrimSpace(strings.ToLower(evData.Text("name")))
	id := evData.CustomID
	split := strings.Split(id, "-")
	cca := split[0]
	country := data.CountryMap[cca]
	name := country.Name
	common := name.Common
	if lower == strings.ToLower(common) || lower == strings.ToLower(name.Official) {
		flag := country.Flag
		streak, _ := strconv.Atoi(split[1])
		embedBuilder := discord.NewEmbedBuilder()
		embedBuilder.SetTitle("You got the country right!")
		embedBuilder.SetDescriptionf("It was **%s**. %s", common, flag)
		embedBuilder.SetColor(0x4dbf36)
		util.SendNewCountryMessages(util.NewCountryData{
			Interaction:     event,
			User:            event.User(),
			EmbedBuilder:    *embedBuilder,
			FollowupContent: fmt.Sprintf("Your guess was correct! It was **%s**. %s", common, flag),
			Streak:          streak + 1,
			Cca:             cca,
			Client:          event.Client().Rest(),
		})
	} else {
		err := event.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent("Your guess was incorrect. Please try again.").
			SetEphemeral(true).
			Build())
		if err != nil {
			log.Error("there was an error while creating a followup: ", err)
		}
	}
}
