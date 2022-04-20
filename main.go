package main

import (
	"log"
	"net/http"
	"os"

	database "go-album-webapp/database"
	web "go-album-webapp/web"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
)

type Album struct {
	ID     int     `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

func main() {

	filename := "albums.db"

	//Check if the data file needs to be created
	log.Printf("Checking if file %v exists", filename)
	fileStatus, err := database.CheckDataFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	if !fileStatus {
		log.Printf("Creating %v...", filename)
		file, err := os.Create(filename)
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
		log.Printf("Created %v successfully.", filename)
	}

	//connect to the sqlite database
	db, err := database.InitDB(filename)
	if err != nil {
		log.Fatal("Could not initialize database")
	}

	//Set up the database tables & structures
	database.CreateTable(db)

	// hand off data control to the database handler
	dbh := web.DBHandler{DB: db}

	log.Println("Starting webserver on localhost:8000")
	mux := mux.NewRouter()

	//Configure functions
	mux.HandleFunc("/", dbh.RenderIndex)
	mux.HandleFunc("/new", dbh.RenderNew)
	mux.HandleFunc("/edit", dbh.RenderEdit)
	mux.HandleFunc("/album/{id}", dbh.DeleteAlbumByID).Methods("Delete")

	//Deal with CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	})

	//Create http handler:
	handler := c.Handler(mux)
	log.Fatal(http.ListenAndServe(":8000", handler))
}
