package gouldian

import (
	"fmt"
	"strings"
)

/*

Node of trie
*/
type Node struct {
	Path string   // substring from the route "owned" by the node
	Heir []*Node  // heir nodes
	Func Endpoint // end point associated with node
	/*
		TODO
		- Wild     *Node    // special node that captures any path
		- Type     int      // Node type
	*/
}

// NewRoutes creates new routing table
func NewRoutes(seq ...Routable) *Node {
	root := &Node{
		Heir: []*Node{
			{
				Path: "/",
				Heir: make([]*Node, 0),
			},
		},
	}

	for _, route := range seq {
		root.appendEndpoint(route())
	}

	return root
}

/*

lookup is hot-path discovery of node at the path
*/
func (root *Node) lookup(path string, values *[]string) (at int, node *Node) {
	node = root
lookup:
	for {
		// leaf node, no futher lookup is possible
		// return current `node`` and position `at` path
		if len(node.Heir) == 0 {
			return
		}

		for _, heir := range node.Heir {
			if len(path[at:]) < len(heir.Path) {
				// No match, path cannot match node
				continue
			}

			if path[at] != heir.Path[0] {
				// No match, path cannot match node
				// this is micro-optimization to reduce overhead of memequal
				continue
			}

			// the node consumers entire path
			if len(heir.Path) == 2 && heir.Path[1] == '*' {
				*values = append(*values, path[at+1:])
				at = len(path)
				node = heir
				return
			}

			if len(heir.Path) == 2 && (heir.Path[1] == ':' || heir.Path[1] == '_' || heir.Path[1] == '*') {
				// the node is a wild-card that matches any path segment
				// let's skip the path until next segment and re-call the value
				p := 1
				max := len(path[at:])
				for p < max && path[at+p] != '/' {
					p++
				}

				if heir.Path[1] == ':' {
					*values = append(*values, path[at+1:at+p])
				}

				at = at + p
				node = heir
				continue lookup
			}

			if path[at:at+len(heir.Path)] == heir.Path {
				// node matches the path, continue lookup
				at = at + len(heir.Path)
				node = heir
				continue lookup
			}
		}

		return
	}
}

/*

appendEndpoint to trie under the path.
Input path is a collection of segments, each segment is either path literal or
wildcard symbol `:` reserved for lenses
*/
func (root *Node) appendEndpoint(path []string, endpoint Endpoint) {
	if len(path) == 0 {
		_, n := root.appendTo("/")

		if n.Func == nil {
			n.Func = endpoint
		} else {
			n.Func = n.Func.Or(endpoint)
		}
		return
	}

	at := 0
	node := root
	for i, segment := range path {
		// `/` required to speed up lookup on the hot-path
		segment = "/" + segment
		at, node = node.appendTo(segment)

		// split the node and add endpoint
		if len(segment[at:]) != 0 {
			split := &Node{
				Path: segment[at:],
				Heir: make([]*Node, 0),
			}
			node.Heir = append(node.Heir, split)
			node = split
		}

		// the last segment needs to be enhanced with endpoint
		if i == len(path)-1 {
			if node.Func == nil {
				node.Func = endpoint
			} else {
				node.Func = node.Func.Or(endpoint)
			}
		}
	}
}

/*

appendTo finds the node in trie where to add path (or segment).
It returns the candidate node and length of "consumed" path
*/
func (root *Node) appendTo(path string) (at int, node *Node) {
	node = root
lookup:
	for {
		if len(node.Heir) == 0 {
			// leaf node, no futher lookup is possible
			// return current `node`` and position `at` path
			return
		}

		for _, heir := range node.Heir {
			prefix := longestCommonPrefix(path[at:], heir.Path)
			at = at + prefix
			switch {
			case prefix == 0:
				// No common prefix, jump to next heir
				continue
			case prefix == len(heir.Path):
				// Common prefix is the node itself, continue lookup into heirs
				node = heir
				continue lookup
			default:
				// Common prefix is shorter than node itself, split is required
				if prefixNode := node.heirByPath(heir.Path[:prefix]); prefixNode != nil {
					// prefix already exists, current node needs to be moved
					// under existing one
					node.Path = node.Path[prefix:]
					prefixNode.Heir = append(prefixNode.Heir, node)
					node = prefixNode
					return
				}

				// prefix does not exist, current node needs to be split
				// the list of heirs needs to be patched
				for j := 0; j < len(node.Heir); j++ {
					if node.Heir[j].Path == heir.Path {
						n := heir
						node.Heir[j] = &Node{
							Path: heir.Path[:prefix],
							Heir: []*Node{n},
						}
						n.Path = heir.Path[prefix:]
						node = node.Heir[j]
						return
					}
				}
			}
		}
		// No heir is found return current node
		return
	}
}

func (root *Node) heirByPath(path string) *Node {
	for i := 0; i < len(root.Heir); i++ {
		if root.Heir[i].Path == path {
			return root.Heir[i]
		}
	}
	return nil
}

/*

Walk through trie, use for debug purposes only
*/
func (root *Node) Walk(f func(int, *Node)) {
	walk(root, 0, f)
}

func walk(node *Node, level int, f func(int, *Node)) {
	f(level, node)
	for _, n := range node.Heir {
		walk(n, level+1, f)
	}
}

// Println outputs trie to console
func (root *Node) Println() {
	root.Walk(
		func(i int, n *Node) {
			fmt.Println(strings.Repeat(" ", i), n.Path)
		},
	)
}

// Endpoint converts trie to Endpoint
func (root *Node) Endpoint() Endpoint {
	return func(ctx *Context) (err error) {
		if ctx.Request == nil {
			return ErrNoMatch
		}

		path := ctx.Request.URL.Path
		ctx.free()

		ctx.values = ctx.values[:0]
		i, node := root.lookup(path, &ctx.values)

		if len(path) == i && node.Func != nil {
			return node.Func(ctx)
		}

		return ErrNoMatch
	}
}

//
// Utils
//

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func longestCommonPrefix(a, b string) (prefix int) {
	max := min(len(a), len(b))
	for prefix < max && a[prefix] == b[prefix] {
		prefix++
	}
	return
}
