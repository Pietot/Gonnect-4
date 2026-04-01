/* ════════════════════════════════════════════════════
   board.js — Pure board logic (no DOM, no rendering)
   ════════════════════════════════════════════════════ */

const ROWS = 6;
const COLS = 7;

/** Priority column order for search (center-first). */
const MOVE_ORDER = [3, 2, 4, 1, 5, 0, 6];

/** Create a fresh empty board (row 0 = bottom). */
function createBoard() {
  return Array.from({ length: ROWS }, () => new Int8Array(COLS));
}

/** Deep-clone a board. */
function cloneBoard(b) {
  return b.map((r) => r.slice());
}

/**
 * Return the lowest empty row in a column, or -1 if full.
 * Row 0 is the bottom of the board.
 */
function getLowestEmpty(b, col) {
  for (let r = 0; r < ROWS; r++) {
    if (b[r][col] === 0) return r;
  }
  return -1;
}

/** Is this column a valid move? */
function isValid(b, col) {
  return col >= 0 && col < COLS && b[ROWS - 1][col] === 0;
}

/** Return all valid columns (in MOVE_ORDER). */
function getValidCols(b) {
  return MOVE_ORDER.filter((c) => isValid(b, c));
}

/**
 * Drop a piece for `player` in `col`.
 * Mutates `b` and returns the row index, or -1 if full.
 */
function drop(b, col, player) {
  const r = getLowestEmpty(b, col);
  if (r === -1) return -1;
  b[r][col] = player;
  return r;
}

/**
 * Check whether `player` has won.
 * Returns an array of 4 [row, col] pairs if yes, otherwise null.
 */
function checkWinner(b, player) {
  // Horizontal
  for (let r = 0; r < ROWS; r++) {
    for (let c = 0; c <= COLS - 4; c++) {
      if (
        b[r][c] === player &&
        b[r][c + 1] === player &&
        b[r][c + 2] === player &&
        b[r][c + 3] === player
      ) {
        return [
          [r, c],
          [r, c + 1],
          [r, c + 2],
          [r, c + 3],
        ];
      }
    }
  }
  // Vertical
  for (let c = 0; c < COLS; c++) {
    for (let r = 0; r <= ROWS - 4; r++) {
      if (
        b[r][c] === player &&
        b[r + 1][c] === player &&
        b[r + 2][c] === player &&
        b[r + 3][c] === player
      ) {
        return [
          [r, c],
          [r + 1, c],
          [r + 2, c],
          [r + 3, c],
        ];
      }
    }
  }
  // Diagonal /
  for (let r = 0; r <= ROWS - 4; r++) {
    for (let c = 0; c <= COLS - 4; c++) {
      if (
        b[r][c] === player &&
        b[r + 1][c + 1] === player &&
        b[r + 2][c + 2] === player &&
        b[r + 3][c + 3] === player
      ) {
        return [
          [r, c],
          [r + 1, c + 1],
          [r + 2, c + 2],
          [r + 3, c + 3],
        ];
      }
    }
  }
  // Diagonal \
  for (let r = 3; r < ROWS; r++) {
    for (let c = 0; c <= COLS - 4; c++) {
      if (
        b[r][c] === player &&
        b[r - 1][c + 1] === player &&
        b[r - 2][c + 2] === player &&
        b[r - 3][c + 3] === player
      ) {
        return [
          [r, c],
          [r - 1, c + 1],
          [r - 2, c + 2],
          [r - 3, c + 3],
        ];
      }
    }
  }
  return null;
}

/** Is the board completely full (draw)? */
function isDraw(b) {
  return b[ROWS - 1].every((v) => v !== 0);
}
