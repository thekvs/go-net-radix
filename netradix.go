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

func NewNetRadixTree() (*NetRadixTree, error) {
	tree := &NetRadixTree{C.New_Radix()}
	if tree == nil {
		return nil, errors.New("couldn't create radix tree")
	}

	return tree, nil
}

func (rtree *NetRadixTree) Close() {
	C.destroy(rtree.tree)
}

func (rtree *NetRadixTree) Add(addr string, udata string) error {
	cstr := C.CString(addr)
	defer C.free(unsafe.Pointer(cstr))

	var errmsg *C.char
	node := C.make_node(rtree.tree, cstr, &errmsg)
	if node == nil {
		return errors.New(C.GoString(errmsg))
	}

	node.data = unsafe.Pointer(C.CString(udata))

	return nil
}

func (rtree *NetRadixTree) SearchBest(addr string) (found bool, udata string, err error) {
	var prefix C.prefix_t
	e := rtree.fillPrefix(&prefix, addr)
	if e != nil {
		return false, "", e
	}

	node := C.radix_search_best(rtree.tree, &prefix)
	if node != nil {
		return true, C.GoString((*C.char)(node.data)), nil
	}

	return false, "", nil
}

func (rtree *NetRadixTree) SearchExact(addr string) (found bool, udata string, err error) {
	var prefix C.prefix_t
	e := rtree.fillPrefix(&prefix, addr)
	if e != nil {
		return false, "", e
	}

	node := C.radix_search_exact(rtree.tree, &prefix)
	if node != nil {
		return true, C.GoString((*C.char)(node.data)), nil
	}

	return false, "", nil
}

func (rtree *NetRadixTree) Remove(addr string) error {
	var prefix C.prefix_t
	err := rtree.fillPrefix(&prefix, addr)
	if err != nil {
		return err
	}

	node := C.radix_search_exact(rtree.tree, &prefix)
	if node != nil {
		C.free(unsafe.Pointer(node.data))
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
