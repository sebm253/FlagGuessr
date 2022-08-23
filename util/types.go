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
	Delete     Action = "delete"
	Details    Action = "details"
)

type HintType int

const (
	Population HintType = iota
	Tlds
	Capitals
)

type NewCountryData struct {
	Interaction     discord.BaseInteraction
	User            discord.User
	EmbedBuilder    discord.EmbedBuilder
	FollowupContent string
	Streak          int
	Cca             string
	Client          rest.Rest
}
