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
      console.log("WASM initialized successfully");
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
 * @param {number}      aiPlayer - 1 or 2
 */
function getAIMove(b, diff, aiPlayer) {
  if (!wasmReady) {
    // Fallback to random move if WASM not ready
    const valid = getValidCols(b);
    return valid[Math.floor(Math.random() * valid.length)];
  }

  const valid = getValidCols(b);

  if (diff === "easy") {
    return valid[Math.floor(Math.random() * valid.length)];
  }

  try {
    const sequence = Array.from(moveHistory).map(col => String(Number(col) + 1)).join("");
    const analysis = window.gonnectAnalyze(sequence);
    console.log("AI Analysis:", analysis);

    if (!analysis.ok || !analysis.analysis || !analysis.analysis.scores) {
      // Fallback to random
      return valid[Math.floor(Math.random() * valid.length)];
    }

    const scores = analysis.analysis.scores;

    const topCandidates = valid
      .filter((col) => scores[col] !== undefined && scores[col] !== -1)
      .sort((a, b) => {
        const scoreA = scores[a] ?? -Infinity;
        const scoreB = scores[b] ?? -Infinity;
        return scoreB - scoreA;
      });

    if (topCandidates.length === 0) {
      return valid[0];
    }

    // For "hard", pick the best move
    if (diff === "hard") {
      return topCandidates[0];
    }

    // For "perfect", pick from top moves with some randomness
    const numTopMoves = Math.max(1, Math.floor(topCandidates.length / 3));
    return topCandidates[Math.floor(Math.random() * Math.min(1, numTopMoves))];
  } catch (err) {
    console.error("Error in getAIMove:", err);
    return valid[0];
  }
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
    const sequence = Array.from(moveHistory).map(col => String(Number(col) + 1)).join("");
    console.log("Requesting analysis for sequence:", sequence);
    const analysis = window.gonnectAnalyze(sequence);
    console.log("Full analysis object:", analysis);
    console.log("analysis.analysis:", analysis.analysis);

    if (!analysis.ok) {
      console.warn("Analysis not ok:", analysis.error);
      return Array(COLS).fill(null);
    }

    if (!analysis.analysis) {
      console.warn("No analysis.analysis field");
      return Array(COLS).fill(null);
    }

    const scores = analysis.analysis.scores;
    console.log("Scores from WASM:", scores, "type:", typeof scores);

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
    console.log("Final column scores:", result);
    return result;
  } catch (err) {
    console.error("Error computing column scores:", err);
    console.error("Stack:", err.stack);
    return Array(COLS).fill(null);
  }
}

// Initialize WASM when the page loads
document.addEventListener("DOMContentLoaded", initWasm);