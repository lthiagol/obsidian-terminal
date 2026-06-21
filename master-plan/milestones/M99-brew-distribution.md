# M99 — Release Automation (Homebrew Formula PR)

**Status:** ⏳ pending (placeholder — detail when reactivated)  
**Phase:** 99 — Future (Low Priority)  
**Priority:** 🔵 Low

## Goal

Automate the release pipeline so that pushing a `v*` tag to `obsidian-terminal` creates a GitHub Release and opens a PR to `lthiagol/homebrew-tap` updating the formula `url` + `sha256` — without manual intervention.

## Problem statement

**What's already done (verified 2026-06-21):**

| Component | State | Location |
|-----------|-------|----------|
| Homebrew tap repo | ✅ exists | `github.com/lthiagol/homebrew-tap` (multi-project structure) |
| `obsidian-terminal` formula | ✅ exists, pinned to `v0.1.0` | `Formula/obsidian-terminal.rb` |
| Tap CI (`brew test-bot`) | ✅ exists, matrix on macOS + Linux | `.github/workflows/tests.yml` in tap repo |
| Project README in tap | ✅ exists | `projects/obsidian-terminal/README.md` |
| obsidian-terminal README | ✅ has `brew install lthiagol/tap/obsidian-terminal` | `README.md` in main repo |
| `v0.1.0` tag | ✅ exists | git tag in main repo |

**What's missing:**

1. **No release workflow in `obsidian-terminal`** — pushing a `v*` tag does nothing except create a tag. No GitHub Release, no formula bump, no release notes.
2. **Formula is pinned to `v0.1.0`** — every new release requires a human to manually:
   - Edit `Formula/obsidian-terminal.rb` in the tap repo
   - Update the `url` to the new tag's tarball
   - Compute the SHA256 of the new tarball
   - Update the `sha256` field
   - Push to the tap repo
   - Wait for `brew test-bot` to verify
3. **No release notes / CHANGELOG** — the v0.1.0 GitHub release (if it exists) has no auto-generated notes.
4. **No version single-source-of-truth** — the version is implicit in the git tag; the formula hardcodes it in the URL.

This milestone closes those gaps with a GitHub Actions release workflow.

## Out of scope

