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
	statement, _ := 
		db.Prepare("CREATE TABLE IF NOT EXISTS packages (name TEXT PRIMARY KEY, description TEXT, depends TEXT, dependants TEXT)")
	statement.Exec()
	
	statement, _ = db.Prepare("CREATE TABLE IF NOT EXISTS dependencies (package TEXT, dependency TEXT)")
	statement.Exec()

	controllers.Database = db

	r := mux.NewRouter()
	r.HandleFunc("/api", controllers.Root)
	r.HandleFunc("/api/index", controllers.Index)
	r.HandleFunc("/api/packages/{package}", controllers.Package)

	r.HandleFunc("/tmpl/index", controllers.IndexTemplate)
	r.HandleFunc("/tmpl/packages/{package}", controllers.PackageTemplate)

	r.PathPrefix("/index").Handler(http.StripPrefix("/index", http.FileServer(http.Dir("client/build/"))))
	r.PathPrefix("/packages/{package}").Handler(http.StripPrefix("/packages/{package}", http.FileServer(http.Dir("client/build/"))))

	fmt.Println("Listening on: 3000")

	http.ListenAndServe(":3000", r)
}