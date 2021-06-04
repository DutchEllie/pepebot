package main

import (
	"github.com/bwmarrin/discordgo"
	"quenten.nl/pepebot/discord/mux"
	"quenten.nl/pepebot/limiter"
)

/*
Middleware chain

Logtoconsole -> loginteraction -> mux -> command

*/

func (app *application) LogToConsole(next mux.Command) mux.Command {
	fn := func(s *discordgo.Session, m *discordgo.MessageCreate) {
		app.infoLog.Printf("%s \tsaid: %s\n", m.Author.Username, m.Content)
		next.Execute(s, m)
	}
	return mux.HandlerFunc(fn)
}

func (app *application) LogInteraction(next mux.Command) mux.Command {
	fn := func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Logging interaction
		a := limiter.NewAction("Any message")
		app.limiter.Logs[m.Author.ID] = append(app.limiter.Logs[m.Author.ID], a)

		// Checking if rate limit exceeded
		err := app.limiter.CheckAllowed(m.Author.ID)
		if err != nil {
			mux.NotFound(s, m)
		} else {
			next.Execute(s, m)
		}
	}
	return mux.HandlerFunc(fn)
}
