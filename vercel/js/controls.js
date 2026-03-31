/* ════════════════════════════════════════════════════════
   controls.js — DOM building, button wiring, event handlers
   Depends on: board.js, render.js, game.js
   ════════════════════════════════════════════════════════ */

// ─── Button references ───────────────────────────────────

const scoreToggleBtn = document.getElementById("scores-toggle");
const resetBtn = document.getElementById("reset-btn");

// ─── Build the board DOM ─────────────────────────────────

function buildDOM() {
  boardEl.innerHTML = "";
  scoreRowEl.innerHTML = "";
  colArrowsEl.innerHTML = "";
  colLayerEl.innerHTML = "";

  // ── Score hint row ──
  for (let c = 0; c < COLS; c++) {
    const sc = document.createElement("div");
    sc.className = "sc";
    sc.id = `sc-${c}`;
    scoreRowEl.appendChild(sc);
  }

  // ── Board cells (rows rendered top-to-bottom visually,
  //    but row 0 is the bottom of the logical board) ──
  for (let r = ROWS - 1; r >= 0; r--) {
    for (let c = 0; c < COLS; c++) {
      const cell = document.createElement("div");
      cell.className = "cell";
      cell.id = `cell-${r}-${c}`;
      boardEl.appendChild(cell);
    }
  }

  // ── Column hover highlight layer ──
  for (let c = 0; c < COLS; c++) {
    const hi = document.createElement("div");
    hi.className = "col-hi";
    hi.id = `chi-${c}`;
    colLayerEl.appendChild(hi);
  }

  // ── Column drop arrows ──
  for (let c = 0; c < COLS; c++) {
    const btn = document.createElement("button");
    btn.className = "arrow-btn";
    btn.innerHTML = "▾";
    btn.id = `arr-${c}`;
    btn.title = `Colonne ${c + 1}`;
    btn.addEventListener("click", () => handleColClick(c));
    colArrowsEl.appendChild(btn);
  }

  // ── Event delegation: clicks on board cells ──
  boardEl.addEventListener("click", (e) => {
    const cell = e.target.closest(".cell");
    if (!cell) return;
    const col = parseInt(cell.id.split("-")[2]);
    handleColClick(col);
  });

  // ── Hover tracking ──
  boardEl.addEventListener("mousemove", (e) => {
    const cell = e.target.closest(".cell");
    if (!cell) return;
    setHover(parseInt(cell.id.split("-")[2]));
  });
  boardEl.addEventListener("mouseleave", () => setHover(-1));

  colArrowsEl.addEventListener("mousemove", (e) => {
    const btn = e.target.closest(".arrow-btn");
    if (!btn) return;
    setHover(parseInt(btn.id.split("-")[1]));
  });
  colArrowsEl.addEventListener("mouseleave", () => setHover(-1));
}

// ─── Opponent selection ──────────────────────────────────

function setOpponent(opp) {
  opponent = opp;
  document.querySelectorAll("[data-opp]").forEach((btn) => {
    btn.className = "btn";
    if (btn.dataset.opp === opp) btn.classList.add("active-blue");
  });
  resetGame();
}

document.querySelectorAll("[data-opp]").forEach((btn) => {
  btn.addEventListener("click", () => setOpponent(btn.dataset.opp));
});

// ─── First player selection ──────────────────────────────

function setFirstPlayer(p) {
  firstPlayer = p;
  document.querySelectorAll("[data-first]").forEach((btn) => {
    btn.className = "btn";
    if (parseInt(btn.dataset.first) === p)
      btn.classList.add(p === 1 ? "active-blue" : "active-red");
  });
  resetGame();
}

document.querySelectorAll("[data-first]").forEach((btn) => {
  btn.addEventListener("click", () =>
    setFirstPlayer(parseInt(btn.dataset.first)),
  );
});

// ─── Reset button ────────────────────────────────────────

resetBtn.addEventListener("click", resetGame);

// ─── Score hints toggle ──────────────────────────────────

scoreToggleBtn.addEventListener("click", () => {
  showScores = !showScores;
  console.log("Toggle scores:", showScores);
  scoreRowEl.classList.toggle("hidden", !showScores);
  scoreToggleBtn.className = "btn btn-block";

  if (showScores) {
    scoreToggleBtn.classList.add("active-green");
    scoreToggleBtn.textContent = "Hide Scores";
    setTimeout(renderColScores, 20);
  } else {
    scoreToggleBtn.textContent = "Show Scores";
    clearColScores();
  }
});

// ─── Initialise ──────────────────────────────────────────

buildDOM();
setOpponent("player"); // also calls resetGame()
setFirstPlayer(1); // also calls resetGame()

// Re-apply active states after the double reset above
document.querySelector('[data-opp="player"]').classList.add("active-blue");
document.querySelector('[data-first="1"]').classList.add("active-blue");

// Load sequence from URL if present
loadSequenceFromURL();

// ─── Load sequence from URL parameter ────────────────────

/**
 * Check if there's a ?seq= parameter and load that sequence.
 */
function loadSequenceFromURL() {
  const url = new URL(window.location);
  const seq = url.searchParams.get("seq");
  if (seq) {
    replaySequence(seq);
  }
}
