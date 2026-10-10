import { isEditable } from "../lib/keyboard.js";

export function init() {
  document.addEventListener("keydown", (event) => {
    if (isEditable(event.target)) return;

    const el = document.querySelector(`[data-focuskey=${CSS.escape(event.key)}]`);
    if (!el) return;
    event.preventDefault();
    el.focus();

    if (el instanceof HTMLInputElement) {
      el.setSelectionRange(0, el.value.length);
    } else if (el instanceof HTMLButtonElement) {
      el.click();
    }
  });
}
