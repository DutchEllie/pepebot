package main

import "github.com/bwmarrin/discordgo"

func newCringe(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "this is a test message right from the new command system!")
}
