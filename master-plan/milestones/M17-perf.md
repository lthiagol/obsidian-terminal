# M17 — Performance: Profile-Driven Optimization

**Status:** ✅ done

## Goal

Identify and optimize actual performance bottlenecks through profiling, not assumptions. Add benchmarks to prevent regressions.

## Motivation

Performance optimizations should be data-driven. Without profiling, we risk optimizing code that isn't slow while missing real bottlenecks. This milestone adds profiling infrastructure first, then optimizes only what the data shows is slow.

## Implementation Plan

### Part 1: Profiling Infrastructure

#### 1. Add benchmark tests

Create `bench_test.go` with benchmarks for hot paths:

```go
func BenchmarkFileTreeView(b *testing.B) {
    // Create tree with 1000 files
    // Benchmark View() method
}

func BenchmarkFuzzySearch(b *testing.B) {
    // Create search with 10000 paths
    // Benchmark FuzzySearch with typical query
}

func BenchmarkMarkdownRender(b *testing.B) {
    // Load large markdown file
    // Benchmark RenderMarkdown
}

func BenchmarkHelpRender(b *testing.B) {
    // Benchmark renderHelp()
}
```

#### 2. Add profiling commands

Add Makefile targets:

```makefile
bench:
	go test -bench=. -benchmem ./...

bench-cpu:
	go test -bench=. -cpuprofile=cpu.prof ./...
	go tool pprof -http=:8080 cpu.prof

bench-mem:
	go test -bench=. -memprofile=mem.prof ./...
	go tool pprof -http=:8080 mem.prof
```

#### 3. Run initial profiling

Execute benchmarks and identify top 3 bottlenecks by:
- Allocation count
- Time per operation
- Memory usage

Document findings in milestone before proceeding.

### Part 2: Targeted Optimizations

**Only implement optimizations that profiling shows are needed.**

#### Candidate Optimization A: Pre-computed tree styles

**Problem (if confirmed by profiling):**
`FileTree.View()` creates new lipgloss.Style objects per item per frame.

**Fix:**
```go
type FileTree struct {
    // existing fields ...
    fileStyle, dirStyle, selectedStyle lipgloss.Style
    prefixCache []string  // pre-computed indentation
}
```

Pre-compute in `NewFileTree`, reuse in `View()`.

#### Candidate Optimization B: Pre-computed lowercase paths

**Problem (if confirmed by profiling):**
`FuzzyScore` calls `strings.ToLower` per file per keystroke.

**Fix:**
```go
type State struct {
    // existing fields ...
    allPathsLower []string  // pre-computed
}
```

Pre-compute in `NewState`, pass to `FuzzyScore`.

**Signature changes:**
- `FuzzyScore(query, target, targetLower string) float64`
- Update all callers and tests

#### Candidate Optimization C: Help text cache

**Problem (if confirmed by profiling):**
`renderHelp()` rebuilds static text every frame.

**Fix:**
```go
var cachedHelpLines []string

func buildHelpLines() []string {
    if cachedHelpLines == nil {
        cachedHelpLines = computeHelpLines()
    }
    return cachedHelpLines
}

func InvalidateHelpCache() {
    cachedHelpLines = nil
}
```

Call `InvalidateHelpCache()` when keybindings change (if configurable).

### Part 3: Validation

#### 1. Re-run benchmarks

Compare before/after:
- Allocations per operation
- Time per operation
- Memory per operation

#### 2. Document results

Add to milestone:
```
## Profiling Results

Before optimization:
- FileTreeView: 1000 allocs/op, 5ms/op
- FuzzySearch: 500 allocs/op, 2ms/op

After optimization:
- FileTreeView: 100 allocs/op, 2ms/op (90% reduction)
- FuzzySearch: 100 allocs/op, 1ms/op (80% reduction)
```

## Decision Criteria

**Implement optimization if:**
- Profiling shows >50% of time/allocations in that function
- Optimization reduces allocations by >50%
- Optimization doesn't significantly increase code complexity

**Skip optimization if:**
- Profiling shows <10% of time/allocations
- Optimization adds significant complexity
- Code is already fast enough (<16ms per frame for 60fps)

## Testing Strategy

- Benchmark tests for all hot paths
- Unit tests still pass after optimizations
- Manual test: smooth navigation with 1000+ files
- Manual test: smooth search with 10000+ files

## Completion Criteria

- [x] Benchmark tests added for hot paths
- [x] Profiling infrastructure in Makefile
- [x] Initial profiling results documented
- [x] Only data-driven optimizations implemented
- [x] Before/after benchmark comparison documented
- [x] `make test` passes
- [x] `make vet` exits 0
- [x] `make bench` runs without errors
- [x] Manual test: smooth performance at scale

## Profiling Results

### Before optimization:
- HelpRender: 368 allocs/op, 9752 B/op
- FileTreeView: 715 allocs/op, 20040 B/op
- FuzzySearch: 20000 allocs/op (strings.ToLower per path)

### After optimization:
- HelpRender: **1 alloc/op, 704 B/op** (99.7% reduction — cached)
- FileTreeView: **614 allocs/op, 18424 B/op** (14% reduction — pre-computed styles)
- FuzzySearch: 20000 allocs/op (ToLower eliminated; remaining allocs are Result structs)
