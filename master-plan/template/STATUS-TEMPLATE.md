# {Project Name} — Build Status

**Last updated:** {YYYY-MM-DD}  
**Language / runtime:** {e.g. Go 1.26+ — see `go.mod`}  
**Architecture review:** [{latest review}](./ARCHITECTURE-REVIEW-{date}.md)  
**Execution plan:** [{phase plan}](./PHASE-{N}-EXECUTION-PLAN.md) _(if applicable)_  
**Templates:** [template/README.md](./template/README.md)

## Charter

### Goals

- **v1:** {bullet list}
- **v2:** {bullet list}

### Non-goals

- {explicit exclusions}

### Planned extras

- {features planned but not in v1 — link to milestones}

## Key decisions

| Decision | Choice |
|----------|--------|
| {area} | {choice} |
| {area} | {choice} |

## Status legend

| Emoji | Meaning |
|-------|---------|
| ⏳ pending | Not started |
| 🚧 in progress | Active work |
| 🟡 partial | Incomplete; follow-up milestone linked |
| ✅ done | Acceptance criteria met |
| ⏸ deferred | Postponed with reactivation criteria |

---

## Progress

### Phase {N}: {Name}

| Milestone | Status | Tests | Started | Completed |
|-----------|--------|-------|---------|-----------|
| M{nn}: {Title} | ⏳ pending | 0 | — | — |

_{Repeat per phase. Keep phases ordered by execution, not numeric id.}_

**Total tests:** {N} — refresh with:

```bash
go test ./... -v -count=1 | grep -c '^--- PASS'
```

---

## Execution order

Group milestones into **batches** with rationale. Reference milestone files, not just numbers.

### Batch 1: {Name} ({priority})

1. **M{nn}** — {one-line scope}
2. **M{nn}** — {one-line scope}

**Rationale:** {why this order}

### Batch 2: {Name}

…

**Minimum shippable batch:** {which batch delivers user value alone}

---

## Milestone dependencies

```
M{a} → M{b}
M{c} + M{d}  (parallel after M{b})
```

## Partial milestones

| Milestone | Done | Remaining |
|-----------|------|-----------|
| M{xx} | {what shipped} | → **M{yy}** |

---

## Keybinding / config references

- [{KEYBINDINGS.md}](../KEYBINDINGS.md)
- [{config example}](../config.yaml.example)

---

## Maintenance checklist

Run after each milestone closure:

- [ ] Milestone file: all acceptance criteria checked
- [ ] Milestone status emoji = STATUS table
- [ ] Test count updated
- [ ] Dates filled (Started / Completed)
- [ ] DESIGN / KEYBINDINGS updated if behavior changed
- [ ] No ✅ with open follow-up unless marked 🟡 partial
