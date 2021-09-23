package main

import (
	"net/http"

	"github.com/bwmarrin/discordgo"
)

func (app *application) addWord(s *discordgo.Session, m *discordgo.MessageCreate, splitCommand []string) {
	/* Check if admin */
	r, err := app.checkIfAdmin(s, m)
	if err != nil {
		app.errorLog.Print(err)
		return
	}
	if !r {
		return
	}
	/* [0] = trigger, [1] is addword, [2] is the word! */
	err = app.contextLength(splitCommand)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Please provide a word to add")
		return
	}

	_, err = app.badwords.InsertNewWord(splitCommand[2], m.GuildID)
	if err != nil {
		app.errorLog.Print(err)
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	err = app.updateAllBadWords()
	if err != nil {
		app.unknownError(err, s, true, m.ChannelID)
		return
	}

	app.successMessage(s, m)
}

func (app *application) removeWord(s *discordgo.Session, m *discordgo.MessageCreate, splitCommand []string) {
	/* Check if admin */
	r, err := app.checkIfAdmin(s, m)
	if err != nil {
		app.errorLog.Print(err)
		return
	}
	if !r {
		return
	}
	/* [0] = trigger, [1] is removeword, [2] is the word! */
	err = app.contextLength(splitCommand)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Please provide a word to remove")
		return
	}

	err = app.badwords.RemoveWord(splitCommand[2], m.GuildID)
	if err != nil {
		app.errorLog.Print(err)
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	err = app.updateAllBadWords()
	if err != nil {
		app.unknownError(err, s, true, m.ChannelID)
		return
	}

	app.successMessage(s, m)
}

func (app *application) addAdmin(s *discordgo.Session, m *discordgo.MessageCreate, splitCommand []string) {
	/* Check if admin */
	r, err := app.checkIfAdmin(s, m)
	if err != nil {
		app.errorLog.Print(err)
		return
	}
	if !r {
		return
	}
	/* [0] = trigger, [1] is addadmin, [2] is the id! */
	err = app.contextLength(splitCommand)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Please provide a role id")
		return
	}

	allRoles, err := s.GuildRoles(m.GuildID)
	if err != nil {
		app.unknownError(err, s, true, m.ChannelID)
		return
	}

	var found bool = false
	var counter int = 0
	for i := 0; i < len(allRoles); i++ {
		if allRoles[i].ID == splitCommand[2] {
			found = true
			counter = i
			break
		}
	}

	if !found {
		s.ChannelMessageSend(m.ChannelID, "This role id does not exist")
		return
	}

	_, err = app.adminroles.AddAdminRole(allRoles[counter].Name, allRoles[counter].ID, m.GuildID)
	if err != nil {
		app.unknownError(err, s, true, m.ChannelID)
		return
	}

	app.successMessage(s, m)

}

func (app *application) removeAdmin(s *discordgo.Session, m *discordgo.MessageCreate, splitCommand []string) {
	/* Check if admin */
	r, err := app.checkIfAdmin(s, m)
	if err != nil {
		app.errorLog.Print(err)
		return
	}
	if !r {
		return
	}
	/* [0] = trigger, [1] is removeadmin, [2] is the id! */
	err = app.contextLength(splitCommand)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Please provide a role id")
		return
	}

	allRoles, err := s.GuildRoles(m.GuildID)
	if err != nil {
		app.unknownError(err, s, true, m.ChannelID)
		return
	}

	var found bool = false
	var counter int = 0
	for i := 0; i < len(allRoles); i++ {
		if allRoles[i].ID == splitCommand[2] {
			found = true
			counter = i
			break
		}
	}

	if !found {
		s.ChannelMessageSend(m.ChannelID, "This role id does not exist")
		return
	}

	err = app.adminroles.RemoveAdminRole(allRoles[counter].Name, allRoles[counter].ID, m.GuildID)
	if err != nil {
		app.unknownError(err, s, true, m.ChannelID)
		return
	}

	app.successMessage(s, m)

}

func (app *application) reloadPepeList(s *discordgo.Session, m *discordgo.MessageCreate) {
	/* Check if admin */
	r, err := app.checkIfAdmin(s, m)
	if err != nil {
		app.errorLog.Print(err)
		return
	}
	if !r {
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Reloading list of pepes")
	url := "http://" + app.pepeServer + "/reload"
	_, err = http.Get(url)
	if err != nil {
		app.errorLog.Print(err)
		s.ChannelMessageSend(m.ChannelID, "An error occured!")
		return
	}

	app.successMessage(s, m)
}
