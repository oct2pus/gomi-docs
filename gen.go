package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/doc"
	"go/token"
	"log"
	"os"
)

//debug command
func lens(pkg *doc.Package) {
	fmt.Printf("==Lens==\nImports: %v\nFilenames: %v\nNotes: %v\nBugs: %v\nConsts: %v\nTypes: %v\nVars: %v\nFuncs: %v\nExamples: %v\n", len(pkg.Imports), len(pkg.Filenames), len(pkg.Notes), len(pkg.Bugs), len(pkg.Consts), len(pkg.Types), len(pkg.Vars[0].Names), len(pkg.Funcs), len(pkg.Examples))
}

// GenComment gets the comments from a GenDecl, returns an empty string if no
// comments exist.
func GenComment(decl *ast.GenDecl) string {
	if decl.Doc != nil {
		return decl.Doc.Text()
	}
	return ""
}

// GenCodeBlock grabs the chunk of code from the source file that matches the
// GenDecl from the AST
func GenCodeBlock(decl *ast.GenDecl, fset *token.FileSet) string {
	pos := fset.Position(decl.TokPos)
	file, err := os.ReadFile(pos.Filename)
	if err != nil {
		log.Fatal(err)
	}
	bp := getBytePos(decl, fset.File(decl.TokPos))
	buf := bytes.NewBuffer(file)
	buf.Next(bp.Start)
	fmt.Printf("reading %v to %v from file %v\n", bp.Start, bp.AbsEnd, pos.Filename)
	return string(buf.Next(bp.Length))
}

type bytePos struct {
	Start, Length, AbsEnd int
}

func getBytePos(decl *ast.GenDecl, file *token.File) bytePos {
	var bp bytePos
	bp.Start = file.Offset(decl.Pos())
	bp.Length = file.Offset(decl.End()) - file.Offset(decl.Pos())
	bp.AbsEnd = file.Offset(decl.End())
	return bp
}
