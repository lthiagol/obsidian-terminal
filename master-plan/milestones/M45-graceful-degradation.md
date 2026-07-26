# M45 — Graceful Degradation

**Status:** ✅ done

## Goal

Handle vault inaccessibility and partial scan failures gracefully without crashing or showing confusing errors.

## Issues

### Vault becomes inaccessible during runtime

If the vault directory is deleted, permissions change, or the filesystem becomes unavailable:
- `checkVaultChanges` may panic on `os.Stat` errors
- `rescanVault` may return nil vault, causing nil pointer dereferences
- Opening notes fails with cryptic errors
- The app continues running in a broken state

### Partial scan failures

If some files/directories can't be read during scanning:
- `ScanVault` returns `scanErrors` but the app doesn't display them clearly
- Users don't know which files failed to load
- No way to retry failed scans

### Missing error recovery

- No "Retry" option when vault scanning fails
- No way to switch to a different vault if the current one is broken
- No indication of which features are disabled due to scan failures

## Design

### Vault state tracking

Add a `vaultState` field to Model:
```go
type VaultState int

const (
    VaultStateOK VaultState = iota
    VaultStatePartial  // some files/dirs failed to scan
    VaultStateBroken   // vault is inaccessible
)
```

### Graceful degradation

1. **VaultStateOK**: Normal operation
2. **VaultStatePartial**: 
   - Show toast with scan errors count
   - Display failed paths in a "Scan Errors" panel (accessible via command palette)
   - Allow retry via Ctrl+R
3. **VaultStateBroken**:
   - Show error screen with clear message
   - Offer options: "Retry", "Switch Vault", "Quit"
   - Disable features that require vault access (search, tags, backlinks)
   - Keep tree visible (showing last known state)

### Error recovery

1. **Retry button**: Re-run `rescanVault` and update state
2. **Switch vault**: Open profile picker to select a different vault
3. **Auto-recovery**: If vault becomes accessible again (e.g., permissions restored), automatically recover

### Scan error display

Add a "Scan Errors" command to the command palette that shows:
- List of failed paths
- Error message for each
- Option to retry

## Files to modify

| File | Changes |
|------|---------|
| `model.go` | Add `vaultState` field; update `checkVaultChanges` and `rescanVault` to handle errors gracefully |
| `vault.go` | Improve error reporting in `ScanVault`; track which paths failed |
| `handlers.go` | Add error recovery handlers (retry, switch vault) |
| `command_palette.go` | Add "Scan Errors" and "Retry Scan" commands |
| `statusbar.go` | Show vault state indicator (OK/Partial/Broken) |
| New: `error_screen.go` | Error screen UI for broken vault state |

## Completion Criteria

- [x] Vault state is tracked (OK/Partial/Broken)
- [x] Partial scan failures show clear error messages with failed paths
- [x] Broken vault shows error screen with recovery options
- [x] Retry option re-scans the vault
- [x] Switch vault option opens profile picker
- [x] Features that require vault are disabled when vault is broken
- [x] No panics or crashes when vault becomes inaccessible
- [x] `make test` passes all tests (add error handling tests)
- [x] `make vet` exits 0

## Completed

2026-06-12

Added `VaultState` type with three states: `VaultStateOK`, `VaultStatePartial`, `VaultStateBroken`.

**checkVaultChanges hardened:**
- Detects vault inaccessibility → sets `VaultStateBroken` with error toast
- Auto-recovers when vault becomes accessible again → rescans
- No longer spams toasts (only fires on state transition)

**rescanVault hardened:**
- Stat/sync failures now set `VaultStateBroken` with toast (was silent return)
- Scan errors update `vaultState` to `Partial`/`OK` based on count
- `vaultStateFrom(scanErrorCount)` helper added

**Error screen:**
- When vault is broken, `renderBrokenVaultScreen()` shows in right panel
- Displays error message, recovery instructions (r/P/q)
- 'r' key triggers rescan from broken state
- Quit from scan errors display blocked (added `!m.scanErrorsVisible` to quit guard)

**Scan errors display:**
- `showScanErrors()` toggles `scanErrorsVisible`
- `renderScanErrors()` renders list of failed paths with error details
- "Scan Errors" command added to command palette (only when errors exist)
- Esc/'q' dismiss the display

**Status bar:**
- Shows " BROKEN" indicator when vault inaccessible
- Shows "(N scan errors)" only when `vaultState == Partial`

9 new tests.

## Estimated Time

1-2 days
