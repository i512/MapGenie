package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
	"log"
	"os"
	"regexp"
)

func main() {
	// Many tools pass their command-line arguments (after any flags)
	// uninterpreted to packages.Load so that it can interpret them
	// according to the conventions of the underlying build system.
	//var conf packages.Config
	//conf.Mode = packages.NeedTypes |
	//	packages.NeedTypesSizes |
	//	packages.NeedSyntax |
	//	packages.NeedTypesInfo |
	//	packages.NeedImports |
	//	packages.NeedName |
	//	packages.NeedFiles |
	//	packages.NeedCompiledGoFiles
	//pkgs, err := packages.Load(&conf, "./...")
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "load: %v\n", err)
	//	os.Exit(1)
	//}
	//if packages.PrintErrors(pkgs) > 0 {
	//	os.Exit(1)
	//}
	//
	//// Print the names of the source files
	//// for each package listed on the command line.
	//for _, pkg := range pkgs {
	//	fmt.Println(pkg.ID, pkg.GoFiles)
	//	for _, s := range pkg.Syntax {
	//		ast.Print(s.)
	//	}
	//}

	fset := token.NewFileSet()
	code, err := os.ReadFile("main.go")
	if err != nil {
		panic(err)
	}
	file, err := parser.ParseFile(fset, "main.go", code, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	//ast.Print(fset, file)

	regex := regexp.MustCompile(`^(?:\w+) map this pls`)

	astutil.Apply(file, nil, func(cursor *astutil.Cursor) bool {
		switch x := cursor.Node().(type) {
		case *ast.Comment:
			//fmt.Println(x.Text)
		case *ast.FuncDecl:
			if x.Doc == nil || !regex.MatchString(x.Doc.Text()) {
				fmt.Println("ignored", x.Name)
				return true
			}

			fmt.Println("will map this func:", x.Name)

			paramList := x.Type.Params.List
			if len(paramList) > 0 {
				fmt.Println("func param types:")
				for _, field := range paramList {
					switch t := field.Type.(type) {
					case *ast.Ident:
						fmt.Println(t.Name)
					default:
						fmt.Println("Unknown")
						ast.Print(fset, t)
					}
				}
				//ast.Print(fset, x.Type.Params)
			}

			results := x.Type.Results
			if results != nil && len(results.List) > 0 {
				fmt.Println("func return types:")
				for _, field := range results.List {
					switch t := field.Type.(type) {
					case *ast.Ident:
						fmt.Println(t.Name)
					default:
						fmt.Println("Unknown")
						ast.Print(fset, t)
					}
				}

			}

			//ast.Print(fset, x)

			return true
		}
		return true
	})

	//printer.Fprint(os.Stdout, fset, file)
}

type From struct {
	A	int
	B	string
	c	int
	D	[]byte
	E	[]int
	F	float64
}

type To struct {
	Other	int
	A	int
	B	string
	D	[]byte
	E	[]int
	F	float64
}

// Bluh bluh
// 1
// 2
// 3
func a() {

}

// b Bluh
func b(To) From {
	return From{}
}

// FromToTo map this pls
func FromToTo(a From) (b To) {
	var result To

	result.A = a.A
	result.B = a.B
	result.D = a.D
	result.E = a.E
	result.F = a.F

	return result
}

// FromFromToTo map this pls
func FromFromToTo(a To) (b From) {
	var result From

	result.D = a.D
	result.E = a.E
	result.F = a.F
	result.A = a.A
	result.B = a.B

	return result
}

func c(From) To {
	return To{}
}
