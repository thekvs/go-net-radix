package netradix

/*
#include "radix.h"
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>

void
free_node_payload(radix_node_t *node, void *cbctx)
{
	free(node->data);
}

void
destroy(radix_tree_t *tree)
{
	Destroy_Radix(tree, free_node_payload, NULL);
}

radix_node_t*
make_node(radix_tree_t *tree, char *addr, const char **errmsg)
{
	radix_node_t *node;
	void         *p;
	prefix_t     prefix;

	p = prefix_pton(addr, -1, &prefix, errmsg);
	if (p == NULL) {
		return NULL;
	}

	node = radix_lookup(tree, &prefix);

	return node;
}
*/
import "C"

import (
	"errors"
	"unsafe"
)

type NetRadixTree struct {
	tree *C.radix_tree_t
}

// NewNetRadixTree creates new NetRadixTree structure.
// Return tuple in which the first element specifies tree structure and the second
// element specifies error object or nil if no error has occured.
func NewNetRadixTree() (*NetRadixTree, error) {
	tree := &NetRadixTree{C.New_Radix()}
	if tree == nil {
		return nil, errors.New("couldn't create radix tree")
	}

	return tree, nil
}

// Close destroys radix tree.
func (rtree *NetRadixTree) Close() {
	C.destroy(rtree.tree)
}

// Add adds network or subnet specification and user defined payload string to the radix tree.
// If no mask width is specified, the longest possible mask is assumed,
// i.e. 32 bits for IPv4 network and 128 bits for IPv6 network.
// On success, returns nil, otherwise returns error object.
func (rtree *NetRadixTree) Add(addr string, udata unsafe.Pointer) error {
	cstr := C.CString(addr)
	defer C.free(unsafe.Pointer(cstr))

	var errmsg *C.char
	node := C.make_node(rtree.tree, cstr, &errmsg)
	if node == nil {
		return errors.New(C.GoString(errmsg))
	}

	node.data = udata

	return nil
}

// SearchBest searches radix tree to find a matching node using usual subnetting rules
// for the address specified. If no mask width is specified, the longest possible mask
// for this type of address (IPv4 or IPv6) is assumed.
// Returns triple in which the first element indicates success of a search,
// the second element returns user payload (or empty string if not found)
// and the third element returns error object in case such an error occured or nil otherwise.
func (rtree *NetRadixTree) SearchBest(addr string) (found bool, udata unsafe.Pointer, err error) {
	var prefix C.prefix_t
	e := rtree.fillPrefix(&prefix, addr)
	if e != nil {
		return false, nil, e
	}

	node := C.radix_search_best(rtree.tree, &prefix)
	if node != nil {
		return true, node.data, nil
	}

	return false, nil, nil
}

// SearchExact searches radix tree to find a matching node. Its semantics are the same as in SearchBest()
// method except that the addr must match a node exactly.
func (rtree *NetRadixTree) SearchExact(addr string) (found bool, udata unsafe.Pointer, err error) {
	var prefix C.prefix_t
	e := rtree.fillPrefix(&prefix, addr)
	if e != nil {
		return false, nil, e
	}

	node := C.radix_search_exact(rtree.tree, &prefix)
	if node != nil {
		return true, node.data, nil
	}

	return false, nil, nil
}

// Remove deletes a node which exactly matches the address given.
// If no errors occured returns nil or error object otherwise.
func (rtree *NetRadixTree) Remove(addr string) error {
	var prefix C.prefix_t
	err := rtree.fillPrefix(&prefix, addr)
	if err != nil {
		return err
	}

	node := C.radix_search_exact(rtree.tree, &prefix)
	if node != nil {
		//C.free(unsafe.Pointer(node.data))
		C.radix_remove(rtree.tree, node)
	}

	return nil
}

func (rtree *NetRadixTree) fillPrefix(prefix *C.prefix_t, addr string) error {
	cstr := C.CString(addr)
	defer C.free(unsafe.Pointer(cstr))

	var errmsg *C.char
	ptr := C.prefix_pton(cstr, -1, prefix, &errmsg)
	if ptr == nil {
		return errors.New(C.GoString(errmsg))
	}

	return nil
}
