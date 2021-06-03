package main

import "github.com/bwmarrin/discordgo"

func (app *application) LogToConsole(next Command) Command {
	fn := func(s *discordgo.Session, m *discordgo.MessageCreate) {
		app.infoLog.Printf("%s   \tsaid: %s\n", m.Author.Username, m.Content)
		next.Execute(s, m)
	}
	return HandlerFunc(fn)
}
