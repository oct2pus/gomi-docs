package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/doc"
	"go/token"
	"log"
	"os"
	"strings"
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
	return string(buf.Next(bp.Length))
}

// GenVarBlock formats the variable and constant field
func GenVarBlock(decl *ast.GenDecl, fset *token.FileSet) string {
	return "```A Block of Code.\n" + GenCodeBlock(decl, fset) + "\n```\n" + GenComment(decl) + "\n"
}

// bytePos is a helper struct to minimize line length when determining offsets
// in a token's position
type bytePos struct {
	Start, Length, AbsEnd int
}

// getBytePos creates a bytePos from a declaration and a file token
func getBytePos(decl *ast.GenDecl, file *token.File) bytePos {
	var bp bytePos
	bp.Start = file.Offset(decl.Pos())
	bp.Length = file.Offset(decl.End()) - file.Offset(decl.Pos())
	bp.AbsEnd = file.Offset(decl.End())
	return bp
}

// tabs2blocks handles comments that have a mixture of unformatted and
// preformated text, such as the Overview which might have some blocks of code
// in it
func tabs2blocks(in string) string {
	lines := strings.Split(in, "\n")
	senil := make([]string, 0)
	tabbed, appended := false, false
	for x := range lines {
		if strings.HasPrefix(lines[x], "	") {
			if !tabbed {
				senil = append(senil, "```A block of Golang code.\n"+lines[x])
				tabbed = true
				appended = true
			}
		} else if tabbed {
			senil = append(senil, "\n```\n"+lines[x])
			tabbed = false
			appended = true
		}
		if !appended {
			senil = append(senil, lines[x])
		}
		appended = false
	}
	out := strings.Join(senil, "\n")
	return out
}
