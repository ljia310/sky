#ifndef _sky_eql_cursor_h
#define _sky_eql_cursor_h

#include <inttypes.h>

#include "cursor.h"
#include "eql_event.h"


//==============================================================================
//
// Definitions
//
//==============================================================================

// The cursor iterates over events in a path.
typedef struct {
    sky_cursor *cursor;
} sky_eql_cursor;


//==============================================================================
//
// Functions
//
//==============================================================================

//======================================
// Lifecycle
//======================================

sky_eql_cursor *sky_eql_cursor_create();

void sky_eql_cursor_free(sky_eql_cursor *cursor);


//======================================
// Iteration
//======================================

void sky_eql_cursor_next(sky_eql_cursor *cursor, sky_eql_event *event);

bool sky_eql_cursor_eof(sky_eql_cursor *cursor);


#endif