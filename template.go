package main

import (
	"go/doc"
	"go/token"
	"log"
	"os"
	"text/template"
)

var tmpl *template.Template

func init() {
	var err error
	tmpl, err = template.ParseFiles("gemini.tmpl")
	if err != nil {
		log.Fatal(err)
	}
}

func Load(pkg *doc.Package, fset *token.FileSet) {
	data := struct {
		Package   string
		Import    string
		Overview  string
		Variables string
	}{Package: pkg.Name,
		Import:    pkg.ImportPath,
		Overview:  tabs2blocks(pkg.Doc),
		Variables: GenVarBlock(pkg.Vars[0].Decl, fset)}
	file, err := os.Create("out.gmi")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	err = tmpl.Execute(file, data)
}
