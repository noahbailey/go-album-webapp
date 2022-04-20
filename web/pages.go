package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

// Render the index page:
func (dbh DBHandler) RenderIndex(w http.ResponseWriter, r *http.Request) {
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


func (dbh DBHandler) RenderEdit(w http.ResponseWriter, r *http.Request) {
	editId, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		log.Fatal(err)
	}
	if r.Method == "POST" {
		//ID and Price must be converted to int and float respectively...
		price, err := strconv.ParseFloat(r.FormValue("price"), 64)
		if err != nil {
			log.Fatal(err)
		}
		submittedData := Album{
			ID: 	editId,
			Title:  r.FormValue("title"),
			Artist: r.FormValue("artist"),
			Price:  price,
		}
		res, err := dbh.updateAlbumByID(submittedData)
		if !res && err != nil {
			log.Fatal(err)
		}
	}
	tpl := template.Must(template.ParseFiles("templates/edit.html"))
	album := dbh.GetAlbumById(editId)
	tpl.Execute(w, album)
}


func (dbh DBHandler) RenderNew(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("templates/new.html"))
	albums := dbh.getAlbums()
	tpl.Execute(w, albums)	
}