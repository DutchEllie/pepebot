package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (app *application) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}
	app.limiter.LogInteraction(m.Author.ID, "messagecreate")

	/* Check if the user is even allowed by the rate limiter */
	err := app.limiter.CheckAllowed(m.Author.ID)
	if err != nil {
		/* normally don't send, but now do for debug purposes. This is the admin bot channel */
		//app.unknownError(err, s, true, "815952128106430514")
		app.infoLog.Printf("Rate limit exceeded by used %s", m.Author.Username)
		channel, err := s.UserChannelCreate(m.Author.ID)
		if err != nil {
			app.errorLog.Printf("Cannot create user PM channel to inform about rate limit, %s", err.Error())
			return
		}
		s.ChannelMessageSend(channel.ID, "You exceeded the rate limit for the server, please stop spamming")
		return
	}

	/* Checking if the message starts with the trigger specified in application struct
	if it does then we start the switch statement to trigger the appropriate function
	if it does not then we check if it contains a triggerword from the database */
	if strings.HasPrefix(m.Content, app.trigger) {
		splitCommand := strings.Split(strings.TrimSpace(m.Content), " ")

		/* If the whole message on it's own is just "!pepe" aka the triggerword, then send a pepe and stop */
		if strings.TrimSpace(m.Content) == app.trigger {
			app.sendPepe(s, m)
			return
		}
		/* This switches on the first word in the message
		it could be anything starting with !pepe */
		if len(splitCommand) > 1 {
			switch splitCommand[1] {
			/* --- Funny commands --- */
			case "cringe":
				app.sendCringe(s, m)
			case "gif":
				app.sendNigelGif(s, m)
			case "monday":
				app.sendMonday(s, m)
			case "tuesday":
				app.sendTuesday(s, m)
			case "wednesday":
				app.sendWednesday(s, m)
			case "github", "source":
				app.sendGithub(s, m)
			case "peski", "rotterdam":
				app.sendPeski(s, m)
			/* --- Bot commands for words --- */
			case "spam":
				app.sendManyPepes(s, m, splitCommand)
			case "stop":
				app.stopRequest(s, m)
			/* --- Bot commands, but only admins --- */
			case "addword":
				app.addWord(s, m, splitCommand)
			case "removeword":
				app.removeWord(s, m, splitCommand)
			case "addadmin":
				app.addAdmin(s, m, splitCommand)
			case "removeadmin":
				app.removeAdmin(s, m, splitCommand)
			case "reload":
				app.reloadPepeList(s, m)
			}

		}
	} else {
		/* If the trigger wasn't the prefix of the message, we need to check all the words for a trigger */
		app.findTrigger(s, m)
	}

}
