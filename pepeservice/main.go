package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger

	pepe_list []string

	pepe_dir string
}

func main() {
	pepe_dir := os.Getenv("PEPE_DIR")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	file, err := os.Open(pepe_dir)
	if err != nil {
		errorLog.Printf("Error opening pepe directory\n")
		return
	}

	pepe_list, err := file.Readdirnames(0)
	if err != nil {
		errorLog.Printf("Error reading pepe directory file names\n")
		return
	}
	file.Close()

	app := &application{
		infoLog:   infoLog,
		errorLog:  errorLog,
		pepe_dir:  pepe_dir,
		pepe_list: pepe_list,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/pepe", app.sendPepe)
	mux.HandleFunc("/reload", app.reloadList)

	app.infoLog.Printf("Starting server at :4000\n")
	err = http.ListenAndServe(":4000", mux)
	log.Fatal(err)

}

func (app *application) sendPepe(w http.ResponseWriter, r *http.Request) {
	// Random number generator
	s := rand.NewSource(time.Now().UnixMicro())
	rd := rand.New(s) // Init pseudorandom generator
	number := rd.Intn(len(app.pepe_list))

	baseURL := "https://cdn.nicecock.eu/pepe/1.00/"
	URL := baseURL + app.pepe_list[number]

	w.Write([]byte(URL))
}

func (app *application) reloadList(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open(app.pepe_dir)
	if err != nil {
		app.errorLog.Printf("Error opening pepe directory\n")
		return
	}
	defer file.Close()
	pepe_list, err := file.Readdirnames(0)
	if err != nil {
		app.errorLog.Printf("Error reading pepe directory file names\n")
		return
	}
	app.pepe_list = pepe_list

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - Reloaded the list of pepes"))
}
