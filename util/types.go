package util

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
)

type Action string

const (
	Guess      Action = "guess"
	NewCountry Action = "new"
	Hint       Action = "hint"
	Details    Action = "details"
)

type HintType int

const (
	Population HintType = iota
	Tlds
	Capitals
	Unknown
)

type NewCountryData struct {
	Interaction     discord.BaseInteraction
	User            discord.User
	FollowupContent string
	Streak          int
	Cca             string
	Client          rest.Rest
}
