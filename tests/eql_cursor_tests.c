#include <stdio.h>
#include <stdlib.h>

#include <eql/eql.h>
#include <eql_path.h>
#include <eql_cursor.h>

#include "minunit.h"
#include "eql_test_util.h"


//==============================================================================
//
// Fixtures
//
//==============================================================================

size_t DATA_LENGTH = 57;
char DATA[] = 
    "\x0a\x00\x00\x00\x31\x00\x00\x00\x01\xa0\x00\x00\x00\x00\x00\x00"
    "\x00\x0b\x00\x02\xa1\x00\x00\x00\x00\x00\x00\x00\x05\x00\x00\x00"
    "\x01\xa3\x66\x6f\x6f\x03\xa2\x00\x00\x00\x00\x00\x00\x00\x0d\x00"
    "\x05\x00\x00\x00\x01\xa3\x62\x61\x72"
;


//==============================================================================
//
// Test Cases
//
//==============================================================================

//--------------------------------------
// Simple Cursor
//--------------------------------------

typedef int64_t (*sky_eql_path_int_func)(sky_eql_path *path);

int test_sky_eql_cursor_execute_simple() {
    eql_module *module = NULL;
    COMPILE_QUERY_1ARG(module, "Path", "path",
        "Int total = 0;\n"
        "Cursor cursor = path.events();\n"
        "for each (Event event in cursor) {\n"
        "  total = total + event.actionId;\n"
        "}\n"
        "return total;"
    );

    // Initialize path.
    sky_eql_path *path = sky_eql_path_create();
    path->path_ptr = &DATA;

    // Execute module.
    sky_eql_path_int_func f = NULL;
    eql_module_get_main_function(module, (void*)(&f));
    int64_t ret = f(path);

    // Validate that the sum of action ids is correct.
    mu_assert_int64_equals(ret, 24LL);

    // Clean up.
    sky_eql_path_free(path);
    eql_module_free(module);
    return 0;
}


//--------------------------------------
// Complex Cursor w/ Map
//--------------------------------------

struct Result {
    int64_t hash_code;
    int64_t id;
    int64_t count;
};

typedef void (*sky_eql_path_map_func)(sky_eql_path *path, eql_map *map);

int test_sky_eql_cursor_execute_with_map() {
    eql_ast_node *type_ref, *var_decl;
    uint32_t arg_count = 2;
    eql_ast_node *args[arg_count];
    
    // Path arg.
    struct tagbstring path_str = bsStatic("path");
    type_ref = eql_ast_type_ref_create_cstr("Path");
    var_decl = eql_ast_var_decl_create(type_ref, &path_str, NULL);
    args[0] = eql_ast_farg_create(var_decl);
    
    // Map arg.
    struct tagbstring data_str = bsStatic("data");
    type_ref = eql_ast_type_ref_create_cstr("Map");
    eql_ast_type_ref_add_subtype(type_ref, eql_ast_type_ref_create_cstr("Int"));
    eql_ast_type_ref_add_subtype(type_ref, eql_ast_type_ref_create_cstr("Result"));
    var_decl = eql_ast_var_decl_create(type_ref, &data_str, NULL);
    args[1] = eql_ast_farg_create(var_decl);

    eql_module *module = NULL;
    COMPILE_QUERY_RAW(module, args, arg_count,
        "[Hashable(\"id\")]\n"
        "class Result {\n"
        "  public Int id;\n"
        "  public Int count;\n"
        "}\n"
        "Cursor cursor = path.events();\n"
        "for each (Event event in cursor) {\n"
        "  Result item = data.get(event.actionId);\n"
        "  item.count = item.count + 1;\n"
        "}\n"
        "return;"
    );

    // Initialize path & map.
    sky_eql_path *path = sky_eql_path_create();
    path->path_ptr = &DATA;
    eql_map *map = eql_map_create();

    // Execute module.
    sky_eql_path_map_func f = NULL;
    eql_module_get_main_function(module, (void*)(&f));
    f(path, map);

    // Validate the contents of the map.
    struct Result *result;
    mu_assert_int64_equals(map->count, 3LL);
    result = map->elements[0];
    mu_assert_int64_equals(result->id, 0LL);
    mu_assert_int64_equals(result->count, 1LL);
    result = map->elements[1];
    mu_assert_int64_equals(result->id, 11LL);
    mu_assert_int64_equals(result->count, 1LL);
    result = map->elements[2];
    mu_assert_int64_equals(result->id, 13LL);
    mu_assert_int64_equals(result->count, 1LL);

    // Clean up.
    sky_eql_path_free(path);
    eql_map_free(map);
    eql_module_free(module);
    return 0;
}


//==============================================================================
//
// Setup
//
//==============================================================================

int all_tests() {
    mu_run_test(test_sky_eql_cursor_execute_simple);
    mu_run_test(test_sky_eql_cursor_execute_with_map);
    return 0;
}

RUN_TESTS()