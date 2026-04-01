/* ═══════════════════════════════════════════════════════
   solver.js — WebAssembly-based game solver
   Uses gonnect4.wasm for AI moves and analysis
   Depends on: board.js
   ═══════════════════════════════════════════════════════ */

let wasmReady = false;
let wasmInitPromise = null;

/**
 * Initialize the WebAssembly module
 */
async function initWasm() {
  if (wasmInitPromise) return wasmInitPromise;

  wasmInitPromise = (async () => {
    if (typeof Go === "undefined") {
      console.error("Go runtime not found. Make sure wasm_exec.js is loaded.");
      return false;
    }

    try {
      const go = new Go();
      const result = await WebAssembly.instantiateStreaming(
        fetch("js/gonnect4.wasm"),
        go.importObject,
      );
      go.run(result.instance);

      if (typeof window.gonnectAnalyze !== "function") {
        console.error("WASM functions not exposed properly");
        return false;
      }

      wasmReady = true;
      return true;
    } catch (err) {
      console.error("Failed to initialize WASM:", err);
      return false;
    }
  })();

  return wasmInitPromise;
}

/**
 * Return the best column for the AI to play.
 * @param {Int8Array[]} b        - Current board
 * @param {string}      diff     - Difficulty key
 */
function getAIMove(b, diff) {
  const valid = getValidCols(b);
  // For "easy", pick a random valid move
  if (diff === "easy") {
    return valid[Math.floor(Math.random() * valid.length)];
  }

  const sequence = Array.from(moveHistory)
    .map((col) => String(Number(col) + 1))
    .join("");

  const response = window.gonnectAnalyze(sequence);
  let scores = response.analysis.scores;

  // For "hard", the higher a column's score, the greater its chance of being chosen.
  if (diff === "hard") {
    const temp = 0.5;
    const candidates = scores
      .map((score, col) => ({ col, score }))
      .filter(({ score }) => score !== -128);

    const maxScore = Math.max(...candidates.map((c) => c.score));
    const exps = candidates.map((c) => Math.exp((c.score - maxScore) / temp));
    const sumExp = exps.reduce((a, b) => a + b, 0);
    const probs = exps.map((e) => e / sumExp);

    const r = Math.random();
    let cumulative = 0;

    for (let i = 0; i < candidates.length; i++) {
      cumulative += probs[i];
      if (r < cumulative) {
        return candidates[i].col;
      }
    }
  }

  const topCandidates = valid
    .filter((col) => scores[col] !== undefined && scores[col] !== -128)
    .sort((a, b) => {
      return scores[b] - scores[a];
    });
  // For "perfect", pick the best move
  return topCandidates[0];
}

/**
 * Compute per-column evaluation scores for the hint display.
 * Returns an array of 7 values (null = invalid column).
 * @param {Int8Array[]} b      - Current board
 * @param {number}      player - Whose turn it is
 */
function computeColScores(b, player) {
  if (!wasmReady) {
    console.warn("WASM not ready for score computation");
    return Array(COLS).fill(null);
  }

  try {
    const sequence = Array.from(moveHistory)
      .map((col) => String(Number(col) + 1))
      .join("");
    const analysis = window.gonnectAnalyze(sequence);

    if (!analysis.ok) {
      console.warn("Analysis not ok:", analysis.error);
      return Array(COLS).fill(null);
    }

    if (!analysis.analysis) {
      console.warn("No analysis.analysis field");
      return Array(COLS).fill(null);
    }

    const scores = analysis.analysis.scores;

    if (!scores) {
      console.warn("No scores returned");
      return Array(COLS).fill(null);
    }

    const result = Array.from({ length: COLS }, (_, col) => {
      if (!isValid(b, col)) return null;
      // Handle both array index and object key
      let score = scores[col];
      if (score === undefined) {
        score = scores[col.toString()];
      }
      return score ?? null;
    });
    return result;
  } catch (err) {
    console.error("Error computing column scores:", err);
    console.error("Stack:", err.stack);
    return Array(COLS).fill(null);
  }
}

// Initialize WASM when the page loads
document.addEventListener("DOMContentLoaded", initWasm);
