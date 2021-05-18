package main

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"runtime/debug"

	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
)

/* --------		DB Helper functions		-------- */

func openDB(dsn string) (*sql.DB, error){
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if orr := db.Ping(); orr != nil {
		return nil, orr
	}

	return db, nil
}

func (app *application) updateAllBadWords() (error) {
	var err error
	app.allBadWords, err = app.badwords.AllWords()
	if err != nil {
		return err
	}

	return nil
}

/* --------		Error Helper functions		-------- */

func (app *application) unknownError(err error, s *discordgo.Session, notifyDiscord bool, channelID string) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
    app.errorLog.Output(2, trace)

	if notifyDiscord {
		msg := fmt.Sprintf("An unknown error occured, error message attached below. Stack trace is in the server logs.\n%s", err.Error())
		s.ChannelMessageSend(channelID, msg)
	}
}

/* --------		Discord Helper functions		-------- */

func (app *application) readAuthToken() (string, error) {
	file, err := os.Open("./discordtoken.txt")
	if err != nil {
		return "", err
	}

	token, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(token), nil
}