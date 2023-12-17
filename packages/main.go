package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
	"os"
	"regexp"
)

func main() {
	//Many tools pass their command-line arguments (after any flags)
	//uninterpreted to packages.Load so that it can interpret them
	//according to the conventions of the underlying build system.
	fset := token.NewFileSet()

	cfg := packages.Config{
		Fset: fset,
		Mode: packages.NeedTypes |
			packages.NeedTypesSizes |
			packages.NeedSyntax |
			packages.NeedTypesInfo |
			packages.NeedImports |
			packages.NeedName |
			packages.NeedFiles |
			packages.NeedCompiledGoFiles,
	}

	pkgs, err := packages.Load(&cfg, "./...")
	if err != nil {
		fmt.Fprintf(os.Stderr, "load: %v\n", err)
		os.Exit(1)
	}
	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}

	// Print the names of the source files
	// for each package listed on the command line.
	for _, pkg := range pkgs {
		fmt.Println("analysis of pkg:", pkg.Name)

		fmt.Println(pkg.ID, pkg.GoFiles)
		for _, file := range pkg.Syntax {
			analyzePkgFile(pkg.Fset, file, pkg)
		}
	}
}

func analyzePkgFile(fset *token.FileSet, file *ast.File, pkg *packages.Package) {
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
						pkgType := pkg.Types.Scope().Lookup(t.Name)
						if s, ok := pkgType.Type().Underlying().(*types.Struct); ok {
							fmt.Println("found struct fields:")
							for i := 0; i < s.NumFields(); i++ {
								f := s.Field(i)
								fmt.Println(f.Name(), f.Exported(), f.Type().String())

								if bt, ok := f.Type().Underlying().(*types.Basic); ok {
									if (bt.Info() & types.IsInteger) == types.IsInteger {
										fmt.Println("found integer", f.Name())
									}
								}
							}
							fmt.Println(s.NumFields())
						}
						//if tn, ok := pkgType.(*types.TypeName); ok {
						//	if s, ok := tn.Type().Underlying().(*types.Struct); ok {
						//		fmt.Println("found struct fields:")
						//		for i := 0; i < s.NumFields(); i++ {
						//			f := s.Field(i)
						//			fmt.Println(f.Name(), f.Exported(), f.Type().String())
						//		}
						//		fmt.Println(s.NumFields())
						//	}
						//}
						fmt.Println(pkgType.String())

						tTypes := pkg.TypesInfo.Types[t]
						tDefs := pkg.TypesInfo.Defs[t]
						tInstances := pkg.TypesInfo.Instances[t]
						tUses := pkg.TypesInfo.Uses[t]
						tImplicits := pkg.TypesInfo.Implicits[t]
						fmt.Println(
							t.Name, "types:", tTypes, "defs:", tDefs, "instances:", tInstances, "uses:",
							tUses, "implicits:", tImplicits,
						)
						ast.Print(fset, t)
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

			return true
		}
		return true
	})
}
