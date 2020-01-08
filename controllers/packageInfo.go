package controllers

import (
	"fmt"
)

type Depends struct {
	Package				string
	Link				int
}

type PackageInfo struct {
	Name	 			string
	Description 		string
	Depends			 	[]Depends
	Dependants 			[]string
}

func (pkg PackageInfo) print() {
	fmt.Printf("Package: %s\n", pkg.Name)
	var dep []string
	for _, d := range pkg.Depends {
		dep = append(dep, d.Package)
	}
	fmt.Printf("Depends: %s\n", dep)
	fmt.Printf("Dependants: %s\n", pkg.Dependants)
	fmt.Printf("Description: %s\n", pkg.Description)
}

func (pkg PackageInfo) insertToDatabase() {
	statement, _ := Database.Prepare("INSERT INTO packages (name, description) VALUES (?, ?)")
	statement.Exec(pkg.Name, pkg.Description)
}