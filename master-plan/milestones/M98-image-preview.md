# M98 — Image Preview

**Status:** ⏳ pending  
**Priority:** Low (85-99 range — address individually)

## Goal

Display inline images in the terminal using sixel or kitty graphics protocol.

## Notes

- Only works on terminals with sixel or kitty protocol support (WezTerm, kitty, iTerm2)
- Image files in the vault can be previewed inline
- Requires detecting terminal capabilities at startup
- May need to resize images to fit the viewport
