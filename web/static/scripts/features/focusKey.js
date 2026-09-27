import { isEditable } from "../lib/keyboard.js";

export function init() {
  document.addEventListener("keydown", (e) => {
    if (isEditable(e.target)) return;

    const el = document.querySelector(`[data-focuskey=${CSS.escape(e.key)}]`);
    if (!(el instanceof HTMLInputElement)) return;

    e.preventDefault();
    el.focus();
    el.setSelectionRange(0, el.value.length);
  });
}
