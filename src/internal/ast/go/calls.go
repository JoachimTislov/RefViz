package ast

import (
	"go/ast"
	goTypes "go/types"

	"golang.org/x/tools/go/packages"
)

type calls struct {
	pkg      *packages.Package
	External []string
	Local    []string
}

func (c *calls) extract(body *ast.BlockStmt) {
	ast.Inspect(body, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		return c.switchNode(call)
	})
}

func (c *calls) switchNode(call *ast.CallExpr) bool {
	switch expr := call.Fun.(type) {
	case *ast.Ident:
		c.Add(expr)

	case *ast.SelectorExpr:
		c.Add(expr.Sel)

	case *ast.FuncLit:
		c.extract(expr.Body)

	case *ast.CallExpr:
		c.switchNode(expr)
	}
	return true
}

func (c *calls) Add(ident *ast.Ident) {
	obj := c.pkg.TypesInfo.ObjectOf(ident)
	if c.isLocal(obj) {
		c.Local = append(c.Local, obj.Name())
	} else {
		c.External = append(c.External, obj.Name())
	}
}

func (c *calls) isLocal(obj goTypes.Object) bool {
	return obj != nil && obj.Pkg() == c.pkg.Types
}
