package main

import (
	"go/ast"
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
		Constants string
	}{
		Package:   pkg.Name,
		Import:    pkg.ImportPath,
		Overview:  tabs2blocks(pkg.Doc),
		Variables: "",
		Constants: ""}
	if varExists(pkg.Consts) {
		data.Constants = GenVarBlock(pkg.Consts[0].Decl, pkg.Consts[0].Doc, fset)
	}
	if varExists(pkg.Vars) {
		data.Variables = GenVarBlock(pkg.Vars[0].Decl, pkg.Vars[0].Doc, fset)
	}
	file, err := os.Create("out.gmi")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	err = tmpl.Execute(file, data)
}

func varExists(val *doc.Value) bool {
	if val != nil {
		return true
	}
	return false
}
