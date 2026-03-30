/* ═══════════════════════════════════════════════════
   game.js — Game state, flow, and shared globals
   Depends on: board.js, solver.js, render.js
   ═══════════════════════════════════════════════════ */

// ─── Shared state (read by render.js and controls.js) ────

let board = createBoard();
let currentPlayer = 1;
let gameOver = false;
let aiThinking = false;
let winCells = [];
let moveHistory = ""; // Track move sequence for WASM

let opponent = "player";
let firstPlayer = 1;
let showScores = false;
let scores = { 1: 0, 2: 0 };

// ─── Helpers ─────────────────────────────────────────────

function getPlayerName(p) {
  if (p === 1) return "Player 1";
  const names = {
    player: "Player 2",
    easy: "Easy AI",
    hard: "Hard AI",
    perfect: "Perfect AI",
  };
  return names[opponent] ?? "Player 2";
}

// ─── Core game flow ──────────────────────────────────────

/**
 * Called when a human clicks a column.
 * Ignored during AI thinking, game over, or when it is the AI's turn.
 */
function handleColClick(col) {
  if (gameOver || aiThinking) return;
  if (!isValid(board, col)) return;
  if (currentPlayer === 2 && opponent !== "player") return;
  placePiece(col);
}

/**
 * Place a piece in `col` for `currentPlayer`, then check for
 * win / draw and hand off to the AI if needed.
 */
function placePiece(col) {
  const r = drop(board, col, currentPlayer);
  if (r === -1) return;

  moveHistory += col;

  renderBoard(r, col);

  // ── Win check ──
  const cells = checkWinner(board, currentPlayer);
  if (cells) {
    winCells = cells;
    scores[currentPlayer]++;
    gameOver = true;
    renderBoard(); // re-render to show glow
    renderStatus();
    renderScoreboard();
    return;
  }

  // ── Draw check ──
  if (isDraw(board)) {
    gameOver = true;
    renderStatus();
    return;
  }

  // ── Switch turn ──
  currentPlayer = currentPlayer === 1 ? 2 : 1;
  renderStatus();

  // Show hint scores if enabled and it's a human's turn next
  const nextIsAI = currentPlayer === 2 && opponent !== "player";
  if (showScores && !nextIsAI) {
    setTimeout(renderColScores, 20);
  }

  // ── Trigger AI ──
  if (nextIsAI) triggerAI();
}

/**
 * Start an asynchronous AI move (with a brief delay so the UI
 * can paint before the heavy computation begins).
 */
function triggerAI() {
  aiThinking = true;
  setThinking(true);
  clearColScores();
  renderBoard(); // disable arrows while thinking

  const delay =
    opponent === "perfect" ? 80 : opponent === "impossible" ? 50 : 30;

  setTimeout(() => {
    const col = getAIMove(board, opponent, 2);
    aiThinking = false;
    setThinking(false);
    placePiece(col);
    if (showScores && !gameOver) setTimeout(renderColScores, 20);
  }, delay);
}

/**
 * Reset the board to its initial state, keeping scores intact.
 * Called by the controls and automatically after changing settings.
 */
function resetGame() {
  board = createBoard();
  moveHistory = "";
  currentPlayer = firstPlayer;
  gameOver = false;
  aiThinking = false;
  winCells = [];

  setThinking(false);
  renderBoard();
  renderStatus();
  clearColScores();

  if (showScores) setTimeout(renderColScores, 20);

  // If the AI goes first
  if (currentPlayer === 2 && opponent !== "player") {
    setTimeout(triggerAI, 300);
  }
}
