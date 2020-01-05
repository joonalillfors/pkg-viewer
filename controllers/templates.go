package controllers

import (
	"html/template"
	"net/http"
	"github.com/gorilla/mux"
	"log"
	"fmt"
)

var IndexTemplate = func (w http.ResponseWriter, r *http.Request) {
	rows, err := Database.Query("SELECT name FROM packages ORDER BY name ASC;")
	if err != nil {
		log.Fatal(err)
	}
	var data []string
	var name string
	for rows.Next() {
		rows.Scan(&name)
		data = append(data, name)
		//fmt.Fprintf(w, fmt.Sprintf("\n%s", name))
	}
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, data)
}

var PackageTemplate = func (w http.ResponseWriter, r *http.Request) {
	packageName := mux.Vars(r)["package"]
	rows, err := Database.Query("SELECT name, description FROM packages WHERE name = ($1)", packageName)
	if err != nil {
		log.Fatal(err)
	}
	var depends, dependant string
	pkg := PackageInfo{}

	if rows.Next() {
		rows.Scan(&pkg.Name, &pkg.Description)
	} else {
		fmt.Fprintf(w, "No such package exists")
	}

	rows, err = Database.Query("SELECT package FROM dependencies WHERE dependency = ($1)", packageName)
	for rows.Next() {
		rows.Scan(&dependant)
		pkg.Dependants = append(pkg.Dependants, dependant)
	}

	rows, err = Database.Query("SELECT dependency FROM dependencies WHERE package = ($1)", packageName)
	for rows.Next() {
		rows.Scan(&depends)
		pkg.Depends = append(pkg.Depends, depends)
	}

	tmpl := template.Must(template.ParseFiles("templates/package.html"))
	tmpl.Execute(w, pkg)
}
