# Architecture & Code Review Template

Use this template when performing a **planning-only** or **plan-then-execute** review of a codebase. Copy to `ARCHITECTURE-REVIEW-{YYYY-MM-DD}.md` and fill every section. Designed to be **project-agnostic** — replace `{placeholders}` with project specifics.

---

## 1. Review metadata

| Field | Value |
|-------|-------|
| **Project** | `{name}` |
| **Reviewer** | `{human or agent id}` |
| **Date** | `{YYYY-MM-DD}` |
| **Branch / commit** | `{branch}` @ `{short-sha}` |
| **Review type** | `planning-only` / `pre-release` / `post-milestone` / `ad-hoc` |
| **Scope** | `{e.g. full repo, subsystem X, milestone M12 closure}` |
| **Out of scope** | `{e.g. no code changes, no dependency upgrades}` |

---

## 2. Pre-review checklist (do before reading code)

Complete these steps and record outputs in the review doc.

- [ ] Read `{README}` — stated goals vs reality
- [ ] Read `{ARCHITECTURE or DESIGN doc}` — module map matches filesystem (`find` / glob)
- [ ] Read `{agent/contributor guide}` — conventions agents should follow
- [ ] Read `{master plan / STATUS}` — milestone status vs milestone files (spot-check 3+ “done” items)
- [ ] Run `{test command}` — pass/fail, count tests, note coverage if available
- [ ] Run `{lint/vet/static analysis}` — pass/fail
- [ ] Run `{build command}` — pass/fail
- [ ] List `{dependency manifest}` — count direct deps, flag new/unused
- [ ] Identify **hard constraints** (e.g. read-only, no network, single binary)

**Evidence block** (paste or link):

```
Tests:     {command} → {N pass}, coverage {X%}
Build:     {command} → {ok/fail}
LOC:       {approx by package/module}
Last plan update: {date on STATUS}
```

---

## 3. Review dimensions

Score each **A / B / C / D / F** with 2–4 sentences of justification. Agents should not skip a dimension because the project is “small.”

| Dimension | Grade | What to inspect |
|-----------|-------|-----------------|
| **Architecture** | | Layering, boundaries, coupling, state ownership, extension points |
| **Correctness** | | Edge cases, error paths, concurrency, state machines, idempotency |
| **Maintainability** | | File size, naming, duplication, god objects, dead code |
| **Tests** | | Meaningful assertions, integration paths, fixtures, flaky patterns |
| **Documentation** | | README, architecture doc, API docs, plan sync, keybindings/config |
| **Performance** | | Hot paths, allocations, I/O patterns, unbounded work, caching |
| **Security** | | Secrets, path traversal, injection, unsafe deserialization, permissions |
| **Operability** | | CI, release, config validation, graceful degradation, observability |
| **Plan hygiene** | | STATUS vs milestones, false “done”, missing deps, scope creep |

---

## 4. What to look for (gap catalog)

Use this as a **hunting list**. Not every item applies to every project; mark N/A explicitly.

### 4.1 Plan ↔ code drift (high signal for agent-maintained repos)

| Symptom | How to detect | Typical fix |
|---------|---------------|-------------|
| Milestone marked done but criteria unchecked | Read milestone file completion section | Reopen as partial; spawn follow-up milestone |
| Architecture doc lists files that don’t exist | Glob vs module map | Update doc or implement extraction milestone |
| Test count in STATUS stale | Re-count test functions / PASS lines | Update STATUS; add CI test summary |
| “Non-goals” contradict pending milestones | Compare STATUS goals vs milestone backlog | Reconcile charter text |
| Duplicate milestone scope | Same issue in M36, M37, M38 | Merge or split with explicit ownership |

### 4.2 Architecture smells

- **God file / god module** — single file > `{threshold}` lines or > `{N}` responsibilities
- **Dual source of truth** — e.g. model field + global for same state
- **Leaky abstraction** — UI layer knows storage format; domain imports UI framework
- **Incomplete consolidation** — “canonical” API exists but N call sites bypass it
- **Implicit state machine** — modes encoded in booleans without documented transitions
- **Package cycle** — especially after extractions

### 4.3 Correctness traps

- **Stack/history bugs** — navigation pushes state twice (manual + shared helper)
- **Copy vs pointer receivers** — mutation on value receiver silently discarded
- **Reload paths** — rescan/refresh re-enters user flows with wrong side effects
- **Nil optional dependencies** — crash when subsystem unavailable
- **Time-based polling** — full rebuild on coarse signal (e.g. directory mtime only)

### 4.4 Test gaps that matter

Prioritize tests that would have **caught shipped bugs**, not coverage for its own sake:

- User workflows end-to-end (input → state → render)
- Regression tests for each fixed production bug
- Error paths (missing files, corrupt config, partial failure)
- Mode-specific keybindings / behavior
- No tests for “helpers” that duplicate 20 lines of setup — **testutil debt**

### 4.5 Performance (proportionate)

- Profile or benchmark **before** large optimizations
- Measure with **realistic fixture sizes** (small test vaults mislead)
- Watch: O(n) work on every keystroke, full scans on timer, re-parse on resize
- Document **decision gates**: optimize only if benchmark exceeds `{threshold}`

