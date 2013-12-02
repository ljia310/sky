package mapper

/*
#cgo LDFLAGS: -llmdb
#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <inttypes.h>
#include <string.h>
#include <lmdb.h>

// The number of microseconds per second.
#define USEC_PER_SEC        1000000

// A bit-mask to extract the microseconds from a Sky timestamp.
#define USEC_MASK           0xFFFFF

// The number of bits that seconds are shifted over in a timestamp.
#define SECONDS_BIT_OFFSET  20

#define badcursordata(MSG, PTR) do {\
    fprintf(stderr, "Cursor pointing at invalid raw event data [" MSG "]: %p\n", PTR); \
    cursor->next_event->timestamp = 0; \
    return false; \
} while(0)

#define debug(M, ...) fprintf(stderr, "DEBUG %s:%d: " M "\n", __FILE__, __LINE__, ##__VA_ARGS__)

typedef struct {
  char *data;
  int64_t sz;
} sky_string;

typedef struct sky_cursor sky_cursor;

typedef struct {
  bool eos;
  bool eof;
  int64_t timestamp;
} sky_event;

struct sky_cursor {
    sky_event *event;
    sky_event *next_event;

    bool blocking_eos;

    void *key_prefix;
    int64_t key_prefix_sz;
    MDB_cursor* lmdb_cursor;
};


bool sky_cursor_iter_object(sky_cursor *cursor, MDB_val *key, MDB_val *data);


// Creates a new cursor instance.
sky_cursor *sky_cursor_new(MDB_cursor *lmdb_cursor, void *key_prefix, int64_t key_prefix_sz)
{
    sky_cursor *cursor = calloc(1, sizeof(sky_cursor));
    if(cursor == NULL) debug("[malloc] Unable to allocate cursor.");

    if(cursor->key_prefix_sz > 0) {
        cursor->key_prefix = malloc(key_prefix_sz);
        if(cursor->key_prefix == NULL) debug("[malloc] Unable to allocate key prefix.");
        memcpy(cursor->key_prefix, key_prefix, key_prefix_sz);
        cursor->key_prefix_sz = key_prefix_sz;
    }

    cursor->lmdb_cursor = lmdb_cursor;

    return cursor;
}

// Removes a cursor from memory.
void sky_cursor_free(sky_cursor *cursor)
{
    if(cursor) {
        if(cursor->key_prefix != NULL) free(cursor->key_prefix);
        cursor->key_prefix = NULL;
        cursor->key_prefix_sz = 0;
        free(cursor);
    }
}

// Moves the cursor to point to the first object. If a prefix is set then
// move to the first object that with the given prefix.
bool sky_cursor_init(sky_cursor *cursor)
{
    int rc;
    MDB_val key, data;
    if(cursor->key_prefix_sz == 0) {
        if((rc = mdb_cursor_get(cursor->lmdb_cursor, &key, &data, MDB_FIRST)) != 0) {
            if(rc != MDB_NOTFOUND) fprintf(stderr, "MDB_FIRST error: %s\n", mdb_strerror(rc));
            return false;
        }

    } else {
        key.mv_data = cursor->key_prefix;
        key.mv_size = cursor->key_prefix_sz;
        if((rc = mdb_cursor_get(cursor->lmdb_cursor, &key, &data, MDB_SET_RANGE)) != 0) {
            if(rc != MDB_NOTFOUND) fprintf(stderr, "MDB_SET_RANGE error: %s\n", mdb_strerror(rc));
            return false;
        }
    }

    return sky_cursor_iter_object(cursor, &key, &data);
}

// Moves the cursor to point to the next object.
bool sky_cursor_next_object(sky_cursor *cursor)
{
    MDB_val key, data;
    int rc = mdb_cursor_get(cursor->lmdb_cursor, &key, &data, MDB_NEXT_NODUP);
    if(rc != 0) {
        if(rc != MDB_NOTFOUND) fprintf(stderr, "MDB_NEXT_NODUP error: %s\n", mdb_strerror(rc));
        return false;
    }

    return sky_cursor_iter_object(cursor, &key, &data);
}

// Sets up object after cursor has already been positioned.
bool sky_cursor_iter_object(sky_cursor *cursor, MDB_val *key, MDB_val *data)
{
    if(cursor->key_prefix_sz > 0 && (key->mv_size < cursor->key_prefix_sz || memcmp(cursor->key_prefix, key->mv_data, cursor->key_prefix_sz) != 0)) {
        return false;
    }

    // Clear the data object if set.
    cursor->blocking_eos = false;

    return true;
}

*/
import "C"

import (
	"github.com/szferi/gomdb"
	"unsafe"
)

func sky_cursor_new(lmdb_cursor *mdb.Cursor, prefix string) *C.sky_cursor {
	cstr_prefix := unsafe.Pointer(C.CString(prefix))
	defer C.free(cstr_prefix)
	return C.sky_cursor_new(lmdb_cursor.MdbCursor(), cstr_prefix, C.int64_t(len(prefix)))
}

func sky_cursor_free(cursor *C.sky_cursor) {
	C.sky_cursor_free(cursor)
}
