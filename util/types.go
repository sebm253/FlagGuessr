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
	MinPopulation   int
	Streak          int
	SliceIndex      int
	Client          rest.Rest
}

type ButtonStateData struct {
	UserID        snowflake.ID   `json:"u"`
	Difficulty    GameDifficulty `json:"d"`
	MinPopulation int            `json:"m"`
	SliceIndex    int            `json:"i"`
	Streak        int            `json:"s"`
	ActionType    ActionType     `json:"a"`
	HintType      HintType       `json:"h"`
}

type GameDifficulty int

const (
	GameDifficultyNormal GameDifficulty = iota
	GameDifficultyHard
)

type ModalStateData struct {
	Difficulty    GameDifficulty `json:"d"`
	MinPopulation int            `json:"m"`
	SliceIndex    int            `json:"i"`
	Streak        int            `json:"s"`
}

type GameStartData struct {
	User          discord.User
	Difficulty    GameDifficulty
	MinPopulation int
	Streak        int
}
