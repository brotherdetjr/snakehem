# Bug Reports - Snakehem Refactoring

This file tracks bugs found during testing of the refactored codebase.

---

## Bug #1: Lobby State Shows Action UI Elements

**Date**: 2026-01-10
**Status**: ✅ FIXED
**Severity**: High
**Reporter**: User

### Steps to Reproduce
1. Run the game
2. Observe the lobby screen

### Expected Behavior
- "PLAYERS PRESS ANY BUTTON TO JOIN" is visible on a green field
- Nothing else should be displayed

### Actual Behavior
- "PLAYERS PRESS ANY BUTTON TO JOIN" is visible on a green field
- **BUG**: Also showing:
  - "THREE" (countdown text)
  - "TARGET SCORE: 999"
  - Timer ticking at the bottom

### Root Cause
In `internal/game/game.go`, the Draw() method is calling both `DrawLobbyUI()` AND `DrawActionUI()` in the LobbyState case. Only `DrawLobbyUI()` should be called.

### Fix
✅ Removed the `DrawActionUI()` call from the LobbyState case in `internal/game/game.go:111`.

### Test Coverage
✅ Created `internal/game/game_test.go` with 3 tests:
- `TestDraw_LobbyState` - Verifies lobby state only calls DrawLobbyUI (Bug #1 fix)
- `TestDraw_ActionState` - Verifies action state only calls DrawActionUI
- `TestDraw_ScoreboardState` - Verifies scoreboard state only calls DrawScoreboardUI

All tests passing ✅

### Files Changed
- `internal/game/game.go` - Removed erroneous DrawActionUI call
- `internal/game/game.go` - Added GameRenderer interface for testability
- `internal/game/game_test.go` - Added test coverage (NEW)

---

## Bug #2: Snake Score Not Displayed in Lobby

**Date**: 2026-01-10
**Status**: ✅ PARTIALLY FIXED (1 of 2 issues)
**Severity**: High
**Reporter**: User

### Steps to Reproduce
1. Run the game
2. Press spacebar or one of the arrow buttons

### Expected Behavior
- A snake link (little square) appears in the upper part of the screen
- "000" label appears at the top of the screen showing current player's score
- Both the link and 000-label are initially red, but quickly fade to white color (when no other keys pressed)

### Actual Behavior
- A snake link (little square) appears in the upper part of the screen ✅
- **BUG #1**: No 000-label is present ❌ → ✅ FIXED
- **BUG #2**: The link has red color, but doesn't fade to white ❌ → 🔴 STILL INVESTIGATING

### Root Cause
1. ✅ **Scores not displayed**: `DrawLobbyUI()` was not calling `DrawScores()`. In the old code (game/draw.go:26), scores were drawn in lobby state.
2. 🔴 **Redness not fading**: TBD - The lobby state Update() method IS calling `ChangeRedness(-0.1, tpsMultiplier)`, so need to investigate why it's not working in practice.

### Fix
✅ **Part 1 - Scores**: Modified `DrawLobbyUI()` in `internal/rendering/composite_renderer.go` to call `DrawScores()` with `isActionState=false` (which uses redness-based coloring).

Updated signature:
- Before: `DrawLobbyUI(screen, snakeCount int)`
- After: `DrawLobbyUI(screen, snakes, elapsedFrames, countdown)`

### Files Changed
- `internal/rendering/composite_renderer.go` - Added DrawScores call to DrawLobbyUI
- `internal/game/game.go` - Updated DrawLobbyUI call with new parameters
- `internal/game/game.go` - Updated GameRenderer interface
- `internal/game/game_test.go` - Updated mock to match new signature

### Test Status
✅ Existing tests still pass - score display logic tested via DrawScores method

### Remaining Issue
🔴 Need to investigate why redness fading isn't working. The Update() loop may not be running, or there may be an initialization issue with the game loop.

---