### 4.6 Documentation agents rely on

Agents read docs instead of code when possible. Stale docs cause **reintroduced bugs**:

- Keybinding / config reference
- Module map / data flow
- “Done” milestones with wrong status
- AGENTS rules contradicting codebase (e.g. “use globals in theme.go” while refactor pending)

---

## 5. Finding format (required)

Every finding uses this structure:

```markdown
### F-{NN}: {Short title}

**Severity:** critical | high | medium | low  
**Category:** correctness | architecture | performance | tests | docs | plan | security  
**Evidence:** `{file:line}` or command output  
**Impact:** {user or maintainer consequence}  
**Recommendation:** {specific fix or milestone}  
**Challenge:** {why might this be wrong or lower priority?}  
**Decision:** {accept / defer / reject — for planning session}
```

Severity rubric:

| Level | Definition |
|-------|------------|
| **Critical** | Data loss, security breach, crash loop, corruption of user state |
| **High** | Broken feature, wrong results, major maintainability blocker |
| **Medium** | Inconsistency, scale risk, test/doc gap with moderate impact |
| **Low** | Polish, minor duplication, nice-to-have optimization |

---

## 6. Challenge your own review (mandatory section)

Before finalizing milestones, answer:

1. **False positives** — Which findings might be acceptable tradeoffs for this project stage?
2. **Ordering** — Are we fixing docs before code that changes behavior (or vice versa)?
3. **Milestone sizing** — Is any item > 3 days? Split it.
4. **Dependency creep** — Do proposed fixes violate project constraints?
5. **Over-engineering** — Are we adding CI, packages, or infra before proving need?
6. **Missing stakeholders** — Did we consider upgrade/migration path for users?
7. **Verification** — Does every milestone have **objective** completion criteria?

Record reversals and rationale in **Section 8: Decision log**.

---

## 7. Milestone creation rules

Create a new milestone when **any** of:

- Multi-file change or > 1 day estimated work
- Architectural boundary move (new package, state ownership change)
- Behavior change users will notice (needs tests + doc update)
- Follow-up to a “partial” milestone closure

Do **not** create a milestone for: typos, single-line fixes, comment-only changes.

### Milestone document minimum (each `M{N}-*.md`)

```markdown
# M{N} — {Title}

**Status:** ⏳ pending | 🚧 in progress | 🟡 partial | ✅ done | ⏸ deferred

## Goal
One sentence.

## Problem statement
What is broken or missing; link to finding F-{NN}.

## Out of scope
What this milestone explicitly does NOT do.

## Dependencies
Blocks / blocked by / parallel-safe.

## Work packages
### WP1 — {name} ({time})
- Step 1…
- Verification: …

## Files to modify
| File | Change |

## Acceptance criteria
- [ ] Objective, testable items only

## Rollback / risk
What could go wrong; how to revert.

## Handoff notes
What the executing agent must read first.
```

After creating milestones:

- [ ] Register in `{STATUS.md}` with phase, priority, dates
- [ ] Update execution order / dependency diagram
- [ ] Link from review doc to each milestone
- [ ] Mark superseded milestones as `partial → M{N}`

---

## 8. Decision log

| ID | Decision | Alternatives considered | Rationale | Date |
|----|----------|-------------------------|-----------|------|
| D-1 | | | | |

---

## 9. Execution handoff package

When planning is “ready for adoption,” the review must produce:

1. **Prioritized backlog** — ordered list with decision gates  
2. **Expanded milestones** — work packages with verification per step  
3. **Updated STATUS** — no false “done”  
4. **Doc sync milestone** — if docs are stale (often M-doc-sync)  
5. **CI milestone** — early if team velocity is high (protect refactors)  
6. **This review file** — linked from STATUS  

Executing agents start with:

1. One milestone only (unless explicitly parallelized)  
2. Read milestone **Handoff notes** + linked findings  
3. Run `{test}` after **each work package**, not only at end  
4. Update milestone checkboxes + STATUS on completion  
5. Do not expand scope — new gaps → new milestone  

---

## 10. Post-review verification (when execution completes)

- [ ] Re-run full review checklist (Section 2)  
- [ ] Compare acceptance criteria vs git diff  
- [ ] Update ARCHITECTURE doc module map  
- [ ] Close or mark partial milestones accurately  
- [ ] Add “Review closure” note with date and test count  

---

## Appendix A — Quick commands (customize per project)

```bash
# Example — replace for your stack
make test && make vet
go test ./... -cover
go test ./... -v -count=1 | grep -c '^--- PASS'
wc -l {source dirs}
rg "TODO|FIXME|HACK" --glob '*.go'
```

## Appendix B — Anti-patterns from past reviews (examples)

- Marking refactor “done” when only the happy path was consolidated  
- Integration tests that never simulate the bug scenario  
- Performance milestones without benchmarks on realistic data  
- Splitting `model.go` and `handlers.go` in the same PR (hard to review)  
- Doc sync before behavior freeze (requires second doc pass)

---

*Template version: 1.0 — copy as-is for new projects; do not treat project-specific examples in other review files as part of this template.*
