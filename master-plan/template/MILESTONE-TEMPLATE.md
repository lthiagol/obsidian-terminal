# M{NN} — {Title}

**Status:** ⏳ pending  
**Phase:** {N} — {phase name}  
**Priority:** 🔴 Immediate | 🟡 High | 🟢 Medium | 🔵 Future  
**Finding:** {F-NN in ARCHITECTURE-REVIEW-*.md, or "—"}  
**Execution plan:** {link to phase plan, or "—"}

## Goal

{One sentence: what user-visible or architectural outcome this achieves.}

## Problem statement

{What is broken, missing, or risky today? Cite evidence: file:line, test gap, or review finding.}

## Out of scope

Explicit boundaries — prevents scope creep during execution:

- {Item 1 — e.g. "No new dependencies"}
- {Item 2 — e.g. "No keybinding changes"}
- {Item 3}

## Dependencies

| Relation | Milestone / artifact |
|----------|----------------------|
| **Blocked by** | {Mxx WPn, or "nothing"} |
| **Blocks** | {Mxx, doc updates, or "nothing"} |
| **Parallel-safe with** | {Mxx, or "nothing — do alone"} |

## Design (approved for execution)

{API sketches, UX decisions, data structures, keybindings. Resolve forks here — not during coding.}

### Key decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| {e.g. placement} | {Option A} | {why} |

---

## Work packages

> **Rule:** One WP ≈ one focused session (2–6h). Run `{test command}` after each WP. One WP per commit/PR recommended.

### WP1 — {Title} ({time estimate})

**Steps:**
1. {Concrete step}
2. {Concrete step}

**Verification:**
- [ ] {Objective check — test name, command, or measurable outcome}
- [ ] `{make test && make vet}` pass (if code changed)

---

### WP2 — {Title} ({time estimate})

**Steps:**
1. …

**Verification:**
- [ ] …

---

{Add WP3…WPn as needed. Delete unused WPs before execution.}

---

## Files to modify

| File | Changes |
|------|---------|
| `{path}` | {summary} |
| `{path}` | **New** — {summary} |

## Test plan

| ID | Scenario | Type | WP |
|----|----------|------|-----|
| T1 | {description} | unit / integration | WP{n} |

## Acceptance criteria (milestone done)

All must be checked before setting status to ✅:

- [ ] All WPs verified
- [ ] {Behavior criterion}
- [ ] {Doc criterion — KEYBINDINGS, help.go, etc.}
- [ ] `{make test && make vet}` pass
- [ ] `STATUS.md` updated (status, dates, test count delta)
- [ ] Milestone file status emoji matches STATUS

## Rollback / risk

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| {risk} | low/med/high | {mitigation} |

## Handoff notes

{What the executing agent must read first; what not to touch; link to related WPs in other milestones.}

## Estimated total

{hours or days — sum of WPs with buffer}

## Completion log

_Fill when done:_

| Field | Value |
|-------|-------|
| Started | {YYYY-MM-DD} |
| Completed | {YYYY-MM-DD} |
| Tests added | {N} |
| Notes | {deviations from plan} |
