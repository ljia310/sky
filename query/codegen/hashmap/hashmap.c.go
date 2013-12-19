package hashmap

/*
#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <inttypes.h>

#define HASHMAP_BUCKET_COUNT 256

typedef enum {
    int_value_type,
    hashmap_value_type,
} sky_hashmap_elem_type_e;

typedef struct sky_hashmap sky_hashmap;

typedef struct {
  int64_t key;
  sky_hashmap_elem_type_e value_type;
  union {
    int64_t int_value;
    sky_hashmap *hashmap_value;
  } value;
} sky_hashmap_elem;

typedef struct {
    sky_hashmap_elem *elements;
    int64_t count;
} sky_hashmap_bucket;

struct sky_hashmap {
    sky_hashmap_bucket buckets[HASHMAP_BUCKET_COUNT];
};

typedef struct {
    sky_hashmap *hashmap;
    int64_t bucket_index;
    int64_t element_index;
} sky_hashmap_iterator;


//--------------------------------------
// Hashmap
//--------------------------------------

// Creates a new hashmap instance.
sky_hashmap *sky_hashmap_new()
{
    sky_hashmap *hashmap = calloc(1, sizeof(sky_hashmap));
    return hashmap;
}

// Frees a hashmap instance.
void sky_hashmap_free(sky_hashmap *hashmap)
{
    // TODO
}

// Retrieves the bucket and element in the hashmap for a key.
inline void sky_hashmap_find(sky_hashmap *hashmap, int64_t key, sky_hashmap_bucket **bucket, sky_hashmap_elem **element)
{
    int64_t i;
    int64_t hash = key % HASHMAP_BUCKET_COUNT;
    *bucket = &hashmap->buckets[hash];
    for (i=0; i<(*bucket)->count; i++) {
        if((*bucket)->elements[i].key == key) {
            *element = &((*bucket)->elements[i]);
            return;
        }
    }
    *element = NULL;
}

// Retrieves the bucket and element in the hashmap for a key.
// The element is created if it's not found.
inline void sky_hashmap_find_or_create(sky_hashmap *hashmap, int64_t key, sky_hashmap_elem_type_e value_type, sky_hashmap_elem **elem)
{
    sky_hashmap_bucket *bucket;
    sky_hashmap_find(hashmap, key, &bucket, elem);

    // Create new element if it didn't exist.
    if((*elem) == NULL) {
        bucket->count++;
        bucket->elements = realloc(bucket->elements, (sizeof(*bucket->elements) * bucket->count));
        *elem = &bucket->elements[bucket->count-1];
        (*elem)->key = key;
        (*elem)->value_type = value_type;
        (*elem)->value.int_value = 0;
    } else if((*elem)->value_type != value_type) {
        // If the existing type is different that the new one then override it.
        if((*elem)->value_type == hashmap_value_type) {
            sky_hashmap_free((*elem)->value.hashmap_value);
        }
        (*elem)->value.hashmap_value = NULL;
        (*elem)->value_type = value_type;
    }
}

// Retrieves an int value for a given key.
// Creates the key if it doesn't already exist.
int64_t sky_hashmap_get(sky_hashmap *hashmap, int64_t key)
{
    sky_hashmap_elem *elem;
    sky_hashmap_find_or_create(hashmap, key, int_value_type, &elem);
    return elem->value.int_value;
}

// Sets an int value for a given key.
void sky_hashmap_set(sky_hashmap *hashmap, int64_t key, int64_t value)
{
    sky_hashmap_elem *elem;
    sky_hashmap_find_or_create(hashmap, key, int_value_type, &elem);
    elem->value.int_value = value;
}

// Retrieves a hashmap value for a given key.
// Creates the key if it doesn't already exist.
sky_hashmap *sky_hashmap_submap(sky_hashmap *hashmap, int64_t key)
{
    sky_hashmap_elem *elem;
    sky_hashmap_find_or_create(hashmap, key, hashmap_value_type, &elem);
    if(elem->value.hashmap_value == NULL) {
        elem->value.hashmap_value = sky_hashmap_new();
    }
    return elem->value.hashmap_value;
}


//--------------------------------------
// Hashmap Iterator
//--------------------------------------

// Creates a new hashmap iterator instance.
sky_hashmap_iterator *sky_hashmap_iterator_new(sky_hashmap *hashmap)
{
    sky_hashmap_iterator *iterator = calloc(1, sizeof(sky_hashmap_iterator));
    iterator->hashmap = hashmap;
    return iterator;
}

// Frees a hashmap iterator instance.
void sky_hashmap_iterator_free(sky_hashmap_iterator *iterator)
{
    if(iterator != NULL) {
        iterator->hashmap = NULL;
        free(iterator);
    }
}

// Returns the next element in the hashmap.
// Returns an EOF flag if no more elements are available.
bool sky_hashmap_iterator_next(sky_hashmap_iterator *iterator, int64_t *key)
{
    sky_hashmap *hashmap = iterator->hashmap;
    sky_hashmap_bucket *bucket = NULL;
    while(true) {
        // If there's no more buckets then return failure.
        if(iterator->bucket_index >= HASHMAP_BUCKET_COUNT) {
            *key = 0;
            return false;
        }

        bucket = &hashmap->buckets[iterator->bucket_index];
        if(iterator->element_index < bucket->count) {
            break;
        }
        iterator->bucket_index++;
        iterator->element_index = 0;
    }

    // Return next element key in bucket.
    *key = bucket->elements[iterator->element_index].key;

    // Increment the element index.
    iterator->element_index++;

    return true;
}


//--------------------------------------
// Benchmark
//--------------------------------------

// FOR BENCHMARKING ONLY.
sky_hashmap *sky_hashmap_benchmark_get(sky_hashmap *hashmap, int64_t n)
{
    int64_t i;
    for (i=0; i<n; i++) {
        sky_hashmap_get(hashmap, i % 256);
    }
}

// FOR BENCHMARKING ONLY.
sky_hashmap *sky_hashmap_benchmark_set(sky_hashmap *hashmap, int64_t n)
{
    int64_t i;
    for (i=0; i<n; i++) {
        sky_hashmap_set(hashmap, i % 256, 100);
    }
}
*/
import "C"
import "unsafe"

