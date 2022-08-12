package util

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
	//FirstLetter
	//LastLetter
)
