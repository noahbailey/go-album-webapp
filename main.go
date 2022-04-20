package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	database "go-album-webapp/database"

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

type dbHandler struct {
	db *sql.DB
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
	dbh := dbHandler{db: db}

	log.Println("Starting webserver on localhost:8000")
	mux := mux.NewRouter()

	//Configure functions
	mux.HandleFunc("/", dbh.renderIndex)
	mux.HandleFunc("/album/{id}", dbh.deleteAlbumByID).Methods("Delete")

	//Deal with CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	})

	//Create http handler:
	handler := c.Handler(mux)
	log.Fatal(http.ListenAndServe(":8000", handler))
}

// Render the index page:
func (dbh dbHandler) renderIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		price, err := strconv.ParseFloat(r.FormValue("price"), 64)
		if err != nil {
			log.Fatal(err)
		}
		submittedData := Album{
			Title:  r.FormValue("title"),
			Artist: r.FormValue("artist"),
			Price:  price,
		}
		id, err := dbh.addAlbum(submittedData)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id)
	}
	tpl := template.Must(template.ParseFiles("templates/index.html"))
	albums := dbh.getAlbums()
	tpl.Execute(w, albums)
}

// ----------------------------
// Utility functions

//Get all albums:
func (dbh dbHandler) getAlbums() (albums []Album) {
	row, err := dbh.db.Query("SELECT * FROM album")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() {
		item := Album{}
		err := row.Scan(&item.ID, &item.Title, &item.Artist, &item.Price)
		if err != nil {
			log.Fatal(err)
		}
		albums = append(albums, item)
	}
	return
}

//Insert a new album to table...
func (dbh dbHandler) addAlbum(album Album) (int64, error) {
	result, err := dbh.db.Exec("INSERT INTO ALBUM (title, artist, price) VALUES (?, ?, ?)", album.Title, album.Artist, album.Price)
	if err != nil {
		return 0, err
	}
	//If it cannot return the ID, something bad happened
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

//Delete a record by ID
//Should probably check if the record exists first
func (dbh dbHandler) deleteAlbumByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	//Don't check if it exists, just delete from the DB
	_, err := dbh.db.Exec("DELETE FROM ALBUM WHERE ID=?", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
