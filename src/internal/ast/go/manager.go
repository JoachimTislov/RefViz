package ast

import (
	"go/ast"
	"go/token"
	"log"
	"strings"

	"github.com/JoachimTislov/RefViz/internal/types"
	"golang.org/x/tools/go/packages"
)

const (
	Const     = "const"
	Var       = "var"
	Func      = "function"
	Struct    = "struct"
	Field     = "field"
	Method    = "method"
	Interface = "interface"
	Type      = "type"
)

type manager struct {
	cfg    *packages.Config
	pkgs   []*packages.Package
	pkgMap types.Packages
	cache  string
}

func NewManager(projectRootPath string) *manager {
	cfg := &packages.Config{
		Mode:  packages.LoadAllSyntax,
		Dir:   projectRootPath,
		Tests: true,
	}
	// wildcard:  "./..." - load all
	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		log.Fatalf("Failed to load packages: %v", err)
	}

	return &manager{
		cfg:    cfg,
		pkgs:   pkgs,
		pkgMap: make(types.Packages),
		cache:  "pkgs.json",
	}
}

func (m *manager) GetSymbols() {
	for _, p := range m.pkgs {

		// Not certain why or where these are loaded from, but
		// this if statement will ignore invalid packages
		if p.Name == "main" && strings.HasSuffix(p.ID, ".test") {
			continue
		}

		pkg := m.pkgMap.Add(p.Name, p.Dir)

		for _, f := range p.Syntax {
			filePath := p.Fset.Position(f.Pos()).Filename
			file := pkg.AddFile(filePath)

			for _, decl := range f.Decls {
				switch d := decl.(type) {
				case *ast.FuncDecl:
					if d.Name.Name == "_" {
						continue
					}
					calls := &calls{pkg: p}
					calls.extract(d.Body)

					file.AddSymbol(d.Name.Name, filePath, Func)
				case *ast.GenDecl:
					for _, spec := range d.Specs {
						switch s := spec.(type) {
						case *ast.TypeSpec:
							if s.Name.Name == "_" {
								continue
							}
							file.AddSymbol(s.Name.Name, filePath, typeOfSymbol(s.Type))
						case *ast.ValueSpec:
							for _, ident := range s.Names {
								if ident.Name == "_" {
									continue
								}
								kind := Var
								if d.Tok == token.CONST {
									kind = Const
								}
								file.AddSymbol(ident.Name, filePath, kind)
							}
						}
					}
				}
			}
		}
	}
}

func typeOfSymbol(expr ast.Expr) string {
	switch expr.(type) {
	case *ast.StructType:
		return Struct
	case *ast.InterfaceType:
		return Interface
	case *ast.FuncType:
		return Func
	default:
		return Type
	}
}

func (m *manager) GetReferences() []types.Ref {
	var references []types.Ref

	for _, pkg := range m.pkgs {
		fset := pkg.Fset
		info := pkg.TypesInfo

		for _, file := range pkg.Syntax {
			ast.Inspect(file, func(n ast.Node) bool {
				ident, ok := n.(*ast.Ident)
				if !ok {
					return true
				}

				obj := info.Uses[ident]
				if obj == nil {
					return true
				}

				// Optional: Filter if you want only global references
				if obj.Pkg() != nil && obj.Pkg().Name() != pkg.PkgPath {
					return true // skip external package references, if desired
				}

				pos := fset.Position(ident.Pos())
				references = append(references, types.Ref{
					MethodName: ident.Name,
					FilePath:   pos.Filename,
				})

				return true
			})
		}
	}

	return references
}
