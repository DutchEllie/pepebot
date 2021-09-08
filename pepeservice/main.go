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

	pepe_dir string
}

func main() {
	pepe_dir := os.Getenv("PEPE_DIR")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		infoLog:  infoLog,
		errorLog: errorLog,
		pepe_dir: pepe_dir,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/pepe", app.sendPepe)

	app.infoLog.Printf("Starting server at :4000\n")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)

}

func (app *application) sendPepe(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open(app.pepe_dir)
	if err != nil {
		app.errorLog.Printf("Error opening pepe directory\n")
		http.Error(w, "Internal Server Error", 500)
		return
	}
	defer file.Close()

	names, err := file.Readdirnames(0)
	if err != nil {
		app.errorLog.Printf("Error reading pepe directory file names\n")
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// Random number generator
	s := rand.NewSource(time.Now().UnixMicro())
	rd := rand.New(s) // Init pseudorandom generator
	number := rd.Intn(len(names))

	baseURL := "https://cdn.nicecock.eu/pepe/1.00/"
	URL := baseURL + names[number]

	w.Write([]byte(URL))
}
