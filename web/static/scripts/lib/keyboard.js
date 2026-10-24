/** @param {EventTarget | null} target */
export function isEditable(target) {
  return target instanceof HTMLElement && target.matches("input, textarea, [contenteditable]");
}

/** @typedef {Partial<Pick<KeyboardEvent, "key" | "metaKey" | "altKey" | "ctrlKey" | "shiftKey">>} KeyboardShortcut */

/**
 * @param {KeyboardShortcut} shortcut
 * @param {KeyboardEvent} event
 * @returns boolean
 */
export function matchShortcut(shortcut, event) {
  const strictShortcut = Object.assign(
    {
      altKey: false,
      ctrlKey: false,
      metaKey: false,
      shiftKey: false,
      key: "",
    },
    shortcut,
  );

  let match = Object.entries(strictShortcut).reduce((match, [key, value]) => {
    return match && event[/** @type {keyof KeyboardEvent} */ (key)] === value;
  }, true);

  return match;
}
