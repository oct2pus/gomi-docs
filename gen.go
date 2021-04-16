package main

import (
	"fmt"
	"go/doc"
)

//debug command
func lens(pkg *doc.Package) {
	fmt.Printf("==Lens==\nImports: %v\nFilenames: %v\nNotes: %v\nBugs: %v\nConsts: %v\nTypes: %v\nVars: %v\nFuncs: %v\nExamples: %v\n", len(pkg.Imports), len(pkg.Filenames), len(pkg.Notes), len(pkg.Bugs), len(pkg.Consts), len(pkg.Types), len(pkg.Vars[0].Names), len(pkg.Funcs), len(pkg.Examples))
}

func genVars(pkg *doc.Package) {
	fmt.Println("Vars:")
	for i, ele := range pkg.Vars {
		fmt.Printf("%v:%v\n", ele.Names[i], ele.Doc)
	}
}
