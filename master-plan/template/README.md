# Master Plan Templates

Use these templates when creating or updating project planning artifacts. They encode the practices from Phase 12 planning (work packages, decision gates, verification per step).

## Files

| Template | Use when |
|----------|----------|
| [MILESTONE-TEMPLATE.md](./MILESTONE-TEMPLATE.md) | Creating any new `milestones/M{N}-*.md` |
| [STATUS-TEMPLATE.md](./STATUS-TEMPLATE.md) | Bootstrapping or major refresh of `STATUS.md` |
| [REVIEW-TEMPLATE.md](../REVIEW-TEMPLATE.md) | Architecture/code reviews (project-agnostic, lives at `master-plan/` root) |

## Workflow

1. **Review** — Copy `REVIEW-TEMPLATE.md` → `ARCHITECTURE-REVIEW-{date}.md`; record findings.
2. **Milestone** — Copy `MILESTONE-TEMPLATE.md` → `milestones/M{N}-slug.md`; fill every section; delete unused WPs.
3. **Register** — Add row to `STATUS.md` (use status legend from `STATUS-TEMPLATE.md`).
4. **Execute** — One work package (WP) per session/PR; run verification after each WP.
5. **Close** — Check all acceptance criteria; update STATUS dates and test count; sync milestone status emoji.

## Status legend (single source of truth)

| Emoji | Meaning |
|-------|---------|
| ⏳ pending | Planned, not started |
| 🚧 in progress | Work underway |
| 🟡 partial | Started or “done” but follow-up milestone remains |
| ✅ done | All acceptance criteria met |
| ⏸ deferred | Intentionally postponed; reactivation criteria documented |
| ⏳ deferred | Same as ⏸; prefer ⏸ for clarity |

## Rules

- Milestone file status must match `STATUS.md` (or say `partial → M{N+1}`).
- Every milestone needs **out of scope** and **dependencies**.
- Complex work uses **WPs** with verification checkboxes — not vague “Steps”.
- Test count in STATUS: `go test ./... -v -count=1 | grep -c '^--- PASS'`
- Do not mark ✅ without checked acceptance criteria in the milestone file.

## Examples in this repo

- Expanded milestone: [M50-navigation-history-fix.md](../milestones/M50-navigation-history-fix.md)
- Feature milestone: [M48-preview-pane.md](../milestones/M48-preview-pane.md), [M49-graph-view.md](../milestones/M49-graph-view.md)
- Execution batching: [PHASE-12-EXECUTION-PLAN.md](../PHASE-12-EXECUTION-PLAN.md)
