/* ═══════════════════════════════════════════════════
   render.js — All DOM rendering (board, status, hints)
   Depends on: board.js, solver.js, state (game.js globals)
   ═══════════════════════════════════════════════════ */

// ─── Element references ───────────────────────────────────

const boardEl = document.getElementById("board");
const scoreRowEl = document.getElementById("score-row");
const colArrowsEl = document.getElementById("col-arrows");
const colLayerEl = document.getElementById("col-layer");
const previewTokenEl = document.getElementById("preview-token");
const statusBar = document.getElementById("status-bar");
const statusDot = document.getElementById("s-dot");
const statusText = document.getElementById("status-text");
const thinkingEl = document.getElementById("thinking");

// ─── Board rendering ─────────────────────────────────────

/**
 * Sync the visual grid with the current board state.
 * Pass newRow/newCol to trigger the drop animation on that cell.
 */
function renderBoard(newRow = -1, newCol = -1) {
  for (let r = 0; r < ROWS; r++) {
    for (let c = 0; c < COLS; c++) {
      const cell = document.getElementById(`cell-${r}-${c}`);
      cell.className = "cell";
      if (board[r][c] === 1) cell.classList.add("p1");
      if (board[r][c] === 2) cell.classList.add("p2");
      if (winCells.some(([wr, wc]) => wr === r && wc === c))
        cell.classList.add("win-cell");

      // Drop animation on the newly placed cell
      if (r === newRow && c === newCol && board[r][c] !== 0) {
        const visualRow = ROWS - 1 - r; // 0 = top
        const cellSizePx =
          document.getElementById("cell-0-0")?.offsetWidth || 64;
        const dropH = (visualRow + 1) * (cellSizePx + 5) + "px";
        cell.style.setProperty("--drop-h", dropH);
        cell.classList.add("dropping");
        cell.addEventListener(
          "animationend",
          () => cell.classList.remove("dropping"),
          { once: true },
        );
      }
    }
  }

  // Disable arrows when the column is full or the game is over
  for (let c = 0; c < COLS; c++) {
    const arr = document.getElementById(`arr-${c}`);
    arr.disabled = gameOver || aiThinking || !isValid(board, c);
  }
}

// ─── Status bar ──────────────────────────────────────────

function renderStatus() {
  statusBar.className = "status-bar";

  if (gameOver) {
    const winner = checkWinner(board, 1) ? 1 : checkWinner(board, 2) ? 2 : 0;
    if (winner === 0) {
      statusBar.classList.add("draw");
      statusDot.className = "s-dot draw";
      statusText.textContent = "It's a draw!";
    } else {
      statusBar.classList.add(`win-p${winner}`);
      statusDot.className = "s-dot win";
      statusText.textContent = `${getPlayerName(winner)} is the winner!`;
    }
  } else {
    statusBar.classList.add(`p${currentPlayer}-turn`);
    statusDot.className = `s-dot p${currentPlayer}`;
    statusText.textContent = `${getPlayerName(currentPlayer)}'s turn`;
  }
}

// ─── Column score hints ───────────────────────────────────

function renderColScores() {
  if (!showScores) return;

  const raw = computeColScores(board, currentPlayer);

  const valid = raw.filter((s) => s !== null);

  // Handle case where all scores are null/invalid
  if (valid.length === 0) {
    for (let c = 0; c < COLS; c++) {
      const el = document.getElementById(`sc-${c}`);
      if (el) {
        el.textContent = "-";
        el.className = "sc score-neutral";
      }
    }
    return;
  }

  const maxS = Math.max(...valid);
  const minS = Math.min(...valid);

  for (let c = 0; c < COLS; c++) {
    const el = document.getElementById(`sc-${c}`);
    el.className = "sc";

    if (raw[c] === null) {
      el.textContent = "—";
      el.classList.add("score-invalid");
      continue;
    }

    const v = raw[c];
    el.textContent = v;

    let cls;
    if (v === maxS && v > 0) cls = "score-best";
    else if (v > 0) cls = "score-good";
    else if (v === 0) cls = "score-neutral";
    else if (v < 0 && v !== minS) cls = "score-bad";
    else cls = "score-worst";
    el.classList.add(cls);

    if (v === maxS) {
      const img = document.createElement("img");
      img.src = "svg/star.svg";
      img.alt = "Best move";
      img.className = "score-star";
      el.appendChild(img);
    }
  }
}

/** Blank out all score hint cells. */
function clearColScores() {
  for (let c = 0; c < COLS; c++) {
    const el = document.getElementById(`sc-${c}`);
    if (el) {
      el.textContent = "";
      el.className = "sc";
    }
  }
}

// ─── Thinking indicator ──────────────────────────────────

function setThinking(on) {
  thinkingEl.classList.toggle("show", on);
}

// ─── Column hover ────────────────────────────────────────

let _hoveredCol = -1;

function setHover(col) {
  _hoveredCol = col;

  // Update column highlights
  for (let c = 0; c < COLS; c++) {
    const hi = document.getElementById(`chi-${c}`);
    if (hi) hi.classList.toggle("hovered", c === col);
  }

  // Show/hide preview token
  if (col === -1 || gameOver || aiThinking || !isValid(board, col)) {
    previewTokenEl.classList.remove("show");
  } else {
    const landingRow = getLowestEmpty(board, col);
    if (landingRow === -1) {
      previewTokenEl.classList.remove("show");
      return;
    }

    // Get position from the actual cell in the landing row
    const cell = document.getElementById(`cell-${landingRow}-${col}`);
    if (!cell) {
      previewTokenEl.classList.remove("show");
      return;
    }

    // Account for board-wrap padding
    previewTokenEl.style.left = cell.offsetLeft + "px";
    previewTokenEl.style.top = cell.offsetTop + "px";
    previewTokenEl.className = `preview-token p${currentPlayer} show`;
  }
}
