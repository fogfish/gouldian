package gouldian

/*

Node of trie
*/
type Node struct {
	Path     string   // substring from the route "owned" by the node
	Heir     []*Node  // heir nodes
	Endpoint Endpoint // end point associated with node
	/*
		TODO
		- Wild     *Node    // special node that captures any path
		- Type     int      // Node type
	*/
}

/*

Walk ...
*/
func Walk(root *Node, f func(int, *Node)) {
	walk(root, 0, f)
}

func walk(node *Node, level int, f func(int, *Node)) {
	f(level, node)
	for _, n := range node.Heir {
		walk(n, level+1, f)
	}
}

var stack = make([]string, 0, 20)

func RGet(root *Node, path string) (int, *Node) {
	i := 0
	node := root
	stack = stack[:0]
	// heir := root.Heir
walk:
	for {
		if len(node.Heir) == 0 {
			return i, node
		}

		for _, n := range node.Heir {
			if len(n.Path) == 0 {
				continue
			}

			// fmt.Println("> ", n.Path, path)
			if len(path[i:]) < len(n.Path) {
				continue
			}

			if path[i] != n.Path[0] {
				continue
			}

			if len(n.Path) == 2 && n.Path[1] == ':' {
				// wild card skip segment
				p := 1
				max := len(path[i:])
				for p < max && path[i+p] != '/' {
					p++
				}
				stack = append(stack, path[i+1:i+p])
				i = i + p
				node = n
				continue walk
			}

			// fmt.Println("> ", n.Path, path)

			if path[i:i+len(n.Path)] == n.Path {
				i = i + len(n.Path)
				node = n
				// heir = n.Heir
				continue walk
			}
		}

		return i, node
	}

}

/*

appendEndpoint to trie under the path.
Input path is a collection of segments, each segment is either path literal or
wildcard symbol `:` reserved for lenses
*/
func (root *Node) appendEndpoint(path []string, endpoint Endpoint) {
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
			if node.Endpoint == nil {
				node.Endpoint = endpoint
			} else {
				node.Endpoint = node.Endpoint.Or(endpoint)
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
walk:
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
				continue walk
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
