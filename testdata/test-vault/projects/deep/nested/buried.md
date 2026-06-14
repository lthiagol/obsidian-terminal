---
title: Deeply Buried Note
tags: [test, nested]
---

# Deeply Buried Note

This note is nested three levels deep to test the file tree navigation.

## Purpose

Verify that the tree component can handle deeply nested directory structures and that users can navigate to and open files at any depth.

## Expected Behavior

1. Root folder (test-vault) is expanded
2. `projects/` is visible, collapsed by default
3. Expand `projects/` → see `deep/`
4. Expand `deep/` → see `nested/`
5. Expand `nested/` → see `buried.md`
6. Select `buried.md` and press Enter to open
