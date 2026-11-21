# Router comparison

Gouldian against other Go routers on the maintained fork of the
[go-http-routing-benchmark](https://github.com/gin-gonic/go-http-routing-benchmark)
suite (GitHub, Google+ and Parse APIs, static routes), Apple M2 Pro, Go 1.27.4,
`GOMAXPROCS=1` as the suite does, 3 runs, median.

## Reproduce

```bash
git clone https://github.com/gin-gonic/go-http-routing-benchmark && cd go-http-routing-benchmark
cp <gouldian>/bench/routers/gouldian_test.go.txt gouldian_test.go
cp <gouldian>/bench/routers/gouldian_check_test.go.txt gouldian_check_test.go
go mod edit -require=github.com/fogfish/gouldian/v3@v3.0.0 -replace=github.com/fogfish/gouldian/v3=<gouldian>
go test -run TestGouldianRoutesAreRouted .      # every route of every API answers 204, never a fast 404
go test -run '^$' -bench '(Gouldian|Gin|Echo|BunRouter|HttpRouter|HttpTreeMux|Chi)_' -benchmem -count 3 -timeout 60m . > results.txt
python3 <gouldian>/bench/routers/rank.py results.txt GithubAll Param5
```

The adapter measures the library the way it is used: the full server
(`httpd.Handler`), every `:param` bound to a struct field by a lens, a handler
returning `nil` (a `204`). `ParamWrite` goes through `FMap` and `ø.Send`.

## Result (ns/op, lower is better, among Gin, Echo, BunRouter, HttpRouter, HttpTreeMux, Chi)

| Benchmark | Gouldian | rank | best | HttpRouter |
| --------- | -------: | :--: | ---- | ---------: |
| Param | 50 | 5 | BunRouter 15 | 40 |
| Param5 | 76 | 4 | Gin 50 | 126 |
| Param20 | 197 | 3 | Gin 146 | 424 |
| GithubParam | 74 | 2 | Gin 58 | 99 |
| GithubAll | 17,686 | **4** | Gin 12,374 | 18,524 |
| GPlusAll | 847 | 5 | BunRouter 458 | 844 |
| ParseAll | 1,778 | 5 | BunRouter 783 | 1,264 |
| GithubStatic | 59 | 6 | BunRouter 23 | 25 |
| StaticAll | 12,305 | 6 | HttpRouter 4,970 | 4,970 |

Gouldian makes no allocations in any routing benchmark. It is strongest on
parameterised routes and weakest on static ones.

## Radix variant

The character radix trie of v2 was ported to the v3 design (method sets,
variables as separate children, iterative descent, backtracking) and measured
against the segment trie of v3 in the same harness. It wins on static routes
(`GithubStatic` 55 vs 60 ns) and loses on parameterised ones (`GithubParam` 88
vs 75 ns, `Param20` 220 vs 206 ns, `GithubAll` 19.4 vs 18.2 µs), so the segment
trie stays. The variant is not part of the repository. The differential test `trie_test.go`
(random routes against a brute force reference) runs against either trie.
