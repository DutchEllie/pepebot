package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

/* The Command interface is a template for the implementation of a command of the discord bot.
The actual command will be executed in the Execute function.
All the actual Command objects will be (similarly to Handlers in the net/http package) put into a CommandMux */
type Command interface {
	Execute(s *discordgo.Session, m *discordgo.MessageCreate)
}

/* The CommandMux struct is a type of mux for Discord commands. It's modelled after the net/http ServeMux */
type CommandMux struct {
	m      map[string]muxEntry
	prefix string
}

func (c *CommandMux) Handler(m *discordgo.MessageCreate) (cmd Command, pattern string) {
	if strings.HasPrefix(m.Content, c.prefix) {

	}
}

func (c *CommandMux) Execute(s *discordgo.Session, m *discordgo.MessageCreate) {

}

/* The muxEntry struct contains the actual Command implementation as well as the pattern (discord command)
it will be matched against */
type muxEntry struct {
	h       Command
	pattern string
}
