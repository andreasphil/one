export function init() {
  document.addEventListener("keydown", (e) => {
    if (e.target instanceof HTMLElement && e.target.matches("input, textarea, [contenteditable]")) {
      return;
    }

    const el = document.querySelector(`[data-focuskey=${CSS.escape(e.key)}]`);
    if (!(el instanceof HTMLInputElement)) return;

    e.preventDefault();
    el.focus();
    el.setSelectionRange(0, el.value.length);
  });
}
