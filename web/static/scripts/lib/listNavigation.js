import { isEditable, matchShortcut } from "./keyboard.js";

/** @param {Element} item */
function isInView(item) {
  const { top, bottom } = item.getBoundingClientRect();
  return bottom > 0 && top < window.innerHeight;
}

/** True while a scroll triggered by select() is still in progress. */
let navigating = false;
let navigatingTimeout = 0;

/**
 * The selected item, unless the user has scrolled it out of view. While a keyboard-triggered scroll
 * is running, positions are in flux, so the selection is trusted as is.
 *
 * @param {string} selector
 */
function getCurrent(selector) {
  const selected = document.querySelector(`${selector}:target`);
  return selected && (navigating || isInView(selected)) ? selected : null;
}

/**
 * @param {string} selector
 * @param {1 | -1} offset
 */
function getRelative(selector, offset) {
  const items = [...document.querySelectorAll(selector)];
  const current = getCurrent(selector);
  if (current) return items[items.indexOf(current) + offset];

  // Items land at the scroll margin when navigated to, so that's the line
  // between "above" and "visible". The pixel of slack absorbs rounding.
  const line = (parseFloat(getComputedStyle(items[0] ?? document.body).scrollMarginTop) || 0) - 1;
  const top = (/** @type {Element} */ item) =>
    (item.querySelector(":scope > header") ?? item).getBoundingClientRect().top;

  return offset > 0 ? items.find((i) => top(i) >= line) : items.findLast((i) => top(i) < line);
}

/** @param {Element | undefined} item */
function select(item) {
  if (!item) return;

  navigating = true;
  clearTimeout(navigatingTimeout);
  const done = () => {
    navigating = false;
    clearTimeout(navigatingTimeout);
  };
  addEventListener("scrollend", done, { once: true });
  navigatingTimeout = setTimeout(done, 600);

  location.replace(`#${item.id}`);
}

/**
 * Keyboard navigation (j/k/J/K) between the elements matching a selector, using the URL fragment
 * as the selection.
 *
 * @param {string} selector Matches the navigable items. Each needs an id.
 * @param {{ activate?: string }} [options] If set, Enter clicks the link matching this selector
 * inside the current item.
 */
export function useListNavigation(selector, { activate } = {}) {
  document.addEventListener("keydown", (e) => {
    if (isEditable(e.target)) return;

    const items = [...document.querySelectorAll(selector)];
    if (!items.length) return;

    if (matchShortcut({ key: "j" }, e)) {
      e.preventDefault();
      select(getRelative(selector, 1));
    } else if (matchShortcut({ key: "k" }, e)) {
      e.preventDefault();
      select(getRelative(selector, -1));
    } else if (matchShortcut({ key: "J", shiftKey: true }, e)) {
      e.preventDefault();
      select(items.at(-1));
    } else if (matchShortcut({ key: "K", shiftKey: true }, e)) {
      e.preventDefault();
      select(items[0]);
    } else if (matchShortcut({ key: "Enter" }, e) && activate) {
      if (e.target instanceof HTMLElement && e.target.matches("a, button")) return;

      const link = getCurrent(selector)?.querySelector(activate);
      if (!(link instanceof HTMLAnchorElement)) return;

      e.preventDefault();
      link.click();
    }
  });
}
