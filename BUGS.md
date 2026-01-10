# Bug Reports - Snakehem Refactoring

This file tracks bugs found during testing of the refactored codebase.

---

## Bug #1: Lobby State Shows Action UI Elements

**Date**: 2026-01-10
**Status**: 🔴 OPEN
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
Remove the `DrawActionUI()` call from the LobbyState case in game.Draw().

### Test Scenario
A test could verify that in lobby state, only lobby-specific UI elements are rendered and action UI elements are not.

---
