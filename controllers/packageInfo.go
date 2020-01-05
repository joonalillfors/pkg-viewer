package controllers

import (
	"fmt"
)

type PackageInfo struct {
	Name	 			string
	Description 		string
	Depends			 	[]string
	Dependants 			[]string
}

func (pkg PackageInfo) print() {
	fmt.Printf("Package: %s\n", pkg.Name)
	fmt.Printf("Depends: %s\n", pkg.Depends)
	fmt.Printf("Dependants: %s\n", pkg.Dependants)
	fmt.Printf("Description: %s\n", pkg.Description)
}

func (pkg PackageInfo) insertToDatabase() {
	statement, _ := Database.Prepare("INSERT INTO packages (name, description, depends, dependants) VALUES (?, ?, ?, ?)")
	var depends, dependants string

	if len(pkg.Depends) > 0 {
		depends = pkg.Depends[0]
		for _, d := range pkg.Depends[1:] {
			depends += fmt.Sprintf(", %s", d)
		}
	} else {
		depends = ""
	}

	if len(pkg.Dependants) > 0 {
		dependants = pkg.Dependants[0]
		for _, d := range pkg.Dependants[1:] {
			dependants += fmt.Sprintf(", %s", d)
		}
	} else {
		dependants = ""
	}
	
	statement.Exec(pkg.Name, pkg.Description, depends, dependants)
}