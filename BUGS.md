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
