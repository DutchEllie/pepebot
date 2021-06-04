package main

import (
	"github.com/bwmarrin/discordgo"
	"quenten.nl/pepebot/discord/mux"
)

func (app *application) LogToConsole(next mux.Command) mux.Command {
	fn := func(s *discordgo.Session, m *discordgo.MessageCreate) {
		app.infoLog.Printf("%s \tsaid: %s\n", m.Author.Username, m.Content)
		next.Execute(s, m)
	}
	return mux.HandlerFunc(fn)
}
