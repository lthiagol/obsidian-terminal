# Code Review: M36-M41 Milestone Refinement

## Summary

The M36-M41 milestones identify real issues, but have several problems:
- **Overlapping scope**: M36, M37, and M38 all touch the same issues
- **Numbering conflicts**: C3 means different things in different milestones
- **Missing critical bugs**: mouse.go open-note paths missing side effects
- **Incorrect analysis**: M39's scroll estimation issue is misdiagnosed
- **Oversized milestones**: M36 and M38 are too large and should be split

## Issues Found

### 1. M36 Scope Problem

M36 claims to fix 6 "critical bugs" but mixes:
- **Quick fixes** (C3 accidental quit, C4 palette filter, C5 SetSize) — 1-2 hours each
- **Major refactoring** (C1 theme system) — 8-10 hours, essentially M37

**Recommendation**: Split M36 into:
- **M36a**: Quick bug fixes (C3, C4, C5, C6) — 1 day
- **M36b**: Theme system refactor (C1) — merge with M37

### 2. M36/M37/M38 Overlap

| Issue | M36 | M37 | M38 |
|-------|-----|-----|-----|
| C1: Theme globals | ✅ | ✅ | |
| C6: applyProfile broken | ✅ | | ✅ (as M10) |
| H2: Open note duplication | | | ✅ |

**Recommendation**: 
- M36 should only include quick fixes (C3, C4, C5)
- M37 should include C1 (theme) + C6 (applyProfile)
- M38 should focus on model.go split + H2 (open note duplication)

### 3. Missing Critical Bug: mouse.go Missing Side Effects

**Severity**: Critical  
**Location**: mouse.go lines 131-140, 143-158

The mouse handlers for opening notes from the tree and search results are **incomplete duplicates** that don't call:
- `m.buildOutline()`
- `m.backlinkPanel = NewBacklinkPanel(...)`
- `m.addRecentNote(...)`

**Impact**: Clicking a file in the tree doesn't update the outline, backlinks, or recents. This is a user-facing bug.

**Recommendation**: Add as **M36d** (part of quick fixes) or create **M36e** as separate milestone.

### 4. Missing Nil Vault Checks

M36's C2 only mentions handlers.go:248, but there are **3 locations** that call `ResolveWikiLink` without nil checks:
1. handlers.go:248 (embed resolver) — mentioned in M36
2. handlers.go:116 (view mode link follow) — **MISSING**
3. command_palette.go:75 (follow link command) — **MISSING**

**Recommendation**: Update M36 to include all 3 locations.

### 5. M39 Scroll Estimation Misdiagnosis

M39 claims `estimateYOffset` has issues with "ANSI escape sequences" but:
- `RenderSegmentsPlain` (markdown.go:1464) only concatenates `seg.Text`
- It does NOT produce ANSI codes
- The real issue is **multi-width characters** (CJK, emoji) making `len(text)` inaccurate

**Recommendation**: Correct M39's description. The fix is the same (use rune count instead of byte count), but the root cause is different.

### 6. Numbering Conflicts

The original code review used C1-C6 for critical bugs. The milestones renumbered them inconsistently:

| Original | M36 | M39 | Issue |
|----------|-----|-----|-------|
| C1: Theme race | C1 | | |
| C2: Nil vault | C2 | | |
| C3: Accidental quit | C3 | | |
| C4: Palette filter | C4 | | |
| C5: SetSize dead | C5 | | |
| C6: applyProfile | C6 | | |
| C7: ANSI bleed | | C3 | **Conflict**: M39 reuses C3 |
| C8: Scroll drift | | L8 | |

**Recommendation**: Use consistent numbering. Keep original C1-C8 for critical/high issues. Use M1-M10 for medium issues. Use L1-L10 for low issues.

### 7. M38 Oversized

M38 tries to:
- Split model.go (868 → 250 lines)
- Create `transitionToNote` to eliminate 6 duplicates
- Fix applyProfile (M10)

This is 3-4 days of work. Should be split into:
- **M38a**: Create `transitionToNote` + fix all 6 call sites (1 day)
- **M38b**: Extract subsystems from model.go (2-3 days)

### 8. Missing Milestones

The review missed several important areas:

#### 8a. Config Validation
**Severity**: Medium  
**Issue**: Config values are not validated. Invalid values are silently ignored.
- `line_spacing: invalid` → silently uses default
- `theme: nonexistent` → error, but no guidance
- `vault_path: /nonexistent` → error, but no suggestion

**Recommendation**: Create **M42: Config Validation**

#### 8b. Graceful Degradation
**Severity**: Medium  
**Issue**: If the vault becomes inaccessible during runtime (deleted, permissions changed), the app crashes or shows confusing errors.

**Recommendation**: Create **M43: Graceful Degradation**

#### 8c. Integration Testing
**Severity**: High  
**Issue**: No end-to-end tests for the full rendering pipeline:
- Parse markdown → Render to ANSI → Viewport softWrap → Display
- User types search → Filter results → Select → Open note → View

**Recommendation**: Create **M44: Integration Test Suite**

