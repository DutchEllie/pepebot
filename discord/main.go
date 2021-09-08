package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"quenten.nl/pepebot/limiter"
	"quenten.nl/pepebot/models/mysql"
)

/* Application struct contains the logging objects.
It also has many methods for the different functions of the bot.
These methods are mostly located in discord.go */
type application struct {
	errorLog    *log.Logger
	infoLog     *log.Logger
	badwords    *mysql.BadwordModel
	adminroles  *mysql.AdminRolesModel
	trigger     string
	allBadWords map[string][]string
	pepeServer  string
	limiter     *limiter.Limiter

	active bool
	stop   bool
}

func main() {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	discordToken := os.Getenv("DISCORD_TOKEN")
	rateLimit := os.Getenv("RATE_LIMIT")
	timeLimit := os.Getenv("TIME_LIMIT")
	pepeServer := os.Getenv("PEPE_SERVER")
	dsn := fmt.Sprintf("%s:%s@tcp(db:3306)/badwords?parseTime=true", dbUser, dbPass)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	rateLim, err := strconv.Atoi(rateLimit)
	if err != nil {
		errorLog.Fatal(err)
	}

	timeLim, err := strconv.Atoi(timeLimit)
	if err != nil {
		errorLog.Fatal(err)
	}

	limiter := &limiter.Limiter{
		RateLimit: rateLim,
		TimeLimit: time.Duration(timeLim * int(time.Second)),
		Logs:      make(map[string][]*limiter.Action),
	}

	app := &application{
		infoLog:    infoLog,
		errorLog:   errorLog,
		badwords:   &mysql.BadwordModel{DB: db},
		adminroles: &mysql.AdminRolesModel{DB: db},
		trigger:    "!pepe",
		limiter:    limiter,
		pepeServer: pepeServer,
	}

	app.allBadWords, err = app.badwords.AllWords()
	if err != nil {
		app.errorLog.Fatal(err)
	}

	/* token, err := app.readAuthToken()
	if err != nil {
		app.errorLog.Fatal(err)
	} */

	discord, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		app.errorLog.Fatal(err)
	}

	discord.AddHandler(app.messageCreate)

	discord.Identify.Intents = discordgo.IntentsGuildMessages

	err = discord.Open()
	if err != nil {
		app.errorLog.Fatal(err)
	}
	defer discord.Close()

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()
}
