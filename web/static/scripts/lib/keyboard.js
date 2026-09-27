/** @param {EventTarget | null} target */
export function isEditable(target) {
  return target instanceof HTMLElement && target.matches("input, textarea, [contenteditable]");
}
