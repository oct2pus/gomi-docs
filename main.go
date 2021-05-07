package main

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"

	_ "git.sr.ht/~adnano/go-gemini"
)

// TODO: Remove log.Fatal(err)s
func main() {
	pkg, fset, err := getDocs("vendor/git.sr.ht/~adnano/go-gemini")
	if err != nil {
		log.Fatal(err)
	}
	//lens(pkg)
	//genVars(pkg)
	//fmt.Println(len(pkg.Types))
	/*for _, ele := range pkg.Types {
		fmt.Printf("%v\n%v\n", ele.Doc, GenCodeBlock(ele.Decl, fset))
	}*/
	//fmt.Printf("Name: %v\nDoc: %v\nImport Path: %v\n", pkg.Name, pkg.Doc, pkg.ImportPath)
	Load(pkg, fset)
}

func getDocs(dirName string) (*doc.Package, *token.FileSet, error) {
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
			//fmt.Printf("Loaded %v/%v\n", dirName, file.Name())
			fullpath := fmt.Sprintf("%v/%v", dirName, file.Name())

			bytes, err := os.ReadFile(fullpath)
			if err != nil {
				log.Fatal(err)
			}
			//fset.AddFile(fullpath, fset.Base(), len(bytes))

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

	docu, err := doc.NewFromFiles(fset, astfiles, mod, doc.AllMethods)
	return docu, fset, err
}
