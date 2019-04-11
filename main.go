package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	database *sql.DB
	logger   *log.Logger
	port     string
	dbFile   string
)

func cacheResponseFor(w http.ResponseWriter, r *http.Request, seconds int) {
	cacheUntil := time.Now().UTC().Add(time.Duration(seconds) * time.Second).Format(http.TimeFormat)
	w.Header().Set("Cache-Control", "public, max-age="+strconv.Itoa(seconds))
	w.Header().Set("Expires", cacheUntil)
}

func dontCacheResponse(w http.ResponseWriter, r *http.Request) {
	cacheUntil := time.Now().UTC().Add(time.Duration(-3600) * time.Second).Format(http.TimeFormat)
	w.Header().Set("Cache-Control", "no-cache, must-revalidate")
	w.Header().Set("Expires", cacheUntil)
}

func init() {
	logger = log.New(os.Stderr, "DRBL: ", log.LstdFlags|log.Lshortfile)
	logger.Println("Starting up")
	flag.StringVar(&port, "port", "8080", "Port for the web service to listen")
	flag.StringVar(&dbFile, "db", "blacklist.db", "Path for sqlite3 db file")
	flag.Parse()
}

func main() {

	var err error

	database, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		panic(err)
	}

	defer database.Close()

	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS blacklist (domain TEXT NOT NULL UNIQUE) ")
	if err != nil {
		panic(err)
	}
	statement.Exec()
	logger.Println("Database File CONNECTED:", dbFile)

	m := http.NewServeMux()

	// All URLs will be handled by this function
	m.HandleFunc("/insert/", insert)
	m.HandleFunc("/search/", search)
	m.HandleFunc("/delete/", delete)

	m.HandleFunc("/whitelist/", delete)
	m.HandleFunc("/blacklist/", insert)
	m.HandleFunc("/test/", search)
	m.HandleFunc("/check/", search)

	logger.Println("Starting Web Service... on PORT", port)
	logger.Fatal(http.ListenAndServe(":8080", m))
}
