package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"pkg-viewer/controllers"
	"fmt"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, _ := sql.Open("sqlite3", "./pkg.db")
	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS packages (name TEXT PRIMARY KEY, description TEXT)")
	statement.Exec()
	
	statement, _ = db.Prepare("CREATE TABLE IF NOT EXISTS dependencies (package TEXT, dependency TEXT, link INTEGER)")
	statement.Exec()

	controllers.Database = db

	r := mux.NewRouter()
	r.PathPrefix("/css").Handler(http.StripPrefix("/css", http.FileServer(http.Dir("css/"))))
	r.HandleFunc("/", controllers.Root)

	r.HandleFunc("/index", controllers.Index)
	r.HandleFunc("/packages/{package}", controllers.Package)

	// extra?
	// r.HandleFunc("/api/index", controllers.Index)
	// r.HandleFunc("/api/packages/{package}", controllers.Package)

	// r.PathPrefix("/index").Handler(http.StripPrefix("/index", http.FileServer(http.Dir("client/build/"))))
	// r.PathPrefix("/packages/{package}").Handler(http.StripPrefix("/packages/{package}", http.FileServer(http.Dir("client/build/"))))

	fmt.Println("Listening on: 3000")

	http.ListenAndServe(":3000", r)
}