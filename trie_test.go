package gouldian_test

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	µ "github.com/fogfish/gouldian/v3"
	"github.com/fogfish/gouldian/v3/mock"
)

// reference implements the documented priority of the trie by brute force:
// at each segment the candidates are tried in order literal, variable, any,
// rest (depth first). Different data structure, same semantic.
func reference(patterns [][]string, segs []string, depth int, alive []int) int {
	if depth == len(segs) {
		for _, id := range alive {
			if len(patterns[id]) == depth {
				return id
			}
		}
		return -1
	}

	pick := func(match func(string) bool) []int {
		var next []int
		for _, id := range alive {
			if len(patterns[id]) > depth && match(patterns[id][depth]) {
				next = append(next, id)
			}
		}
		return next
	}

	if next := pick(func(s string) bool { return s == segs[depth] }); len(next) > 0 {
		if id := reference(patterns, segs, depth+1, next); id >= 0 {
			return id
		}
	}
	if next := pick(func(s string) bool { return s == ":" }); len(next) > 0 {
		if id := reference(patterns, segs, depth+1, next); id >= 0 {
			return id
		}
	}
	if next := pick(func(s string) bool { return s == "_" }); len(next) > 0 {
		if id := reference(patterns, segs, depth+1, next); id >= 0 {
			return id
		}
	}
	for _, id := range alive {
		if len(patterns[id]) == depth+1 && patterns[id][depth] == "*" {
			return id
		}
	}
	return -1
}

// Trie level conformance: the radix trie, the variables, the wildcards and
// the backtracking agree with the reference on random routes and paths.
func TestTrieAgreesWithReference(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	literals := []string{"a", "b", "ab", "abc", "abd", "x", "xy", "user", "users", "u"}
	kinds := []string{":", "_"}

	var hits, total int
	defer func() {
		// the test must exercise matches, not only misses
		if hits < total/5 {
			t.Errorf("only %d of %d requests are routed", hits, total)
		}
	}()

	for trial := 0; trial < 3000; trial++ {
		// random distinct patterns
		var patterns [][]string
		seen := map[string]bool{}
		for n := 1 + rnd.Intn(8); n > 0; n-- {
			size := rnd.Intn(5)
			p := make([]string, 0, size)
			for i := 0; i < size; i++ {
				switch r := rnd.Intn(10); {
				case r < 6:
					p = append(p, literals[rnd.Intn(len(literals))])
				case r < 9:
					p = append(p, kinds[rnd.Intn(len(kinds))])
				default:
					if i == size-1 {
						p = append(p, "*")
					} else {
						p = append(p, "_")
					}
				}
			}
			if key := strings.Join(p, "/"); !seen[key] {
				seen[key] = true
				patterns = append(patterns, p)
			}
		}

		routes := make([]µ.Routable, len(patterns))
		for id, p := range patterns {
			id, p := id, p
			routes[id] = func() µ.Route {
				return µ.Route{
					Method:   "GET",
					Path:     p,
					Endpoint: func(*µ.Context) error { return mock.Output(200, fmt.Sprint(id))(nil) },
				}
			}
		}
		trie := µ.NewRoutes(routes...).Endpoint()

		all := make([]int, len(patterns))
		for i := range all {
			all[i] = i
		}

		for k := 0; k < 30; k++ {
			segs := make([]string, rnd.Intn(6))
			for i := range segs {
				if rnd.Intn(4) == 0 {
					segs[i] = fmt.Sprint("v", rnd.Intn(3))
				} else {
					segs[i] = literals[rnd.Intn(len(literals))]
				}
			}
			path := "/" + strings.Join(segs, "/")

			want := reference(patterns, segs, 0, all)
			got := -1
			if out, ok := trie(mock.Input(mock.URL(path))).(*µ.Output); ok {
				fmt.Sscan(string(out.Body), &got)
			}

			total++
			if want >= 0 {
				hits++
			}
			if got != want {
				t.Fatalf("trial %d: routes %v, path %q: trie matched #%d, reference #%d", trial, patterns, path, got, want)
			}
		}
	}
}
