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
			gateway.WithPresence(gateway.NewWatchingPresence("your guesses", discord.OnlineStatusOnline, false))),
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
		_ = event.CreateMessage(util.GetCountryCreate(event.User().ID, util.HintType(0)))
	}
}

func onButton(event *events.ComponentInteractionCreate) {
	id := event.Data.CustomID()
	split := strings.Split(id, "-")
	action := util.Action(split[0])
	user := split[1]
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
	userID := event.User().ID
	if user != userID.String() {
		err := event.CreateMessage(messageBuilder.
			SetContent("You can't interact with games of other users! Launch your own game by using </flag:1007718785345667284>.").
			SetEphemeral(true).
			Build())
		if err != nil {
			log.Error("there was an error while creating error message: ", err)
		}
		return
	}
	client := event.Client().Rest()
	if action == util.Guess {
		err := event.CreateModal(discord.NewModalCreateBuilder().
			SetCustomID(cca).
			SetTitle("Guess the flag!").
			AddActionRow(discord.NewShortTextInput("name", "Country name").
				WithRequired(true)).
			Build())
		if err != nil {
			log.Error("there was an error while creating modal: ", err)
		}
	} else if action == util.NewCountry {
		messageUpdate := util.GetCountryUpdate(userID, util.HintType(0))
		err := event.UpdateMessage(messageUpdate)
		if err != nil {
			log.Error("there was an error while updating message with new country: ", err)
		}
		content := fmt.Sprintf("The country was **%s**.", name)
		util.SendFollowup(client, event.BaseInteraction, content, util.GetDetailsButton(userID, cca), util.GetDeleteButton(userID, cca))
	} else if action == util.Delete {
		err := client.DeleteMessage(event.ChannelID(), event.Message.ID)
		if err != nil {
			log.Error("there was an error while deleting message: ", err)
		}
	} else if action == util.Hint {
		i, _ := strconv.Atoi(split[3])
		hintType := util.HintType(i)
		var hint string
		if hintType == util.Population {
			hint = fmt.Sprintf("The population of this country is %s.", util.FormatPopulation(country))
		} else if hintType == util.Tlds {
			hint = fmt.Sprintf("The Top Level Domains of this country are **%s**.", strings.Join(country.Tlds, ", "))
		} else if hintType == util.Capitals {
			hint = fmt.Sprintf("The capitals of this country are **%s**.", strings.Join(country.Capitals, ", "))
		}
		messageUpdate := discord.MessageUpdate{
			Embeds: &event.Message.Embeds,
		}
		messageUpdate.Components = &[]discord.ContainerComponent{discord.ActionRowComponent(util.GetGuessButtons(userID, cca, hintType+1))}
		err := event.UpdateMessage(messageUpdate)
		if err != nil {
			log.Error("there was an error while updating message with new hint: ", err)
		}
		util.SendFollowup(client, event.BaseInteraction, hint, util.GetDeleteButton(userID, cca))
	}
}

func onModal(event *events.ModalSubmitInteractionCreate) {
	evData := event.Data
	lower := strings.ToLower(evData.Text("name"))
	cca := evData.CustomID
	country := data.CountryMap[cca]
	name := country.Name
	common := name.Common
	messageBuilder := discord.NewMessageCreateBuilder()
	var err error
	if lower == strings.ToLower(common) || lower == strings.ToLower(name.Official) {
		userID := event.User().ID
		err = event.UpdateMessage(util.GetCountryUpdate(userID, util.HintType(0)))
		if err != nil {
			log.Error("there was an error while updating original message: ", err)
		}
		client := event.Client().Rest()
		content := fmt.Sprintf("Your guess was correct! It was **%s**. New flag to guess has been prepared!", common)
		util.SendFollowup(client, event.BaseInteraction, content, util.GetDetailsButton(userID, cca), util.GetDeleteButton(userID, cca))
	} else {
		err = event.CreateMessage(messageBuilder.
			SetContent("Your guess was incorrect. Please try again.").
			SetEphemeral(true).
			Build())
	}
	if err != nil {
		log.Error("there was an error while creating a followup: ", err)
	}
}
