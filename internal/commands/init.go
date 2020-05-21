package commands

import (
	"github.com/auttaja/gommand"
)

var cmds []*gommand.Command

var infoCategory = &gommand.Category{
	Name:        "Information",
	Description: "General commands to retrieve info.",
}
var modCategory = &gommand.Category{
	Name:        "Moderation",
	Description: "Commands made to ease managing and moderating servers.",
}
var utilCategory = &gommand.Category{
	Name:        "Utilities",
	Description: "Simple utilities to do with the bot.",
}

// Register registers all commands in the cmds slice.
func Register(router *gommand.Router) {
	for _, v := range cmds {
		router.SetCommand(v)
	}
}
