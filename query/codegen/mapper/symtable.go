package mapper

import (
	"fmt"

	"github.com/skydb/sky/query/ast"
)

// symtable represents a scope inside the query and allows variable
// declarations to be looked up while performing codegen.
type symtable struct {
	parent *symtable
	decls map[string]*ast.VarDecl
}

func newSymtable(parent *symtable) *symtable {
	return &symtable{
		parent:parent,
		decls: make(map[string]*ast.VarDecl),
	}
}

// find looks up a declaration by name. If not found in the current scope
// then declaration is searched for higher up the scope hierarchy.
func (tbl *symtable) find(name string) *ast.VarDecl {
	if tbl.decls[name] != nil {
		return tbl.decls[name]
	} else if tbl.parent != nil {
		return tbl.parent.find(name)
	}
	return nil
}

// add creates a new entry for the declaration in the symbol table. If an
// entry already exists then an error is returned.
func (tbl *symtable) add(decl *ast.VarDecl) error {
	if tbl.decls[decl.Name] != nil {
		return fmt.Errorf("duplicate symbol in scope: %s", decl.Name)
	}
	tbl.decls[decl.Name] = decl
	return nil
}
