/*****************************************************************************
 * router.go
 * Name: Amir Mokhammed-Ali
 * NetId: am70
 *****************************************************************************/

package http_router

import (
	"net/http"
	"strings"
	"fmt"
)

// A trie that contains different routes added to the router
type Node struct {
	next 		map[string](*Node) 	// key: directory name, value: subtrie rooted at the given directory (key "*" corresponds to a capture name)
	args		[]string		// a list of captures (:name) met along the path from the root to the current node
	handlers	map[string](http.HandlerFunc) // Stores handlers for paths ending in this node (key: method, value: handler)
					
}

func newNode() *Node {
	return &Node {
		make(map[string](*Node)),
		make([]string, 0),
		make(map[string](http.HandlerFunc)),
	}	
}

// HTTPRouter stores the information necessary to route HTTP requests
type HTTPRouter struct {
	root	*Node	// root of the trie
}

// NewRouter creates a new HTTP Router, with no initial routes
func NewRouter() *HTTPRouter {
	return &HTTPRouter {
		newNode(),
	}
}

// AddRoute adds a new route to the router, associating a given method and path
// pattern with the designated http handler.
func (router *HTTPRouter) AddRoute(method string, pattern string, handler http.HandlerFunc) {
	splitFunc := func (r rune) bool {
		return r == '/'
	}
	patternSlice := strings.FieldsFunc(pattern, splitFunc)
	// quickly resolve the edge case 
	if len(patternSlice) == 0 {
		patternSlice = append(patternSlice, "/")
	}

	
	curNode := router.root
	args := make([]string, 0)
	var key string
	for index, dir := range patternSlice {
		// capture or not
		if dir[0] == ':' {
			key = "*"
			args = append(args, dir[1:])
		} else {
			key = dir
		}
		// go down the trie
		if _, ok := curNode.next[key]; !ok {
			curNode.next[key] = newNode()
		}
		curNode = curNode.next[key]
		// add handler if curNode is the last node in the path
		if (index == len(patternSlice) - 1) {
			curNode.handlers[method] = handler
			curNode.args = args
		}
	}	

}
// ServeHTTP writes an HTTP response to the provided response writer
// by invoking the handler associated with the route that is appropriate
// for the provided request.
func (router *HTTPRouter) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	splitFunc := func (r rune) bool {
		return r == '/'
	}
	patternSlice := strings.FieldsFunc(request.URL.Path, splitFunc)	
	// quickly resolve the edge case
	if len(patternSlice) == 0 {
		patternSlice = append(patternSlice, "/")
	}
	// find the required path
	node, args := router.findRoute(router.root, patternSlice, 0)
	if node == nil {
		http.NotFound(response, request)
		return
	}

	handler, ok := node.handlers[request.Method]
	if !ok {
		http.NotFound(response, request)
		return
	}

	query := request.URL.Query()
	for index, capture := range node.args {
		query.Add(capture, args[index])
	}
	request.URL.RawQuery = query.Encode()
	handler(response, request)
}


func (router *HTTPRouter) findRoute (root *Node, patternSlice []string, depth int) (*Node, []string) {
	if depth == len(patternSlice) {
		return root, make([]string, 0)
	}

	dir := patternSlice[depth]
	args := make([]string, 0)			// for storing capture values (e.g. "Amir")
	nextNode, ok := root.next[dir]
	
	// try non-captures first
	if ok { 	
		node, nextArgs := router.findRoute(nextNode, patternSlice, depth + 1)
		if node != nil {
			args = append(args, nextArgs...)
			return node, args
		}
	}


	// try captures if no success
	nextNode, ok = root.next["*"]	
	if !ok {
		return nil, make([]string, 0) 		// report error
	}
		
	args = append(args, dir)
	node, nextArgs := router.findRoute(nextNode, patternSlice, depth + 1)
	if node == nil {
		args = args[:len(args) - 1] 		// remove previously added capture value
		return nil, make([]string, 0)
	}

	args = append(args, nextArgs...)
	fmt.Println(node)
	return node, args


	






}
