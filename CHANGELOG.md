## [2.1.0](https://github.com/Quad4-Software/tagparser) (2026-04-18)

Quad4 fork: module path, layout, toolchain, CI, performance, and tests.

This is the first **git tag** whose `go.mod` declares `github.com/Quad4-Software/tagparser/v2`, so consumers can depend on **`v2.1.0`** with plain `go get`.

### Module and imports

- Module path is `github.com/Quad4-Software/tagparser/v2`.
- Import the API as `github.com/Quad4-Software/tagparser/v2/pkg/tagparser`.

### Toolchain

- Go **1.26.2**.

### Layout and repository

- Library sources moved to **`pkg/tagparser/`**; **`internal/parser`** and **`internal`** unchanged in role but updated import paths.
- Added **`.gitea/workflows`** (`ci.yml`, `scan.yml`) and **`scripts/ci/`** (setup-go, gosec, govulncheck, trivy, test-all, scan-all) aligned with other Quad4 Go libraries.
- Removed **`.travis.yml`**.
- **`Makefile`**: `go vet ./...` and tests against `./...`.
- **`.gitignore`**: local `GOMODCACHE` / `GOCACHE` / `GOTMPDIR` cache dirs when building on constrained disks.

### Documentation and legal

- **`README.md`**: fork notice, Quad4 maintenance, install/import paths, layout, testing instructions (fuzz / bench / `-tags=appengine`).
- **`LICENSE`**: retained upstream copyright; added Quad4 fork line.
- **`pkg/tagparser/doc.go`**: package overview.

### Performance

- **Segment buffer reuse**: `tagParser.buf` is reused across key/value/quoted segments via `b := p.buf[:0]` and `p.buf = b` writebacks at split points; each segment converts to a fresh string with `string(b)` before the next segment overwrites the backing array (no aliasing).
- **Struct shrink**: removed the per-call `key` field (now passed as a parameter to `parseValue` / `parseQuotedValue`) and reordered fields so `tagParser` stays in the same 64-byte size class even with the new `buf` field.
- **`HasOption`**: short-circuits when `Options` is `nil`.

### Testing

- **Unit**: `tagparser_test.go` (table cases), `tagparser_more_test.go` (`HasOption` + option values), `invariant_test.go` (Parse never `nil`, idempotent across all table cases).
- **Examples**: `example_test.go`.
- **Property-based** via [pbt](https://github.com/Quad4-Software/pbt): `pbt_test.go` (no panic, determinism, `HasOption` matches `Options`).
- **Fuzz**: `FuzzParse`, `FuzzDeterminism`, `FuzzUntrustedTag` (NULs, long runs, bidi/unicode).
- **Stress** (skipped with `-short`): concurrent parses (64 goroutines), 256 KiB inputs, 4096 repeated segments, 2000-deep parens.
- **Internal**: `internal/parser/parser_test.go` (Parser API + `BenchmarkParser_ReadSep`), `internal/convert_test.go` (`StringToBytes`/`BytesToString` equivalence), `internal/unsafe_invariants_test.go` (`!appengine` `len`/`cap` invariants), `internal/fuzz_test.go` (`FuzzBytesToString`, `FuzzStringToBytes`, `FuzzConvertRoundtrip`).
- **Benchmarks**: `BenchmarkParse` sub-cases, `BenchmarkParseParallel`, `BenchmarkParseParallel_long`, `BenchmarkParseThroughput`.

### Security

- **`gosec ./...`**: 0 issues. Single `// #nosec G103` on `internal/unsafe.go` with rationale for unsafe string/byte views.
- **`govulncheck ./...`**: no vulnerabilities.
- **`internal/unsafe.go`**: `//go:build` tags for non-App Engine; safe equivalent in `internal/safe.go` for `-tags=appengine`.

### Code

- Ran **`go fix ./...`** where applicable.

## Historical upstream

Prior releases under `github.com/vmihailenco/tagparser/v2` are described in the [upstream repository](https://github.com/vmihailenco/tagparser).
