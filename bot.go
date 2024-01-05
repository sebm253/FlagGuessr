package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
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
	"github.com/disgoorg/json"
	"github.com/disgoorg/log"
)

func main() {
	countryData := &data.CountryData{}
	countryData.Populate()

	log.SetLevel(log.LevelInfo)
	log.Info("starting the bot...")
	log.Info("disgo version: ", disgo.Version)

	client, err := disgo.New(os.Getenv("FLAG_GUESSR_TOKEN"),
		bot.WithGatewayConfigOpts(gateway.WithIntents(gateway.IntentsNone),
			gateway.WithPresenceOpts(gateway.WithWatchingActivity("your guesses"))),
		bot.WithCacheConfigOpts(cache.WithCaches(cache.FlagsNone)),
		bot.WithEventListeners(&events.ListenerAdapter{
			OnApplicationCommandInteraction: func(event *events.ApplicationCommandInteractionCreate) {
				onCommand(event, countryData)
			},
			OnComponentInteraction: func(event *events.ComponentInteractionCreate) {
				onButton(event, countryData)
			},
			OnModalSubmit: func(event *events.ModalSubmitInteractionCreate) {
				onModal(event, countryData)
			},
		}))
	if err != nil {
		log.Fatal("error while building disgo instance: ", err)
	}

	defer client.Close(context.TODO())

	if err := client.OpenGateway(context.TODO()); err != nil {
		log.Fatal("error while connecting to the gateway: ", err)
	}

	log.Info("flag guessr bot is now running.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}

func onCommand(event *events.ApplicationCommandInteractionCreate, countryData *data.CountryData) {
	interactionData := event.SlashCommandInteractionData()
	if interactionData.CommandName() == "flag" {
		ephemeral, ok := interactionData.OptBool("hide")
		if !ok {
			ephemeral = true
		}
		_ = event.CreateMessage(util.GetCountryCreate(&util.GameStartData{
			User:          json.Ptr(event.User()),
			Difficulty:    util.GameDifficulty(interactionData.Int("difficulty")),
			MinPopulation: interactionData.Int("min-population"),
			Ephemeral:     ephemeral,
		}, countryData))
	}
}

func onButton(event *events.ComponentInteractionCreate, countryData *data.CountryData) {
	var stateData util.ButtonStateData
	buttonID := event.Data.CustomID()
	_ = json.Unmarshal([]byte(buttonID), &stateData)

	messageBuilder := discord.NewMessageCreateBuilder().SetEphemeral(true)
	actionType := stateData.ActionType
	countryIndex := stateData.SliceIndex
	country := countryData.Countries[countryIndex]
	if actionType == util.ActionTypeDetails {
		err := event.CreateMessage(messageBuilder.
			SetContentf("Viewing details for **%s** %s %s", country.Name.Common, country.Flag, util.GetCountryInfo(country)).
			Build())
		if err != nil {
			log.Error("there was an error while creating details message: ", err)
		}
		return
	}
	if stateData.UserID != event.User().ID {
		err := event.CreateMessage(messageBuilder.
			SetContent("You can't interact with games of other users! Launch your own game by using </flag:1007718785345667284>.").
			Build())
		if err != nil {
			log.Error("there was an error while creating error message: ", err)
		}
		return
	}
	difficulty := stateData.Difficulty
	minPopulation := stateData.MinPopulation
	ephemeral := stateData.Ephemeral
	client := event.Client().Rest()
	switch actionType {
	case util.ActionTypeGuess:
		marshalledData, _ := json.Marshal(util.ModalStateData{
			Difficulty:    difficulty,
			MinPopulation: minPopulation,
			SliceIndex:    countryIndex,
			Ephemeral:     ephemeral,
			Streak:        stateData.Streak,
		})
		err := event.Modal(discord.NewModalCreateBuilder().
			SetCustomID(string(marshalledData)).
			SetTitle("Guess the country!").
			AddActionRow(discord.NewShortTextInput("input", "Country name").
				WithPlaceholder("This field is case-insensitive.").
				WithRequired(true)).
			Build())
		if err != nil {
			log.Error("there was an error while creating modal: ", err)
		}
	case util.ActionTypeNewCountry:
		util.SendGameUpdates(&util.NewCountryData{
			Interaction:     event,
			FollowupContent: fmt.Sprintf("You skipped a country. It was **%s**. %s", country.Name.Common, country.Flag),
			Difficulty:      difficulty,
			MinPopulation:   minPopulation,
			Ephemeral:       ephemeral,
			SliceIndex:      countryIndex,
			Client:          client,
			CountryData:     countryData,
		})
	case util.ActionTypeDelete:
		if err := client.DeleteMessage(event.Channel().ID(), event.Message.ID); err != nil {
			log.Error("there was an error while deleting message: ", err)
		}
	case util.ActionTypeHint:
		hintType := stateData.HintType
		switch hintType {
		case util.HintTypeDrivingSide:
			messageBuilder.SetContentf("This country drives on the **%s**.", country.Car.Side)
		case util.HintTypePopulation:
			messageBuilder.SetContentf("The population of this country is %s.", util.FormatPopulation(country))
		case util.HintTypeCapitals:
			capitals := country.Capitals
			if len(capitals) == 0 {
				messageBuilder.SetContent("This country has no capitals.")
			} else {
				messageBuilder.SetContentf("The capitals of this country are **%s**.", strings.Join(capitals, ", "))
			}
		case util.HintTypeTlds:
			tlds := country.Tlds
			if len(tlds) == 0 {
				messageBuilder.SetContent("This country has no Top Level Domains.")
			} else {
				messageBuilder.SetContentf("The Top Level Domains of this country are **%s**.", strings.Join(tlds, ", "))
			}
		}
		stateData.HintType = hintType + 1
		err := event.UpdateMessage(discord.NewMessageUpdateBuilder().
			AddActionRow(util.GetGuessButtons(stateData)...).
			Build())
		if err != nil {
			log.Error("there was an error while updating message after hint usage: ", err)
		}
		_, err = client.CreateFollowupMessage(event.ApplicationID(), event.Token(), messageBuilder.Build())
		if err != nil {
			log.Error("there was an error while creating hint message: ", err)
		}
	}
}

func onModal(event *events.ModalSubmitInteractionCreate, countryData *data.CountryData) {
	eventData := event.Data

	var stateData util.ModalStateData
	modalID := eventData.CustomID
	_ = json.Unmarshal([]byte(modalID), &stateData)

	streak := stateData.Streak
	difficulty := stateData.Difficulty
	countryIndex := stateData.SliceIndex
	country := countryData.Countries[countryIndex]
	countryInput := eventData.Text("input")
	countryInputLow := strings.TrimSpace(strings.ToLower(countryInput))
	newCountryData := &util.NewCountryData{
		Interaction:   event,
		Difficulty:    difficulty,
		MinPopulation: stateData.MinPopulation,
		Ephemeral:     stateData.Ephemeral,
		SliceIndex:    countryIndex,
		Client:        event.Client().Rest(),
		CountryData:   countryData,
	}
	if countryInputLow == strings.ToLower(country.Name.Common) || countryInputLow == strings.ToLower(country.Name.Official) {
		newCountryData.FollowupContent = fmt.Sprintf("Your guess was **correct**! It was **%s**. %s", country.Name.Common, country.Flag)
		newCountryData.Streak = streak + 1
		util.SendGameUpdates(newCountryData)
	} else if difficulty == util.GameDifficultyNormal {
		err := event.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent("Your guess was **incorrect**. Please try again.").
			SetEphemeral(true).
			Build())
		if err != nil {
			log.Error("there was an error while creating a followup: ", err)
		}
	} else if difficulty == util.GameDifficultyHard {
		if streak == 0 {
			newCountryData.FollowupContent = fmt.Sprintf("Your guess was **incorrect**. It was %s. %s", country.Name.Common, country.Flag)
		} else {
			newCountryData.FollowupContent = fmt.Sprintf("Your guess was **incorrect** and you've lost your streak of **%d**! It was **%s**. %s", streak, country.Name.Common, country.Flag)
		}
		util.SendGameUpdates(newCountryData)
	}
}
