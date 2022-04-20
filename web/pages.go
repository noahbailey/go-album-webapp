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
