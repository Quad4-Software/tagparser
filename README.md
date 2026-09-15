# tagparser

Fork of [github.com/vmihailenco/tagparser](https://github.com/vmihailenco/tagparser), maintained by Quad4 under the module path `github.com/Quad4-Software/tagparser`. The upstream helper is small and stable. This fork updates the Go toolchain, module path, layout, CI, and tests.

The parser reads comma-separated tag names and `key:value` options, including quoted segments and parentheses inside values.

## Install

```bash
go get github.com/Quad4-Software/tagparser@latest
```

For local development against a checkout, point a replace directive at it:

```go
replace github.com/Quad4-Software/tagparser => ../tagparser
```

```go
import "github.com/Quad4-Software/tagparser/pkg/tagparser"
```

## Layout

| Path | Contents |
|------|----------|
| `pkg/tagparser` | Public API: `Parse`, `Tag`, `HasOption` |
| `internal/parser` | Incremental parser over bytes |
| `internal` | `StringToBytes` and `BytesToString`, unsafe on non-App Engine builds |
| `scripts/ci` | Local CI parity scripts: `test-all.sh`, `scan-all.sh`, `setup-go.sh` |

## Testing

- Unit tests are table-driven cases in `tagparser_test.go`, `invariant_test.go`, and `tagparser_more_test.go`. The low-level scanner has its own tests in `internal/parser/parser_test.go`.
- Property tests using pbt live in `pbt_test.go` and cover no-panic guarantees, determinism, and `HasOption` against a plain map.
- Fuzz targets are `FuzzParse` and `FuzzDeterminism` on the parser, plus `FuzzUntrustedTag` for hostile inputs (NULs, long runs, bidi and unicode). The `internal` package has `FuzzBytesToString`, `FuzzStringToBytes`, and `FuzzConvertRoundtrip` asserting the unsafe conversions match plain string and byte copy semantics.
- `internal/unsafe_invariants_test.go` checks `len` and `cap` on `StringToBytes` for non-App Engine builds. Run the safe build with `go test -tags=appengine ./internal/...`.
- `stress_test.go` covers concurrent parses, long inputs, many comma-separated segments, and deep parentheses. These are skipped under `-short`.
- Benchmarks run with `go test ./... -bench=. -benchmem -run=^$` and include `BenchmarkParse` subcases, parallel and throughput benches, and the `internal/parser` `BenchmarkParser_ReadSep`.

Example fuzz runs:

```bash
go test ./pkg/tagparser -fuzz=FuzzParse -fuzztime=30s
go test ./internal -fuzz=FuzzConvertRoundtrip -fuzztime=30s
```

## License

BSD 2-clause. See [LICENSE](LICENSE). Original copyright remains with the vmihailenco authors. Fork maintenance is attributed here.
