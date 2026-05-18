# tagparser (Quad4 fork)

This repository is a **fork** of [github.com/vmihailenco/tagparser](https://github.com/vmihailenco/tagparser) (v2 API), **maintained by Quad4** at `github.com/Quad4-Software/tagparser`. The upstream helper is small and stable; this fork updates the Go toolchain, module path, layout, CI, and tests.

Import the library as:

```go
import "github.com/Quad4-Software/tagparser"
```

Install:

```bash
go get github.com/Quad4-Software/tagparser/v2@latest
```

The module path is `github.com/Quad4-Software/tagparser/v2` (the `/v2` suffix matches the major version). Library sources live under **`pkg/tagparser/`**; **`internal/parser`** holds the low-level byte scanner; **`internal`** provides App Engine–safe or unsafe string/byte helpers.

## Features

- Parses comma-separated tag names and `key:value` options, including quoted segments and parentheses in values.

## Layout

| Path | Purpose |
|------|---------|
| `pkg/tagparser` | Public API (`Parse`, `Tag`, `HasOption`) |
| `internal/parser` | Incremental parser over bytes |
| `internal` | `StringToBytes` / `BytesToString` (unsafe on non-App Engine) |
| `scripts/ci` | Local CI parity with Gitea (`test-all.sh`, `scan-all.sh`, `setup-go.sh`, …) |
| `.gitea/workflows` | CI (`ci.yml`) and security scan (`scan.yml`) |

## Testing

- **Unit**: table-driven cases in `tagparser_test.go`, `invariant_test.go`, `tagparser_more_test.go`; low-level **`internal/parser`** tests in `internal/parser/parser_test.go`.
- **Property-based**: [pbt](https://github.com/Quad4-Software/pbt) in `pbt_test.go` (no panic, determinism, `HasOption` vs map).
- **Fuzz**: `FuzzParse` / `FuzzDeterminism` on the parser; `FuzzUntrustedTag` for hostile-style inputs (NULs, long runs, bidi/unicode). **`internal`**: `FuzzBytesToString`, `FuzzStringToBytes`, `FuzzConvertRoundtrip` assert unsafe (or safe) conversions match `string` / `[]byte` copy semantics. Example:  
  `go test ./pkg/tagparser -fuzz=FuzzParse -fuzztime=30s`  
  `go test ./internal -fuzz=FuzzConvertRoundtrip -fuzztime=30s`
- **Unsafe vs App Engine**: `internal/unsafe_invariants_test.go` (`!appengine`) checks `len`/`cap` on `StringToBytes`. Run the safe build with `go test -tags=appengine ./internal/...`.
- **Stress** (skipped with `-short`): concurrent parses, long inputs, many comma-separated segments, deep parentheses in `stress_test.go`.
- **Benchmarks**: `go test ./... -bench=. -benchmem -run=^$` — includes `BenchmarkParse` subcases, parallel and throughput benches, and `internal/parser` `BenchmarkParser_ReadSep`.

## License

BSD 2-clause; see [LICENSE](LICENSE). Original copyright remains with the vmihailenco authors; fork maintenance is attributed in this README.