- Pre-built binaries — the formula builds from source via `go build`; no bottle/build matrix needed in the obsidian-terminal repo (the tap's `test-bot` handles bottle building on PRs)
- Cross-compilation — Go builds natively on the user's machine via Homebrew's `go` dependency
- Other package managers (APT, RPM, AUR, etc.)
- Auto-detection of breaking changes for semver bumping — human decides the version number
- goreleaser — the project has no existing `.goreleaser.yml`; adding it would be a separate decision. This milestone uses plain GitHub Actions.

## Dependencies

| Relation | Milestone / artifact |
|----------|----------------------|
| **Blocked by** | nothing (tap repo + formula already exist) |
| **Blocks** | nothing |
| **Parallel-safe with** | everything — pure CI/infra work, no code changes |
| **External repos involved** | `github.com/lthiagol/homebrew-tap` (formula PR target) |

## Design (approved for execution — detailed 2026-06-21)

### Release workflow architecture

```
Push `v*` tag to obsidian-terminal
   ↓
GitHub Actions: release.yml
   ├── 1. Create GitHub Release (auto-generated notes from commits since last tag)
   ├── 2. Download source tarball: refs/tags/{tag}.tar.gz
   ├── 3. Compute SHA256 of tarball
   ├── 4. Via gh api (no clone):
   │     ├── GET current Formula/obsidian-terminal.rb contents + file SHA
   │     ├── GET main branch SHA
   │     ├── POST new branch `bump-obsidian-terminal-{tag}` from main
   │     ├── Decode formula, sed-replace url + sha256, re-encode
   │     ├── PUT updated file to the new branch
   │     └── gh pr create (triggers tap's test-bot CI on the PR)
   └── 5. PR is ready for review + merge in lthiagol/homebrew-tap
```

**No git clone of the tap repo anywhere.** Every operation on the tap repo goes through the GitHub REST API via `gh api`. The only checkout is `actions/checkout@v4` on the obsidian-terminal repo itself (for `fetch-depth: 0` release notes).

### Secret required

| Secret name | Scope | Stored in | Purpose |
|-------------|-------|-----------|---------|
| `HOMEBREW_TAP_TOKEN` | fine-grained PAT on `lthiagol/homebrew-tap` only | `obsidian-terminal` repo settings → Secrets and variables → Actions | Allows the workflow to create branches, update files, and open PRs to the tap repo via `gh api` |

**Why not `GITHUB_TOKEN`?** The default `GITHUB_TOKEN` is scoped to the obsidian-terminal repo only — it can't touch `lthiagol/homebrew-tap`. The PAT bridges the two repos.

**PAT creation steps (for handoff notes):**
1. Go to GitHub → Settings → Developer settings → Personal access tokens → Fine-grained tokens
2. Create a token scoped to `lthiagol/homebrew-tap` only (not all repos)
3. Permissions: Contents (read+write) — needed for GET file, POST branch, PUT file. Pull requests (read+write) — needed for `gh pr create`.
4. Copy the token, add it as a repository secret in `lthiagol/obsidian-terminal` → Settings → Secrets and variables → Actions → New repository secret → Name: `HOMEBREW_TAP_TOKEN`

### Release workflow file

**Path:** `.github/workflows/release.yml` in `obsidian-terminal` repo

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

permissions:
  contents: write   # for creating the GitHub Release in the obsidian-terminal repo

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0  # need full history for auto-generated release notes

      - name: Set tag variable
        id: tag
        run: echo "version=${GITHUB_REF#refs/tags/}" >> $GITHUB_OUTPUT

      - name: Create GitHub Release
        env:
          GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        run: |
          gh release create "${{ steps.tag.outputs.version }}" \
            --title "${{ steps.tag.outputs.version }}" \
            --generate-notes \
            --verify-tag

      - name: Compute SHA256 of source tarball
        id: sha
        run: |
          TAG="${{ steps.tag.outputs.version }}"
          URL="https://github.com/lthiagol/obsidian-terminal/archive/refs/tags/${TAG}.tar.gz"
          curl -fsSL "$URL" -o source.tar.gz
          SHA=$(sha256sum source.tar.gz | awk '{print $1}')
          echo "sha256=$SHA" >> $GITHUB_OUTPUT
          echo "url=$URL" >> $GITHUB_OUTPUT

      - name: Bump formula in homebrew-tap via API
        env:
          GH_TOKEN: ${{ secrets.HOMEBREW_TAP_TOKEN }}
          TAG: ${{ steps.tag.outputs.version }}
          NEW_URL: ${{ steps.sha.outputs.url }}
          NEW_SHA: ${{ steps.sha.outputs.sha256 }}
        run: |
          TAP="lthiagol/homebrew-tap"
          FORMULA="Formula/obsidian-terminal.rb"
          BRANCH="bump-obsidian-terminal-${TAG}"
          
          # 1. Get current file contents + file SHA (needed for PUT)
          FILE_JSON=$(gh api "repos/${TAP}/contents/${FORMULA}")
          FILE_SHA=$(echo "$FILE_JSON" | jq -r '.sha')
          FILE_B64=$(echo "$FILE_JSON" | jq -r '.content')
          
          # 2. Get main branch SHA (to branch from)
          MAIN_SHA=$(gh api "repos/${TAP}/git/refs/heads/main" | jq -r '.object.sha')
          
          # 3. Create new branch from main
          gh api "repos/${TAP}/git/refs" \
            -X POST \
            -f ref="refs/heads/${BRANCH}" \
            -f sha="${MAIN_SHA}"
          
          # 4. Decode formula, replace url + sha256, re-encode
          echo "$FILE_B64" | base64 -d > formula.rb
          sed -i -E "s|url \".*\"|url \"${NEW_URL}\"|" formula.rb
          sed -i -E "s|sha256 \".*\"|sha256 \"${NEW_SHA}\"|" formula.rb
          NEW_B64=$(base64 -w 0 formula.rb)
          
          # Show the diff for log traceability
          echo "--- Updated formula ---"
          cat formula.rb
          
          # 5. Update file on the new branch via API
          gh api "repos/${TAP}/contents/${FORMULA}" \
            -X PUT \
            -f message="Bump obsidian-terminal to ${TAG}" \
            -f content="${NEW_B64}" \
            -f branch="${BRANCH}" \
            -f sha="${FILE_SHA}"
          
          # 6. Open PR
          gh pr create \
            --repo "${TAP}" \
            --base main \
            --head "${BRANCH}" \
            --title "Bump obsidian-terminal to ${TAG}" \
            --body "Automated bump triggered by tag push in [obsidian-terminal](${GITHUB_SERVER_URL}/${GITHUB_REPOSITORY}/releases/tag/${TAG})."
```

**Notes on the workflow:**
- `jq` and `base64` are pre-installed on `ubuntu-latest` runners — no setup needed
- `sed -i -E` is GNU sed (default on ubuntu-latest)
- No `git clone`, no `git push`, no git credentials configured — all tap-repo mutations go through `gh api`
- The PR body links back to the release for traceability
- `--label "automated"` omitted by default — add it if an "automated" label exists in the tap repo
- The PUT contents API requires the current file SHA (`FILE_SHA`) to prevent clobbering concurrent edits — we fetched it in step 1
- The workflow prints the updated formula to the log so you can verify the sed replacements without inspecting the PR

### Versioning strategy

- **Tags:** `vMAJOR.MINOR.PATCH` (semver) — e.g., `v0.2.0`, `v1.0.0`
- **Pre-releases:** `vMAJOR.MINOR.PATCH-rc.N` or `vMAJOR.MINOR.PATCH-beta.N` — the workflow triggers on `v*` so pre-release tags also trigger; add `--prerelease` flag conditionally if needed
- **No VERSION file** — the git tag is the source of truth; the formula extracts it from the URL
- **Release notes:** auto-generated from commits between the previous tag and the new tag (GitHub's `--generate-notes` feature)

### Key decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Use GitHub Actions, not goreleaser | Yes | No existing goreleaser config; plain Actions is simpler for this single-binary Go project |
| Build from source in formula (no pre-built binaries) | Yes (already done) | Go cross-compiles trivially; `brew install` builds from source using the `go` dep — no need to maintain a release artifact matrix |
| Open a PR to the tap, not push directly to main | Yes | PR triggers the tap's `brew test-bot` CI to verify the formula before merge — safer than direct push |
| Use a fine-grained PAT scoped to only the tap repo | Yes | Minimum-privilege; the token can't touch the obsidian-terminal repo or any other |
| **Pure `gh api` for tap mutations, no git clone** | Yes | Smaller blast radius (no git push credentials — just API calls); no clone step (faster, less runner disk); easier to audit (every API call is logged) |
| Tag is the source of truth for version | Yes | No VERSION file to forget to update; git tag is the release |
| `--generate-notes` for release body | Yes | Free changelog from commit history; no manual release notes needed |
| Trigger on `v*` (all version tags) | Yes | Catches pre-releases too; add conditional `--prerelease` if needed later |

---

## Work packages

### WP1 — Add release workflow (2h)

**Steps:**
1. Create `.github/workflows/release.yml` in the obsidian-terminal repo (see Design section for exact content)
2. Verify the workflow syntax: `actionlint .github/workflows/release.yml` (if actionlint installed) or just visual review
3. **Do not push a tag yet** — the secret `HOMEBREW_TAP_TOKEN` must be set first (WP2)
4. Commit the workflow on a branch, open a PR in obsidian-terminal, merge after review

**Verification:**
- [ ] `.github/workflows/release.yml` exists and is valid YAML
- [ ] Workflow triggers only on `v*` tag push (not on every push to main)
- [ ] No secrets referenced that don't exist yet (only `HOMEBREW_TAP_TOKEN`, which is set in WP2)
- [ ] `make test && make vet` pass (no code changes — workflow file only)

---

### WP2 — Create PAT and add as repo secret (30m)

**Steps:**
1. Go to GitHub → Settings → Developer settings → Personal access tokens → Fine-grained tokens
2. Click "Generate new token"
3. Settings:
   - **Resource owner:** `lthiagol`
   - **Repository access:** Only select repositories → `lthiagol/homebrew-tap`
   - **Permissions:** Contents (Read and write), Pull requests (Read and write)
   - **Expiration:** 1 year (renew before expiry)
4. Copy the generated token
5. Go to `lthiagol/obsidian-terminal` → Settings → Secrets and variables → Actions → New repository secret
   - **Name:** `HOMEBREW_TAP_TOKEN`
   - **Value:** (paste the token)
6. Verify the secret is visible in the repo settings (names only — values are hidden)

**Verification:**
- [ ] Fine-grained PAT created with scope limited to `lthiagol/homebrew-tap`
- [ ] PAT has Contents (write) + Pull requests (write) permissions
- [ ] `HOMEBREW_TAP_TOKEN` secret added to `lthiagol/obsidian-terminal` repo
- [ ] Secret is not committed to any file (only in GitHub UI)

---

### WP3 — First end-to-end release test (1h)

**Steps:**
1. Ensure all code to be released is merged to `main`
2. Choose a version number (e.g., `v0.2.0` if new features since `v0.1.0`, or `v0.1.1` for fixes)
3. Tag locally: `git tag -a v0.2.0 -m "Release v0.2.0"`
4. Push the tag: `git push origin v0.2.0`
5. Watch the Actions tab in `lthiagol/obsidian-terminal` — the "Release" workflow should start
6. Verify the workflow:
   - Step "Create GitHub Release" succeeds → check the Releases page for the new release with auto-generated notes
   - Step "Compute SHA256" succeeds → check the workflow log for the computed hash
   - Step "Bump formula in homebrew-tap" succeeds → check `lthiagol/homebrew-tap` for the new PR
7. Verify the PR in the tap repo:
   - Title: `Bump obsidian-terminal to v0.2.0`
   - Body: links back to the release
   - Diff: only `Formula/obsidian-terminal.rb` changed (url + sha256)
   - Tap CI (`brew test-bot`) starts on the PR
8. Wait for `brew test-bot` to pass, then merge the PR
9. Test locally: `brew install lthiagol/tap/obsidian-terminal` (after `brew update`)

**Verification:**
- [ ] GitHub Release created at `https://github.com/lthiagol/obsidian-terminal/releases/tag/v0.2.0`
- [ ] Auto-generated release notes present
- [ ] PR opened in `lthiagol/homebrew-tap` with correct url + sha256
- [ ] `brew test-bot` passes on the PR
- [ ] After merge: `brew update && brew install lthiagol/tap/obsidian-terminal` installs v0.2.0
- [ ] `obsidian-terminal --version` (or `--help`) confirms the new version

---

### WP4 — Document release process (30m)

**Steps:**
1. Create `RELEASE.md` in the obsidian-terminal repo:
   ```markdown
   # Releasing obsidian-terminal
   
   ## Prerequisites
   - All changes for the release are merged to `main`
   - `make test && make vet` pass
   - You've decided on a semver version number
   
   ## Release steps
   1. Tag: `git tag -a vMAJOR.MINOR.PATCH -m "Release vMAJOR.MINOR.PATCH"`
   2. Push: `git push origin vMAJOR.MINOR.PATCH`
   3. Watch the Release workflow in the Actions tab
   4. Wait for the formula PR to appear in `lthiagol/homebrew-tap`
   5. Wait for `brew test-bot` to pass on the PR
   6. Merge the PR
   7. Verify: `brew update && brew upgrade lthiagol/tap/obsidian-terminal`
   
   ## Rollback
   - If the workflow fails: delete the tag (`git push origin :vMAJOR.MINOR.PATCH`), fix the issue, re-tag
   - If the formula PR fails `test-bot`: do not merge — investigate the formula, push fixes to the PR branch
   - If users have already installed the bad version: yank the tag, release a patch version immediately
   ```
2. Optionally add a "Release" section to README pointing to `RELEASE.md`
3. Update STATUS.md: M99 → ✅

**Verification:**
- [ ] `RELEASE.md` exists and describes the one-command release process
- [ ] README links to `RELEASE.md` (optional)
- [ ] `STATUS.md` M99 → ✅ with completion log

---

## Files to modify

| File | Changes |
|------|---------|
| `.github/workflows/release.yml` | **New** — release + formula PR workflow |
| `RELEASE.md` | **New** — release process documentation |
| `README.md` | Optional: link to `RELEASE.md` in a "Releasing" section |
| `STATUS.md` | M99 → ✅ with completion log |

**No changes to the tap repo** — the workflow opens PRs to it, but the tap repo itself is already set up.

## Test plan

| ID | Scenario | Type | WP |
|----|----------|------|-----|
| T1 | Workflow YAML is valid | lint | WP1 |
| T2 | PAT created with correct scope | manual | WP2 |
| T3 | Tag push triggers Release workflow | e2e | WP3 |
| T4 | GitHub Release created with notes | e2e | WP3 |
| T5 | Formula PR opened in tap repo | e2e | WP3 |
| T6 | `brew test-bot` passes on PR | e2e | WP3 |
| T7 | `brew install` works after merge | e2e | WP3 |

## Acceptance criteria (milestone done)

- [ ] WP1–WP4 complete
- [ ] `.github/workflows/release.yml` merged to main
- [ ] `HOMEBREW_TAP_TOKEN` secret set in obsidian-terminal repo
- [ ] One successful end-to-end release (tag → release → formula PR → test-bot pass → merge → `brew install` works)
- [ ] `RELEASE.md` documents the process
- [ ] `STATUS.md` updated: M99 → ✅

## Rollback / risk

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| PAT leaks | low | Fine-grained scope (tap repo only); stored as GitHub secret, never in code; rotate if compromised |
| Workflow fails on tag push | medium | WP3 is the test; if it fails, delete the tag, fix, re-tag |
| `brew test-bot` fails on the formula PR | medium | Don't merge — investigate; common causes: sha256 mismatch, URL 404, build failure with new Go version |
| Source tarball URL changes format | low | GitHub's `archive/refs/tags/{tag}.tar.gz` format is stable; if GitHub changes it, update the workflow |
| Auto-generated release notes are noisy | low | Can edit the release manually after creation; or switch to manual notes in the workflow |
| Pre-release tag triggers full release | low | Add `--prerelease` conditional: `if [[ "$TAG" == *"-rc"* ]] || [[ "$TAG" == *"-beta"* ]]; then PRERELEASE="--prerelease"; else PRERELEASE=""; fi` and pass to `gh release create` |

**Rollback:** Delete the tag (`git push origin :vX.Y.Z`) and delete the GitHub Release. The formula PR can be closed without merging if the release is aborted.

## Handoff notes

**Read first:**
- This milestone file (especially the workflow YAML in the Design section)
- The existing formula: `Formula/obsidian-terminal.rb` in the tap repo — to understand the `url`/`sha256` line format the workflow's `sed` targets
- The tap's CI: `.github/workflows/tests.yml` in the tap repo — to understand what `test-bot` will check

**Do not:**
- Add goreleaser — not in scope; the project has no existing goreleaser config
- Build pre-built binaries — the formula builds from source; no artifact matrix needed
- Push the formula directly to main on the tap repo — always open a PR so `test-bot` runs first
- Clone the tap repo in the workflow — use `gh api` for all tap-repo mutations (the workflow is designed API-only)
- Commit the PAT to any file — only in GitHub repo secrets

**When stuck:**
- If the workflow's `sed` doesn't match the formula lines: check the exact format in `Formula/obsidian-terminal.rb` — the lines look like `url "..."` and `sha256 "..."` (double-quoted). Adjust the regex if the formula uses single quotes. The updated formula is printed to the workflow log for verification.
- If `gh release create` fails with "release already exists": delete the existing release first (`gh release delete vX.Y.Z`) or use `--draft` for testing.
- If `gh api` calls fail with 403: verify the PAT has `Contents: write` on the tap repo (for GET file, POST branch, PUT file).
- If `gh pr create` fails with 403: verify the PAT has `Pull requests: write` on the tap repo.
- If `brew test-bot` fails with sha256 mismatch: the workflow computed the wrong hash — verify by manually downloading the tarball and running `sha256sum`. The computed hash is printed in the workflow log (`echo "sha256=$SHA"`).
- If the PUT contents API fails with 409 "does not match": the file SHA changed between GET and PUT (someone pushed to the tap between steps 1 and 5). Re-run the workflow.

## Estimated total

4 hours (2h WP1 + 30m WP2 + 1h WP3 + 30m WP4)

## Priority

🔵 Low — Phase 99; execute when prioritized. The tap and formula already work for v0.1.0; this milestone automates future releases.

## Completion log

_Fill when done:_

| Field | Value |
|-------|-------|
| Started | — |
| Completed | — |
| Tests added | 0 (infra only — e2e verification is manual) |
| First release via automation | {tag, e.g. v0.2.0} |
| Notes | {any deviations; tap PR URL for first release} |
