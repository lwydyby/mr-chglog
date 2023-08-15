package cmd

import (
	"io"
	"strings"
)

type CLIContext struct {
	WorkingDir       string
	Stdout           io.Writer
	Stderr           io.Writer
	ConfigPath       string
	Template         string
	OutputPath       string
	Silent           bool
	NoColor          bool
	NoEmoji          bool
	NoCaseSensitive  bool
	Query            string
	NextTag          string
	TagFilterPattern string
	RepositoryURL    string
	Token            string
	Sort             string
	AI               bool
	AIType           string
	PushBot          bool
	AppID            string
	AppSecret        string
	ChatID           []string
	BotTitle         string
	Update           bool
}

func (c CLIContext) IsSingleTag() bool {
	return len(strings.Split(c.Query, "..")) == 1 && len(c.Query) != 0
}

type InitContext struct {
	WorkingDir    string
	ProjectID     string
	Token         string
	RepositoryURL string
	Stdout        io.Writer
	Stderr        io.Writer
}
