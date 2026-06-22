# Releasing obsidian-terminal

## Prerequisites
- All changes for the release are merged to `main`
- `make test && make vet` pass
- You've decided on a semver version number

## Automated release

Pushing a `v*` tag triggers `.github/workflows/release.yml`, which:

1. Creates a GitHub Release with auto-generated release notes
2. Downloads the source tarball and computes its SHA256
3. Opens a PR in `lthiagol/homebrew-tap` updating the formula's `url` + `sha256`
4. The tap's `brew test-bot` CI runs on the PR automatically

## Release steps

```bash
# 1. Tag
git tag -a vMAJOR.MINOR.PATCH -m "Release vMAJOR.MINOR.PATCH"

# 2. Push
git push origin vMAJOR.MINOR.PATCH

# 3. Watch the Release workflow in the Actions tab
#    https://github.com/lthiagol/obsidian-terminal/actions

# 4. Wait for the formula PR to appear in lthiagol/homebrew-tap
#    https://github.com/lthiagol/homebrew-tap/pulls

# 5. Wait for brew test-bot to pass on the PR, then merge

# 6. Verify
brew update && brew upgrade lthiagol/tap/obsidian-terminal
```

## Rollback

- If the workflow fails: delete the tag (`git push origin :vMAJOR.MINOR.PATCH`), fix the issue, re-tag
- If the formula PR fails `test-bot`: do not merge — investigate the formula, push fixes to the PR branch
- If users have already installed the bad version: delete the tag, release a patch version immediately