// Hashmap wraps the underlying C struct.
type Hashmap struct {
    C *C.sky_hashmap
}

// New creates a new Hashmap instance.
func New() *Hashmap {
    return &Hashmap{C.sky_hashmap_new()}
}

// Retrieves the number of buckets used by the hashmap.
func BucketCount() int {
    return C.HASHMAP_BUCKET_COUNT
}

// Free releases the underlying C memory.
func (h *Hashmap) Free() {
    if h.C != nil {
        C.free(unsafe.Pointer(h.C))
        h.C = nil
    }
}

// Get retrieves the int value for a given key.
func (h *Hashmap) Get(key int64) int64 {
    return int64(C.sky_hashmap_get(h.C, C.int64_t(key)))
}

// Set sets an int value for a given key.
func (h *Hashmap) Set(key int64, value int64) {
    C.sky_hashmap_set(h.C, C.int64_t(key), C.int64_t(value))
}

// Submap retrieves the hashmap value for a given key.
func (h *Hashmap) Submap(key int64) *Hashmap {
    return &Hashmap{C.sky_hashmap_submap(h.C, C.int64_t(key))}
}



// Iterator wraps the underlying C struct.
type Iterator struct {
    C *C.sky_hashmap_iterator
}

// NewIterator creates a new HashmapIterator instance.
func NewIterator(h *Hashmap) *Iterator {
    return &Iterator{C.sky_hashmap_iterator_new(h.C)}
}

// Free releases the underlying C memory.
func (i *Iterator) Free() {
    if i.C != nil {
        C.free(unsafe.Pointer(i.C))
        i.C = nil
    }
}

// Next retrieves the next key and a success flag.
func (i *Iterator) Next() (int64, bool) {
    var key C.int64_t
    success := bool(C.sky_hashmap_iterator_next(i.C, &key))
    return int64(key), success
}



// benchmark runs get() N number of times within the C context.
func benchmarkGet(h *Hashmap, n int64) {
    C.sky_hashmap_benchmark_get(h.C, C.int64_t(n))
}

// benchmark runs set() N number of times within the C context.
func benchmarkSet(h *Hashmap, n int64) {
    C.sky_hashmap_benchmark_set(h.C, C.int64_t(n))
}
