/* ════════════════════════════════════════════════════════
   modal.js — Modal handling and section navigation
   ════════════════════════════════════════════════════════ */

const modalOverlay = document.getElementById("modal-overlay");
const modal = document.getElementById("modal");
const modalCloseBtn = document.getElementById("modal-close");
const infoBtn = document.getElementById("info-btn");
const modalTabs = document.querySelectorAll(".modal-tab");
const modalSections = document.querySelectorAll(".modal-section");

// ── Open modal ──
function openModal() {
  modalOverlay.classList.add("show");
}

// ── Close modal ──
function closeModal() {
  modalOverlay.classList.remove("show");
}

// ── Switch section ──
function switchSection(sectionId) {
  // Remove active from all tabs and sections
  modalTabs.forEach((tab) => tab.classList.remove("active"));
  modalSections.forEach((section) => section.classList.remove("active"));

  // Add active to the clicked tab and corresponding section
  document
    .querySelector(`[data-section="${sectionId}"]`)
    ?.classList.add("active");
  document.getElementById(sectionId)?.classList.add("active");
}

// ── Event listeners ──
infoBtn.addEventListener("click", openModal);
modalCloseBtn.addEventListener("click", closeModal);

// Close modal when clicking on overlay (not on modal content)
modalOverlay.addEventListener("click", (e) => {
  if (e.target === modalOverlay) closeModal();
});

// Close modal with Escape key
document.addEventListener("keydown", (e) => {
  if (e.key === "Escape") closeModal();
});

// Tab navigation
modalTabs.forEach((tab) => {
  tab.addEventListener("click", () => {
    const sectionId = tab.dataset.section;
    switchSection(sectionId);
  });
});
