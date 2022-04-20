package web

import (
	"log"
	"net/http"
	"database/sql"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

//Delete a record by ID
//Should probably check if the record exists first
func (dbh DBHandler) DeleteAlbumByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	//Don't check if it exists, just delete from the DB
	_, err := dbh.DB.Exec("DELETE FROM ALBUM WHERE ID=?", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

//Get all albums:
func (dbh DBHandler) getAlbums() (albums []Album) {
	row, err := dbh.DB.Query("SELECT * FROM album")
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

func (dbh DBHandler) GetAlbumById(id int) (album Album) {
	// id := mux.Vars(r)["id"]
	var alb Album
	row := dbh.DB.QueryRow("SELECT * FROM album WHERE id = ?", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			//w.WriteHeader(http.StatusNotFound)
			log.Fatal(err)
		}
		//should something go here?
	}
	return alb
}

//Insert a new album to table...
func (dbh DBHandler) addAlbum(album Album) (int64, error) {
	result, err := dbh.DB.Exec("INSERT INTO ALBUM (title, artist, price) VALUES (?, ?, ?)", album.Title, album.Artist, album.Price)
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

func (dbh DBHandler) updateAlbumByID(album Album) (bool, error) {
	_, err := dbh.DB.Exec("UPDATE ALBUM SET title = ?, artist = ?, price = ? WHERE ID = ?", album.Title, album.Artist, album.Price, album.ID)
	if err != nil {
		return false, err
	}
	return true, nil
}
