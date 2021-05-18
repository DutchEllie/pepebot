package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func (app *application) sendPepe(s *discordgo.Session, m *discordgo.MessageCreate) {
	resp, err := http.Get("http://bbwroller.com/random")
	if err != nil {
		app.errorLog.Print(err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		app.errorLog.Print(err)
		return
	}
	rep, err := regexp.Compile("/static.*\\.jpg")
	if err != nil {
		app.errorLog.Print(err)
		return
	}
	pepes := rep.FindAllString(string(body), 200)
	if pepes == nil {
		app.errorLog.Printf("No pepes were found\n")
		return
	}
	randomIndex := rand.Intn(35)
	url := "https://bbwroller.com"
	url += pepes[randomIndex]

	_, err = s.ChannelMessageSend(m.ChannelID, url)
	if err != nil {
		app.errorLog.Print(err)
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

func (app *application) sendTuesday(s *discordgo.Session, m *discordgo.MessageCreate) {
	if time.Now().Weekday().String() != "Tuesday" {
		_, err := s.ChannelMessageSend(m.ChannelID, "This command only works on tuesdays")
		if err != nil {
			app.errorLog.Print(err)
			return
		}
	}
	_, err := s.ChannelMessageSend(m.ChannelID, "https://cdn.nicecock.eu/TBT.webm")
	if err != nil {
		app.errorLog.Print(err)
		return
	}
}

func (app *application) findTrigger(s *discordgo.Session, m *discordgo.MessageCreate) {
	/* Finding for every word in the allBadWords map of string slices
	Check if the message contains that word
	if it doesn't continue,
	if it does then get the word from the database, update the new time, format a message and send it */
	for i := 0; i < len(app.allBadWords[m.GuildID]); i++{
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
			_ ,err = s.ChannelMessageSend(m.ChannelID, message) 
			if err != nil {
				app.errorLog.Print(err)
				return
			}
			break
		}
	}
}

func formatTimeCheck(last time.Time) string{
	now := time.Now()
	sinceLast := now.Sub(last)
	var realSeconds uint64 = uint64(sinceLast.Seconds())
	var seconds, minutes, hours, days uint64
	realBackup := realSeconds
	days = realSeconds / ( 24 * 3600 )
	realSeconds -= days * ( 24 * 3600 )
	hours = realSeconds / 3600
	realSeconds -= hours * 3600
	minutes = realSeconds / 60
	realSeconds -= minutes * 60
	seconds = realSeconds
	if realBackup < 60{
		if seconds == 1{
			return fmt.Sprintf("%d second", seconds)
		}
		return fmt.Sprintf("%d seconds", seconds)
	}else if realBackup > 60 && realBackup < 3600 {
		if seconds == 1 && minutes == 1{
			return fmt.Sprintf("%d minute and %d second", minutes, seconds)
		}else if minutes == 1 && seconds != 1{
			return fmt.Sprintf("%d minute and %d seconds", minutes, seconds)
		}else if minutes != 1 && seconds == 1{
			return fmt.Sprintf("%d minutes and %d second", minutes, seconds)
		}
		return fmt.Sprintf("%d minutes and %d seconds", minutes, seconds)
	}else if realBackup > 60 && realBackup < ( 24 * 3600 ){
		if hours == 1 && minutes == 1 && seconds == 1{
			return fmt.Sprintf("%d hour, %d minute and %d second", hours, minutes, seconds)
		}else if hours == 1 && minutes == 1 && seconds != 1 {
			return fmt.Sprintf("%d hour, %d minute and %d seconds", hours, minutes, seconds)
		}else if hours == 1 && minutes != 1 && seconds == 1{
			return fmt.Sprintf("%d hour, %d minutes and %d second", hours, minutes, seconds)
		}else if hours == 1 && minutes != 1 && seconds != 1{
			return fmt.Sprintf("%d hour, %d minutes and %d seconds", hours, minutes, seconds)
		}else if hours != 1 && minutes == 1 && seconds == 1{
			return fmt.Sprintf("%d hours, %d minute and %d second", hours, minutes, seconds)
		}else if hours != 1 && minutes == 1 && seconds != 1{
			return fmt.Sprintf("%d hours, %d minute and %d seconds", hours, minutes, seconds)
		}else if hours != 1 && minutes != 1 && seconds == 1{
			return fmt.Sprintf("%d hours, %d minutes and %d second", hours, minutes, seconds)
		}else if hours != 1 && minutes != 1 && seconds != 1{
			return fmt.Sprintf("%d hours, %d minutes and %d seconds", hours, minutes, seconds)
		}
		return fmt.Sprintf("%d hours, %d minutes and %d seconds", hours, minutes, seconds)
	}else if realBackup > ( 24 * 3600 ){
		if days != 1 && hours != 1 && minutes != 1 && seconds != 1{
			return fmt.Sprintf("%d days, %d hours, %d minutes and %d seconds", days, hours, minutes, seconds)
		}else if days != 1 && hours != 1 && minutes != 1 && seconds == 1{
			return fmt.Sprintf("%d days, %d hours, %d minutes and %d second", days, hours, minutes, seconds)
		}else if days != 1 && hours != 1 && minutes == 1 && seconds != 1{
			return fmt.Sprintf("%d days, %d hours, %d minute and %d seconds", days, hours, minutes, seconds)
		}else if days != 1 && hours != 1 && minutes == 1 && seconds == 1{
			return fmt.Sprintf("%d days, %d hours, %d minute and %d second", days, hours, minutes, seconds)
		}else if days != 1 && hours == 1 && minutes != 1 && seconds != 1{
			return fmt.Sprintf("%d days, %d hour, %d minutes and %d seconds", days, hours, minutes, seconds)
		}else if days != 1 && hours == 1 && minutes != 1 && seconds == 1{
			return fmt.Sprintf("%d days, %d hour, %d minutes and %d second", days, hours, minutes, seconds)
		}else if days != 1 && hours == 1 && minutes == 1 && seconds != 1{
			return fmt.Sprintf("%d days, %d hour, %d minute and %d seconds", days, hours, minutes, seconds)
		}else if days != 1 && hours == 1 && minutes == 1 && seconds == 1{
			return fmt.Sprintf("%d days, %d hour, %d minute and %d second", days, hours, minutes, seconds)
		}else if days == 1 && hours != 1 && minutes != 1 && seconds != 1{
			return fmt.Sprintf("%d day, %d hours, %d minutes and %d seconds", days, hours, minutes, seconds)
		}else if days == 1 && hours != 1 && minutes != 1 && seconds == 1{
			return fmt.Sprintf("%d day, %d hours, %d minutes and %d second", days, hours, minutes, seconds)
		}else if days == 1 && hours != 1 && minutes == 1 && seconds != 1{
			return fmt.Sprintf("%d day, %d hours, %d minute and %d seconds", days, hours, minutes, seconds)
		}else if days == 1 && hours != 1 && minutes == 1 && seconds == 1{
			return fmt.Sprintf("%d day, %d hours, %d minute and %d second", days, hours, minutes, seconds)
		}else if days == 1 && hours == 1 && minutes != 1 && seconds != 1{
			return fmt.Sprintf("%d day, %d hour, %d minutes and %d seconds", days, hours, minutes, seconds)
		}else if days == 1 && hours == 1 && minutes != 1 && seconds == 1{
			return fmt.Sprintf("%d day, %d hour, %d minutes and %d second", days, hours, minutes, seconds)
		}else if days == 1 && hours == 1 && minutes == 1 && seconds != 1{
			return fmt.Sprintf("%d day, %d hour, %d minute and %d seconds", days, hours, minutes, seconds)
		}else if days == 1 && hours == 1 && minutes == 1 && seconds == 1{
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

func (app *application) contextLength(splitCommand []string) (error) {
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