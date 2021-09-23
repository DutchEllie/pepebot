package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func (app *application) getPepeLink() (string, error) {
	url := "http://" + app.pepeServer + "/pepe"
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (app *application) sendPepe(s *discordgo.Session, m *discordgo.MessageCreate) {
	url, err := app.getPepeLink()
	if err != nil {
		app.errorLog.Print(err)
		return
	}

	_, err = s.ChannelMessageSend(m.ChannelID, url)
	if err != nil {
		app.errorLog.Print(err)
		return
	}
}

func (app *application) sendCringe(s *discordgo.Session, m *discordgo.MessageCreate) {
	_, err := s.ChannelMessageSend(m.ChannelID, "https://cdn.nicecock.eu/cringe.webm")
	if err != nil {
		app.errorLog.Print(err)
		return
	}
}

func (app *application) sendNigelGif(s *discordgo.Session, m *discordgo.MessageCreate) {
	msg := "<@77516941199159296> kun je die gif verwijderen van die pickup truck die naar de camera rijdt want bij mij zorg ie voor dat discord opnieuw opstart. ik weet niet of iemand anders dit heeft maar als iemand weet hoe dit komt en een andere oplossing weet hoor ik het graag."
	_, err := s.ChannelMessageSend(m.ChannelID, msg)
	if err != nil {
		app.errorLog.Print(err)
		return
	}
}

func (app *application) sendMonday(s *discordgo.Session, m *discordgo.MessageCreate) {
	if time.Now().Weekday().String() != "Monday" {
		_, err := s.ChannelMessageSend(m.ChannelID, "This command only works on mondays")
		if err != nil {
			app.errorLog.Print(err)
			return
		}
		return
	}
	_, err := s.ChannelMessageSend(m.ChannelID, "https://www.youtube.com/watch?v=EkALyaMjoXw")
	if err != nil {
		app.errorLog.Print(err)
		return
	}
}

func (app *application) sendTuesday(s *discordgo.Session, m *discordgo.MessageCreate) {
	if time.Now().Weekday().String() != "Tuesday" {
		_, err := s.ChannelMessageSend(m.ChannelID, "This command only works on tuesdays")
		if err != nil {
			app.errorLog.Print(err)
			return
		}
		return
	}
	_, err := s.ChannelMessageSend(m.ChannelID, "https://cdn.nicecock.eu/TBT.webm")
	if err != nil {
		app.errorLog.Print(err)
		return
	}
}

func (app *application) sendWednesday(s *discordgo.Session, m *discordgo.MessageCreate) {
	_, month, day := time.Now().Date()
	if month.String() == "May" && day == 19 {
		_, err := s.ChannelMessageSend(m.ChannelID, "https://www.youtube.com/watch?v=z21HOwUk5oM")
		if err != nil {
			app.errorLog.Print(err)
			return
		}
		return
	}
	if time.Now().Weekday().String() != "Wednesday" {
		_, err := s.ChannelMessageSend(m.ChannelID, "This command only works on wednesdays")
		if err != nil {
			app.errorLog.Print(err)
			return
		}
		return
	}
	_, err := s.ChannelMessageSend(m.ChannelID, "https://www.youtube.com/watch?v=du-TY1GUFGk")
	if err != nil {
		app.errorLog.Print(err)
		return
	}
}

func (app *application) sendGithub(s *discordgo.Session, m *discordgo.MessageCreate) {
	_, err := s.ChannelMessageSend(m.ChannelID, "My code is hosted publicly over at https://github.com/DutchEllie/pepebot")
	if err != nil {
		app.errorLog.Print(err)
		return
	}
}

func (app *application) sendPeski(s *discordgo.Session, m *discordgo.MessageCreate) {
	_, err := s.ChannelMessageSend(m.ChannelID, "https://www.youtube.com/watch?v=P0jHTCJYm44")
	if err != nil {
		app.errorLog.Print(err)
		return
	}
}

func (app *application) sendManyPepes(s *discordgo.Session, m *discordgo.MessageCreate, splitCommand []string) {
	override := false

	/* [0] is !pepe, [1] is spam, [2] is amount, [3] is override*/

	if len(splitCommand) <= 2 {
		app.errorLog.Printf("spam command had no numeral argument")
		s.ChannelMessageSend(m.ChannelID, "This command requires a numeral as a second argument, which is between 1 and 3")
		return
	}

	if len(splitCommand) > 3 {
		if splitCommand[3] == "override" {
			/* Check if admin */
			r, err := app.checkIfAdmin(s, m)
			if err != nil {
				app.errorLog.Print(err)
				return
			}
			if r {
				//s.ChannelMessageSend(m.ChannelID, "You have to be admin to override, not overriding")
				override = true
			}

		}
	}

	val, err := strconv.Atoi(splitCommand[2])
	if err != nil {
		app.errorLog.Printf("spam command had a non-numeral as argument")
		s.ChannelMessageSend(m.ChannelID, "This command requires a numeral as a second argument")
		return
	}

	if (val <= 0 || val > 3) && !override {
		s.ChannelMessageSend(m.ChannelID, "The amount has to be > 0 and < 4")
		return
	} else if val <= 0 && override {
		s.ChannelMessageSend(m.ChannelID, "I know you're admin and all, but you still have to provide a positive integer amount of pepes to send...")
		return
	}

	app.active = true

	var msg string = ""
	for i := 0; i < val; i++ {
		if app.stop {
			app.stop = false
			break
		}
		link, err := app.getPepeLink()
		if err != nil {
			app.errorLog.Print(err)
			return
		}

		if len(msg+link) > 512 {
			s.ChannelMessageSend(m.ChannelID, msg)
			msg = ""
			time.Sleep(time.Millisecond * 500)
		}

		msg += link
		msg += "\n"
	}

	s.ChannelMessageSend(m.ChannelID, msg)

	app.active = false
}

func (app *application) stopRequest(s *discordgo.Session, m *discordgo.MessageCreate) {
	if app.active {
		app.stop = true
		s.ChannelMessageSend(m.ChannelID, "Emergency stop called, hopefully I stop now")
	} else {
		app.stop = false
		s.ChannelMessageSend(m.ChannelID, "But I wasn't doing anything!")
	}

}

func (app *application) findTrigger(s *discordgo.Session, m *discordgo.MessageCreate) {
	/* Finding for every word in the allBadWords map of string slices
	Check if the message contains that word
	if it doesn't continue,
	if it does then get the word from the database, update the new time, format a message and send it */
	for i := 0; i < len(app.allBadWords[m.GuildID]); i++ {
		if strings.Contains(strings.ToLower(m.Content), strings.ToLower(app.allBadWords[m.GuildID][i])) {
			/* Found the bad word */
			word, err := app.badwords.GetWord(strings.ToLower(app.allBadWords[m.GuildID][i]), m.GuildID)
			if err != nil {
				app.errorLog.Print(err)
				s.ChannelMessageSend(m.ChannelID, err.Error())
			}
			format := formatTimeCheck(word.LastSaid)
			_, err = app.badwords.UpdateLastSaid(word.Word, word.ServerID)
			if err != nil {
				app.errorLog.Print(err)
				return
			}
			user := m.Author.Mention()
			eyesEmoji := ":eyes:"
			message := fmt.Sprintf("%s mentioned the forbidden word '%s'. They broke a streak of %s...\nYou better watch out, I am always watching %s", user, word.Word, format, eyesEmoji)
			_, err = s.ChannelMessageSend(m.ChannelID, message)
			if err != nil {
				app.errorLog.Print(err)
				return
			}
			break
		}
	}
}

func formatTimeCheck(last time.Time) string {
	now := time.Now()
	sinceLast := now.Sub(last)
	var realSeconds uint64 = uint64(sinceLast.Seconds())
	var seconds, minutes, hours, days uint64
	realBackup := realSeconds
	days = realSeconds / (24 * 3600)
	realSeconds -= days * (24 * 3600)
	hours = realSeconds / 3600
	realSeconds -= hours * 3600
	minutes = realSeconds / 60
	realSeconds -= minutes * 60
	seconds = realSeconds
	if realBackup < 60 {
		if seconds == 1 {
			return fmt.Sprintf("%d second", seconds)
		}
		return fmt.Sprintf("%d seconds", seconds)
	} else if realBackup > 60 && realBackup < 3600 {
		if seconds == 1 && minutes == 1 {
			return fmt.Sprintf("%d minute and %d second", minutes, seconds)
		} else if minutes == 1 && seconds != 1 {
			return fmt.Sprintf("%d minute and %d seconds", minutes, seconds)
		} else if minutes != 1 && seconds == 1 {
			return fmt.Sprintf("%d minutes and %d second", minutes, seconds)
		}
		return fmt.Sprintf("%d minutes and %d seconds", minutes, seconds)
	} else if realBackup > 60 && realBackup < (24*3600) {
		if hours == 1 && minutes == 1 && seconds == 1 {
			return fmt.Sprintf("%d hour, %d minute and %d second", hours, minutes, seconds)
		} else if hours == 1 && minutes == 1 && seconds != 1 {
			return fmt.Sprintf("%d hour, %d minute and %d seconds", hours, minutes, seconds)
		} else if hours == 1 && minutes != 1 && seconds == 1 {
			return fmt.Sprintf("%d hour, %d minutes and %d second", hours, minutes, seconds)
		} else if hours == 1 && minutes != 1 && seconds != 1 {
			return fmt.Sprintf("%d hour, %d minutes and %d seconds", hours, minutes, seconds)
		} else if hours != 1 && minutes == 1 && seconds == 1 {
			return fmt.Sprintf("%d hours, %d minute and %d second", hours, minutes, seconds)
		} else if hours != 1 && minutes == 1 && seconds != 1 {
			return fmt.Sprintf("%d hours, %d minute and %d seconds", hours, minutes, seconds)
		} else if hours != 1 && minutes != 1 && seconds == 1 {
			return fmt.Sprintf("%d hours, %d minutes and %d second", hours, minutes, seconds)
		} else if hours != 1 && minutes != 1 && seconds != 1 {
			return fmt.Sprintf("%d hours, %d minutes and %d seconds", hours, minutes, seconds)
		}
		return fmt.Sprintf("%d hours, %d minutes and %d seconds", hours, minutes, seconds)
	} else if realBackup > (24 * 3600) {
		if days != 1 && hours != 1 && minutes != 1 && seconds != 1 {
			return fmt.Sprintf("%d days, %d hours, %d minutes and %d seconds", days, hours, minutes, seconds)
		} else if days != 1 && hours != 1 && minutes != 1 && seconds == 1 {
			return fmt.Sprintf("%d days, %d hours, %d minutes and %d second", days, hours, minutes, seconds)
		} else if days != 1 && hours != 1 && minutes == 1 && seconds != 1 {
			return fmt.Sprintf("%d days, %d hours, %d minute and %d seconds", days, hours, minutes, seconds)
		} else if days != 1 && hours != 1 && minutes == 1 && seconds == 1 {
			return fmt.Sprintf("%d days, %d hours, %d minute and %d second", days, hours, minutes, seconds)
		} else if days != 1 && hours == 1 && minutes != 1 && seconds != 1 {
			return fmt.Sprintf("%d days, %d hour, %d minutes and %d seconds", days, hours, minutes, seconds)
		} else if days != 1 && hours == 1 && minutes != 1 && seconds == 1 {
			return fmt.Sprintf("%d days, %d hour, %d minutes and %d second", days, hours, minutes, seconds)
		} else if days != 1 && hours == 1 && minutes == 1 && seconds != 1 {
			return fmt.Sprintf("%d days, %d hour, %d minute and %d seconds", days, hours, minutes, seconds)
		} else if days != 1 && hours == 1 && minutes == 1 && seconds == 1 {
			return fmt.Sprintf("%d days, %d hour, %d minute and %d second", days, hours, minutes, seconds)
		} else if days == 1 && hours != 1 && minutes != 1 && seconds != 1 {
			return fmt.Sprintf("%d day, %d hours, %d minutes and %d seconds", days, hours, minutes, seconds)
		} else if days == 1 && hours != 1 && minutes != 1 && seconds == 1 {
			return fmt.Sprintf("%d day, %d hours, %d minutes and %d second", days, hours, minutes, seconds)
		} else if days == 1 && hours != 1 && minutes == 1 && seconds != 1 {
			return fmt.Sprintf("%d day, %d hours, %d minute and %d seconds", days, hours, minutes, seconds)
		} else if days == 1 && hours != 1 && minutes == 1 && seconds == 1 {
			return fmt.Sprintf("%d day, %d hours, %d minute and %d second", days, hours, minutes, seconds)
		} else if days == 1 && hours == 1 && minutes != 1 && seconds != 1 {
			return fmt.Sprintf("%d day, %d hour, %d minutes and %d seconds", days, hours, minutes, seconds)
		} else if days == 1 && hours == 1 && minutes != 1 && seconds == 1 {
			return fmt.Sprintf("%d day, %d hour, %d minutes and %d second", days, hours, minutes, seconds)
		} else if days == 1 && hours == 1 && minutes == 1 && seconds != 1 {
			return fmt.Sprintf("%d day, %d hour, %d minute and %d seconds", days, hours, minutes, seconds)
		} else if days == 1 && hours == 1 && minutes == 1 && seconds == 1 {
			return fmt.Sprintf("%d day, %d hour, %d minute and %d second", days, hours, minutes, seconds)
		}
		return fmt.Sprintf("%d days, %d hours, %d minutes and %d seconds", days, hours, minutes, seconds)
	}
	return "error"
}

func (app *application) checkIfAdmin(s *discordgo.Session, m *discordgo.MessageCreate) (bool, error) {
	authorMemberInfo, err := s.GuildMember(m.GuildID, m.Author.ID)
	if err != nil {
		return false, err
	}

	roleIDs, err := app.adminroles.GetAdminRoleIDs()
	if err != nil {
		return false, err
	}

	for i := 0; i < len(authorMemberInfo.Roles); i++ {
		for j := 0; j < len(roleIDs); j++ {
			if authorMemberInfo.Roles[i] == roleIDs[j] {
				return true, nil
			}
		}
	}

	app.infoLog.Printf("The user %s tried to perform an admin command without an admin role, purge them", m.Author)
	_, err = s.ChannelMessageSend(m.ChannelID, "You aren't authorized to perform this function, this incident has been reported.")
	if err != nil {
		return false, err
	}

	return false, nil
}

func (app *application) contextLength(splitCommand []string) error {
	if !(len(splitCommand) > 2) {
		app.errorLog.Printf("The command's context was not enough.\n")
		return errors.New("not enough context")
	}
	return nil
}

func (app *application) successMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	_, err := s.ChannelMessageSend(m.ChannelID, "Success!")
	if err != nil {
		app.unknownError(err, s, true, m.ChannelID)
		return
	}
}
