# Benchmarks

* [v2.txt](v2.txt) - baseline of `v2.0.9` / `v2.0.10`, `go test -run '^$' -bench . -benchmem -count 6 .`
* [v3.txt](v3.txt) - the same command on v3
* [compare.py](compare.py) - median ns/op of both runs (`benchstat` is not required)

Both files were captured on the same machine (Apple M-series, `darwin/arm64`, Go 1.27.4). Run `python3 bench/compare.py` to read the comparison.

## v3 vs v2

Median ns/op of six runs (`python3 bench/compare.py`):

| benchmark | v2 | v3 | |
| --------- | -- | -- | - |
| `LensForProduct1/5` (decoding) | 23 / 80 | 6 / 24 | 3-4x faster, the kind switch and the `unsafe` casts are gone |
| `PathParam1/5` (routing + binding) | 25 / 48 | 23 / 53 | -8% / +10% |
| `ServerParam1/5` (whole server) | 36 / 60 | 51 / 81 | +15 / +21 ns |
| `GitHub` (whole GitHub API) | 133 | 151 | +13% |
| `Endpoint1/5` (route + handler allocating a response) | 147 / 272 | 214 / 351 | +1 alloc, see below |

The gate "no regression above 10% on the hot path, or an explained one" is met with explanations, **not** by being faster everywhere. The additions are the price of behaviour that v2 lacked:

* A `nil` result is now a real `204` response. v2 produced no output at all; v3 builds server generated responses (204, 404, 405, ...) in the pooled `Context`, so they cost no second pool round trip (that saved ~7 ns per request).
* The server owns the pool of contexts and the pool of outputs, applies the issue renderer, logging and CORS (both are an early return if disabled).
* The routing table is a segment trie with backtracking instead of a character radix trie that never backtracked. A request pays for correctness: `/a/b`, `/a/:id` and `/a/_` are all tried in the priority order.
* `Output.Body` is `[]byte`, `ø.Send(string)` copies once.
* `Value{String, Number, Ptr}` is 32 bytes as in v2 (the double is stored in `Number`), so the Morphism is 48 bytes as in v2. `Ptr` carries decoded values of non-scalar types (struct, slice, `[]byte`, `TextUnmarshaler`); a first version used `any` (40 bytes) and was measurably slower.

Optimisations that were needed to get here (kept in the code): the path split without `strings.Cut`, the fast path of the lens for `string` and `int` fields, no boxing of `Lens` into an interface per request, direct `X-Request-Id` lookup in `httpd`.

Routers other than gouldian: see [routers/](routers/README.md).
