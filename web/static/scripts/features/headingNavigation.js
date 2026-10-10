import { isEditable } from "../lib/keyboard.js";

const SELECTOR = "#note h2[id]";

function getItems() {
  return [...document.querySelectorAll(SELECTOR)];
}

/** @param {Element} item */
function isInView(item) {
  const { top, bottom } = item.getBoundingClientRect();
  return bottom > 0 && top < window.innerHeight;
}

/** True while a scroll triggered by select() is still in progress. */
let navigating = false;
let navigatingTimeout = 0;

/**
 * The selected heading, unless the user has scrolled it out of view. While a keyboard-triggered
 * scroll is running, positions are in flux, so the selection is trusted as is.
 */
function getCurrent() {
  const selected = document.querySelector(`${SELECTOR}:target`);
  return selected && (navigating || isInView(selected)) ? selected : null;
}

/** @param {1 | -1} offset */
function getRelative(offset) {
  const items = getItems();
  const current = getCurrent();
  if (current) return items[items.indexOf(current) + offset];

  // Headings land at the scroll margin when navigated to, so that's the line
  // between "above" and "visible". The pixel of slack absorbs rounding.
  const line = (parseFloat(getComputedStyle(items[0] ?? document.body).scrollMarginTop) || 0) - 1;
  const top = (/** @type {Element} */ item) => item.getBoundingClientRect().top;

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

export function init() {
  document.addEventListener("keydown", (e) => {
    if (e.altKey || e.ctrlKey || e.metaKey || isEditable(e.target)) return;
    if (!document.querySelector("#note")) return;

    if (e.key === "j") {
      e.preventDefault();
      select(getRelative(1));
    } else if (e.key === "k") {
      e.preventDefault();
      select(getRelative(-1));
    } else if (e.key === "J") {
      e.preventDefault();
      select(getItems()[0]);
    } else if (e.key === "K") {
      e.preventDefault();
      select(getItems().at(-1));
    }
  });
}