#### 8d. Performance Profiling
**Severity**: Low  
**Issue**: Known performance concerns not addressed:
- `DefaultKeys()` allocates on every tree update (L3)
- Model struct (30 fields) passed by value in Update/View
- No profiling data to guide optimization

**Recommendation**: Create **M45: Performance Profiling & Optimization**

### 9. M41 Scope Problem

M41 mixes:
- Dead code removal (IconVertical, IconDiamond, truncateContent)
- Godoc comments (15+ symbols)
- Hardcoded color fixes (9 places)
- Performance fix (KeyMap caching)
- UX fix (cursor reset)

**Recommendation**: Split into:
- **M41a**: Dead code removal (1 hour)
- **M41b**: Godoc comments (2 hours)
- **M41c**: Hardcoded color fix (1 hour)
- **M41d**: Performance & UX fixes (1 hour)

## Refined Milestone Plan

### Phase 7: Critical Bug Fixes (Priority: 🔴 Immediate)

| # | Milestone | Scope | Est. Time |
|---|-----------|-------|-----------|
| **M36** | Quick Bug Fixes | C3 (quit), C4 (palette), C5 (SetSize), mouse.go side effects, nil vault checks (3 sites) | 1 day |
| **M37** | Theme System Refactor | C1 (globals → model fields), C6 (applyProfile pointer receiver) | 2 days |

### Phase 8: Architecture Improvements (Priority: 🟡 High)

| # | Milestone | Scope | Est. Time |
|---|-----------|-------|-----------|
| **M38** | Eliminate Open-Note Duplication | Create `transitionToNote`, fix all 6 call sites | 1 day |
| **M39** | Split model.go | Extract pin_handler, outline_render, daily_handler, recent_handler, profile_handler | 2-3 days |
| **M40** | ANSI Wrapping Fixes | Style preservation across wrap boundaries, scroll estimation fix (multi-width chars) | 1 day |

### Phase 9: Code Quality (Priority: 🟢 Medium)

| # | Milestone | Scope | Est. Time |
|---|-----------|-------|-----------|
| **M41** | Config & Parser Hardening | YAML indent handling, remove duplicate heading parsers, extract constants | 1 day |
| **M42** | Dead Code & Cleanup | Remove unused exports, fix hardcoded colors, cache KeyMap, fix cursor reset | 1 day |
| **M43** | Godoc & Documentation | Add comments to 15+ exported symbols | 2 hours |

### Phase 10: Robustness (Priority: 🔵 Future)

| # | Milestone | Scope | Est. Time |
|---|-----------|-------|-----------|
| **M44** | Config Validation | Validate all config values, provide helpful error messages | 1 day |
| **M45** | Graceful Degradation | Handle vault inaccessibility, partial scan failures | 1 day |
| **M46** | Integration Test Suite | End-to-end tests for rendering pipeline, search flow, mode transitions | 2 days |
| **M47** | Performance Profiling | Profile hot paths, fix KeyMap allocation, optimize Model passing | 1 day |

## Execution Order

1. **M36** (Quick Bug Fixes) — do first, unblocks everything
2. **M37** (Theme Refactor) — required before M39 (ANSI fixes depend on theme being stable)
3. **M38** (Open-Note Duplication) — reduces duplication, makes M39 easier
4. **M39** (Split model.go) — large refactor, do after M38
5. **M40** (ANSI Wrapping) — independent, can be done anytime after M37
6. **M41-M43** (Code Quality) — can be done in parallel or sequentially
7. **M44-M47** (Robustness) — future work, lower priority

## Dependencies

```
M36 (Quick Fixes)
  ↓
M37 (Theme Refactor)
  ↓
M38 (Open-Note Duplication)
  ↓
M39 (Split model.go)
  ↓
M40 (ANSI Wrapping) — can start after M37
  ↓
M41-M43 (Code Quality) — independent, can be parallel
  ↓
M44-M47 (Robustness) — future
```

## Corrections to Existing Milestones

### M36 Corrections
- Remove C1 (theme refactor) — move to M37
- Remove C6 (applyProfile) — move to M37
- Add mouse.go missing side effects bug
- Add 2 missing nil vault check locations

### M37 Corrections
- Add C6 (applyProfile) from M36
- Clarify that this is a major refactor (2 days, not 1)

### M38 Corrections
- Remove applyProfile fix (already in M37)
- Split into M38a (transitionToNote) and M38b (extract files)

### M39 Corrections
- Fix C3 numbering conflict — use original C7
- Correct scroll estimation diagnosis (multi-width chars, not ANSI)

### M40 Corrections
- No major changes needed

### M41 Corrections
- Split into M41a-M41d (dead code, godoc, colors, perf/UX)

## New Milestones to Create

- **M42**: Config Validation
- **M43**: Graceful Degradation
- **M44**: Integration Test Suite
- **M45**: Performance Profiling & Optimization

## Conclusion

The original M36-M41 plan identified real issues but had scope problems, overlaps, and missed some critical bugs. The refined plan:
- Splits oversized milestones
- Removes overlaps
- Adds missing critical bugs
- Corrects misdiagnoses
- Adds missing robustness milestones
- Provides clear execution order with dependencies

Total estimated work: **18-22 days** (down from original estimate of 25+ days due to better scoping)
