#include <stdlib.h>
#include "dbg.h"

#include "node.h"

//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

// Creates an AST node for a function argument declaration.
//
// var_decl - The variable declaration.
//
// Returns a function argument node.
eql_ast_node *eql_ast_farg_create(struct eql_ast_node *var_decl)
{
    eql_ast_node *node = malloc(sizeof(eql_ast_node)); check_mem(node);
    node->type = EQL_AST_TYPE_FARG;
    node->parent = NULL;
    node->line_no = node->char_no = 0;
    node->generated = false;
    node->farg.var_decl = var_decl;
    if(var_decl != NULL) var_decl->parent = node;
    
    return node;

error:
    eql_ast_node_free(node);
    return NULL;
}

// Frees a variable declaration AST node from memory.
//
// node - The AST node to free.
void eql_ast_farg_free(struct eql_ast_node *node)
{
    if(node->farg.var_decl) {
        eql_ast_node_free(node->farg.var_decl);
    }
    node->farg.var_decl = NULL;
}

// Copies a node and its children.
//
// node - The node to copy.
// ret  - A pointer to where the new copy should be returned to.
//
// Returns 0 if successful, otherwise returns -1.
int eql_ast_farg_copy(eql_ast_node *node, eql_ast_node **ret)
{
    int rc;
    check(node != NULL, "Node required");
    check(ret != NULL, "Return pointer required");

    eql_ast_node *clone = eql_ast_farg_create(NULL);
    check_mem(clone);

    rc = eql_ast_node_copy(node->farg.var_decl, &clone->farg.var_decl);
    check(rc == 0, "Unable to copy var decl");
    if(clone->farg.var_decl) clone->farg.var_decl->parent = clone;
    
    *ret = clone;
    return 0;

error:
    eql_ast_node_free(clone);
    *ret = NULL;
    return -1;
}


//--------------------------------------
// Codegen
//--------------------------------------

// Recursively generates LLVM code for the function argument AST node.
//
// node    - The node to generate an LLVM value for.
// module  - The compilation unit this node is a part of.
// value   - A pointer to where the LLVM value should be returned.
//
// Returns 0 if successful, otherwise returns -1.
int eql_ast_farg_codegen(eql_ast_node *node, eql_module *module,
                         LLVMValueRef *value)
{
    int rc;

    check(node != NULL, "Node required");
    check(node->type == EQL_AST_TYPE_FARG, "Node type must be 'function argument'");
    check(node->farg.var_decl != NULL, "Function argument declaration required");
    check(module != NULL, "Module required");

    // Delegate LLVM generation to the variable declaration.
    rc = eql_ast_node_codegen(node->farg.var_decl, module, value);
    check(rc == 0, "Unable to codegen function argument");
    
    return 0;

error:
    return -1;
}


//--------------------------------------
// Preprocessor
//--------------------------------------

// Preprocesses the node.
//
// node   - The node.
// module - The module that the node is a part of.
//
// Returns 0 if successful, otherwise returns -1.
int eql_ast_farg_preprocess(eql_ast_node *node, eql_module *module)
{
    int rc;
    check(node != NULL, "Node required");
    check(module != NULL, "Module required");
    
    // Preprocess variable declaration.
    rc = eql_ast_node_preprocess(node->farg.var_decl, module);
    check(rc == 0, "Unable to preprocess function argument variable declaration");
    
    return 0;

error:
    return -1;
}


//--------------------------------------
// Type refs
//--------------------------------------

// Computes a list of type references used by the node.
//
// node      -  The node.
// type_refs - A pointer to an array of type refs.
// count     - A pointer to where the number of type refs is stored.
//
// Returns 0 if successful, otherwise returns -1.
int eql_ast_farg_get_type_refs(eql_ast_node *node,
                               eql_ast_node ***type_refs,
                               uint32_t *count)
{
    int rc;
    check(node != NULL, "Node required");
    check(type_refs != NULL, "Type refs return pointer required");
    check(count != NULL, "Type ref count return pointer required");

    if(node->farg.var_decl != NULL) {
        rc = eql_ast_node_get_type_refs(node->farg.var_decl, type_refs, count);
        check(rc == 0, "Unable to add function argument type refs");
    }

    return 0;
    
error:
    eql_ast_node_type_refs_free(type_refs, count);
    return -1;
}


//--------------------------------------
// Dependencies
//--------------------------------------

// Computes a list of class names that this AST node depends on.
//
// node         - The node to compute dependencies for.
// dependencies - A pointer to an array of dependencies.
// count        - A pointer to where the number of dependencies is stored.
//
// Returns 0 if successful, otherwise returns -1.
int eql_ast_farg_get_dependencies(eql_ast_node *node,
                                  bstring **dependencies,
                                  uint32_t *count)
{
    int rc;
    check(node != NULL, "Node required");
    check(dependencies != NULL, "Dependencies return pointer required");
    check(count != NULL, "Dependency count return pointer required");

    if(node->farg.var_decl != NULL) {
        rc = eql_ast_node_get_dependencies(node->farg.var_decl, dependencies, count);
        check(rc == 0, "Unable to add function argument dependencies");
    }

    return 0;
    
error:
    eql_ast_node_dependencies_free(dependencies, count);
    return -1;
}


//--------------------------------------
// Validation
//--------------------------------------

// Validates the AST node.
//
// node   - The node to validate.
// module - The module that the node is a part of.
//
// Returns 0 if successful, otherwise returns -1.
int eql_ast_farg_validate(eql_ast_node *node, eql_module *module)
{
    int rc;
    check(node != NULL, "Node required");
    check(module != NULL, "Module required");
    
    // Validate variable declaration.
    rc = eql_ast_node_validate(node->farg.var_decl, module);
    check(rc == 0, "Unable to validate function argument variable declaration");
    
    return 0;

error:
    return -1;
}


//--------------------------------------
// Debugging
//--------------------------------------

// Append the contents of the AST node to the string.
// 
// node - The node to dump.
// ret  - A pointer to the bstring to concatenate to.
//
// Return 0 if successful, otherwise returns -1.s
int eql_ast_farg_dump(eql_ast_node *node, bstring ret)
{
    int rc;
    check(node != NULL, "Node required");
    check(ret != NULL, "String required");

    // Append dump.
    check(bcatcstr(ret, "<farg>\n") == BSTR_OK, "Unable to append dump");

    // Recursively dump children.
    if(node->farg.var_decl != NULL) {
        rc = eql_ast_node_dump(node->farg.var_decl, ret);
        check(rc == 0, "Unable to dump variable declaration");
    }

    return 0;

error:
    return -1;
}