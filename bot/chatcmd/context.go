package chatcmd

import (
	"github.com/a-random-lemurian/go-owot"
	"github.com/spf13/viper"
	"golang.org/x/time/rate"
)

type Context interface {
	Chat(message string)
	IsTrusted(username string) bool
	Viper() *viper.Viper
}

type Command struct {
	Function  CommandFunction
	Ratelimit *rate.Limiter
}

type CommandFunction func(Context, *owot.MessageChat)

type CommandRegistry interface {
	RegisterCommand(name string, fn CommandFunction)
}
