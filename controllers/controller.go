package controllers

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"strings"
	"html/template"
	"github.com/gorilla/mux"
	"log"
)

// Read status file to sqlite db
var Root = func (w http.ResponseWriter, r *http.Request) {
	// Clear packages table
	statement, _ := Database.Prepare("DELETE FROM packages;")
	statement.Exec()

	// Clear dependencies table
	statement, _ = Database.Prepare("DELETE FROM dependencies;")
	statement.Exec()

	tmpl := template.Must(template.ParseFiles("templates/load.html"))

	// Read status file
	content, err := ioutil.ReadFile("status")
	if err != nil {
		tmpl.Execute(w, false)
		return
		//log.Fatal(err)
	}

	var packages []string

	// Split sections of package information
	splitted := strings.Split(string(content), "\n\n")

	for _, line := range splitted {
		// Split lines of a package information
		pkg := strings.Split(line, "\n")
		p := PackageInfo{}
		for _, s := range pkg {

			// Line starts with "Package:"
			if strings.HasPrefix(s, "Package:") {
				p.Name = strings.Split(s, ": ")[1]

			// Line starts with "Description:"
			} else if strings.HasPrefix(s, "Description:") {
				p.Description = strings.Split(s, ": ")[1]

			// Line starts with "Depends:"
			} else if strings.HasPrefix(s, "Depends:") {
				depends := strings.Split(s, ": ")[1]
				depending := strings.Split(depends, ", ")
				seen := make(map[string]struct{}, len(depending))
				for _, d := range depending {
					dependency := strings.Split(d, " (")[0]

					// Check if alternative dependancy exists
					dep := strings.Split(dependency, " | ")
					for _, dp := range dep {

						// Only add distinct dependencies
						if _, ok := seen[dp]; !ok {
							statement, _ := Database.Prepare("INSERT INTO dependencies (package, dependency, link) VALUES (?, ?, ?)")
							statement.Exec(p.Name, dp, 0)
							seen[dependency] = struct{}{}
						}
					}
				}
			} else {
				desc := strings.Split(s, ": ")[0]
				// Line is continuation of description
				if strings.HasPrefix(desc, " ") {
					p.Description += fmt.Sprintf("\n%s", desc)
				}
			}
		}

		// If package has a name, add to database (first round of the loop would add package with no information)
		if p.Name != "" {
			packages = append(packages, p.Name)
			p.insertToDatabase()
		}
	}

	// Update existing dependencies
	for _, pkg := range packages {
		statement, _ := Database.Prepare("UPDATE dependencies SET link = 1 WHERE dependency = (?)")
		statement.Exec(pkg)
	}

	tmpl.Execute(w, true)
}

// Index page
var Index = func (w http.ResponseWriter, r *http.Request) {
	rows, err := Database.Query("SELECT name FROM packages ORDER BY name ASC;")
	if err != nil {
		log.Fatal(err)
	}
	var data []string
	var name string
	for rows.Next() {
		rows.Scan(&name)
		data = append(data, name)
	}
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, data)
}

// Package page
var Package = func (w http.ResponseWriter, r *http.Request) {
	packageName := mux.Vars(r)["package"]
	rows, err := Database.Query("SELECT name, description FROM packages WHERE name = ($1)", packageName)
	if err != nil {
		log.Fatal(err)
	}
	var dependant string

	depends := Depends{}
	pkg := PackageInfo{}

	if rows.Next() {
		rows.Scan(&pkg.Name, &pkg.Description)
	} else {
		fmt.Fprintf(w, "No such package exists")
		return
	}

	rows, err = Database.Query("SELECT package FROM dependencies WHERE dependency = ($1)", packageName)
	for rows.Next() {
		rows.Scan(&dependant)
		pkg.Dependants = append(pkg.Dependants, dependant)
	}

	rows, err = Database.Query("SELECT dependency, link FROM dependencies WHERE package = ($1)", packageName)
	for rows.Next() {
		rows.Scan(&depends.Package, &depends.Link)
		pkg.Depends = append(pkg.Depends, depends)
	}

	tmpl := template.Must(template.ParseFiles("templates/package.html"))
	tmpl.Execute(w, pkg)
}