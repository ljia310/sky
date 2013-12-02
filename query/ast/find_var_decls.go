package ast

import (
	"fmt"
)

// FindVarDecls retrieves a list of all unique variable declarations within
// an AST node. An error is returned if there are variables declared with
// different types, associations or property ids.
func FindVarDecls(node Node) (VarDecls, error) {
	v := new(varDeclVisitor)
	Walk(v, node)
	return v.decls, v.err
}

// varDeclVisitor is a visitor that retrieves a list of all unique
// variable declarations within a given AST node.
type varDeclVisitor struct {
	decls VarDecls
	err   error
}

func (v *varDeclVisitor) Visit(node Node) Visitor {
	if decl, ok := node.(*VarDecl); ok {
		if v.add(decl); v.err != nil {
			return nil
		}
	}
	return v
}

// add appends the node to the list of declarations unless the variable has
// already been declared. If there is an existing declaration that doesn't
// match then an error is set on the visitor.
func (v *varDeclVisitor) add(node *VarDecl) {
	if v.err != nil {
		return
	}

	// Find existing declaration.
	for _, decl := range v.decls {
		if decl.Name != node.Name {
			continue
		}

		var err error
		if decl.Id != node.Id {
			err = fmt.Errorf("Declaration error on '%s': mismatched id: %d != %d", decl.Name, decl.Id, node.Id)
		} else if decl.DataType != node.DataType {
			err = fmt.Errorf("Declaration error on '%s': mismatched data type: %s != %s", decl.Name, decl.DataType, node.DataType)
		} else if decl.Association != node.Association {
			err = fmt.Errorf("Declaration error on '%s': mismatched association: %s != %s", decl.Name, decl.Association, node.Association)
		}
		if err != nil {
			v.err = err
			v.decls = nil
		}

		return
	}

	// Append if no declaration exists.
	v.decls = append(v.decls, node)
}