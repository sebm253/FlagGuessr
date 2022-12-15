package util

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/snowflake/v2"
)

type ActionType int

const (
	ActionTypeGuess ActionType = iota
	ActionTypeNewCountry
	ActionTypeHint
	ActionTypeDelete
	ActionTypeDetails
)

type HintType int

const (
	HintTypePopulation HintType = iota
	HintTypeTlds
	HintTypeCapitals
	HintTypeUnknown
)

type NewCountryData struct {
	Interaction     discord.BaseInteraction
	User            discord.User
	FollowupContent string
	Difficulty      GameDifficulty
	Ephemeral       bool
	Streak          int
	Cca             string
	Client          rest.Rest
}

type ButtonStateData struct {
	UserID     snowflake.ID   `json:"u"`
	Difficulty GameDifficulty `json:"d"`
	Cca        string         `json:"c"`
	Streak     int            `json:"s"`
	ActionType ActionType     `json:"a"`
	HintType   HintType       `json:"h"`
}

type GameDifficulty int

const (
	GameDifficultyNormal GameDifficulty = iota
	GameDifficultyHard
)

type ModalStateData struct {
	Difficulty GameDifficulty `json:"d"`
	Cca        string         `json:"c"`
	Streak     int            `json:"s"`
}
