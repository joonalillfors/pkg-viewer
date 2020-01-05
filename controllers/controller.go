package controllers

import (
	"net/http"
	"fmt"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"strings"
)

var Root = func (w http.ResponseWriter, r *http.Request) {
	statement, _ := Database.Prepare("DELETE FROM packages;")
	statement.Exec()
	content, err := ioutil.ReadFile("status")
	if err != nil {
		log.Fatal(err)
	}
	splitted := strings.Split(string(content), "\n\n")
	for _, line := range splitted {
		pkg := strings.Split(line, "\n")
		p := PackageInfo{}
		dependants := make(map[string][]string)
		for _, s := range pkg {
			if strings.HasPrefix(s, "Package:") {
				p.Name = strings.Split(s, ": ")[1]
			} else if strings.HasPrefix(s, "Description:") {
				p.Description = strings.Split(s, ": ")[1]
			} else if strings.HasPrefix(s, "Depends:") {
				depends := strings.Split(s, ": ")[1]
				depending := strings.Split(depends, ", ")
				seen := make(map[string]struct{}, len(depending))
				for _, d := range depending {
					dependency := strings.Split(d, " (")[0]
					dependants[dependency] = append(dependants[dependency], p.Name)
					// Only add distinct dependencies
					if _, ok := seen[dependency]; !ok {
						p.Depends = append(p.Depends, dependency)
						statement, _ := Database.Prepare("INSERT INTO dependencies (package, dependency) VALUES (?, ?)")
						statement.Exec(p.Name, dependency)
						seen[dependency] = struct{}{}
					}
				}
			} else {
				desc := strings.Split(s, ": ")[0]
				if strings.HasPrefix(desc, " ") {
					p.Description += fmt.Sprintf("\n%s", desc)
				}
			}
		}
		if p.Name != "" {
			p.insertToDatabase()
		}
	}
	fmt.Fprintf(w, splitted[1])
}

var Index = func (w http.ResponseWriter, r *http.Request) {
	rows, err := Database.Query("SELECT name FROM packages ORDER BY name ASC;")
	if err != nil {
		log.Fatal(err)
	}
	var res []string
	var name string
	for rows.Next() {
		rows.Scan(&name)
		res = append(res, name)
		//fmt.Fprintf(w, fmt.Sprintf("\n%s", name))
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

var Package = func (w http.ResponseWriter, r *http.Request) {
	packageName := mux.Vars(r)["package"]
	rows, err := Database.Query("SELECT name, description, depends FROM packages WHERE name = ($1)", packageName)
	if err != nil {
		log.Fatal(err)
	}
	var name, description, depends string
	if rows.Next() {
		rows.Scan(&name, &description, &depends)
		fmt.Fprintf(w, fmt.Sprintf("Package: %s\n\nDescription: %s\n\nDepends: %s\n", name, description, depends))
	} else {
		fmt.Fprintf(w, "No such package exists")
	}
}