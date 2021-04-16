package main

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"strings"

	_ "git.sr.ht/~adnano/go-gemini"
)

// TODO: Remove log.Fatal(err)s
func main() {
	pkg, err := getDocs("vendor/git.sr.ht/~adnano/go-gemini")
	if err != nil {
		log.Fatal(err)
	}
	//lens(pkg)
	genVars(pkg)
	//fmt.Printf("Name: %v\nDoc: %v\nImport Path: %v\n", pkg.Name, pkg.Doc, pkg.ImportPath)
}

func getDocs(dirName string) (*doc.Package, error) {
	//dirName := "vendor/git.sr.ht/~adnano/go-gemini"
	fset := token.NewFileSet()
	astfiles := [](*ast.File){}
	files, err := os.ReadDir(dirName)
	var mod string
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".go") {
			fullpath := fmt.Sprintf("%v/%v", dirName, file.Name())

			bytes, err := ioutil.ReadFile(fullpath)
			if err != nil {
				log.Fatal(err)
			}
			astfile, err := parser.ParseFile(fset, fullpath, bytes, parser.ParseComments)
			if err != nil {
				log.Fatal(err)
			}
			astfiles = append(astfiles, astfile)
		} else if file.Name() == "go.mod" {
			fullpath := fmt.Sprintf("%v/%v", dirName, file.Name())
			oFile, err := os.Open(fullpath)
			if err != nil {
				log.Fatal(err)
			}
			scanner := bufio.NewScanner(oFile)
			scanner.Split(bufio.ScanLines)
			scanner.Scan()
			mod = scanner.Text()
			mod = strings.Split(mod, " ")[1]
		}
	}
	// this assumes go.mod exists
	return doc.NewFromFiles(fset, astfiles, mod, doc.AllMethods)
}
