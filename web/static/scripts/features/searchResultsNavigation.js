import { isEditable } from "../lib/keyboard.js";

/** @param {Element} list */
function getSelected(list) {
  const target = list.querySelector(":scope > li:target");
  return target instanceof HTMLLIElement ? target : null;
}

/** @param {Element} list */
function getItems(list) {
  return [...list.querySelectorAll(":scope > li")];
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
 * The selected item, unless the user has scrolled it out of view. While a
 * keyboard-triggered scroll is running, positions are in flux, so the
 * selection is trusted as is.
 *
 * @param {Element} list
 */
function getCurrent(list) {
  const selected = getSelected(list);
  return selected && (navigating || isInView(selected)) ? selected : null;
}

/**
 * @param {Element} list
 * @param {1 | -1} offset
 */
function getRelative(list, offset) {
  const items = getItems(list);
  const current = getCurrent(list);
  if (current) return items[items.indexOf(current) + offset];

  // Headings land at the scroll margin when navigated to, so that's the line
  // between "above" and "visible". The pixel of slack absorbs rounding.
  const line = (parseFloat(getComputedStyle(items[0] ?? list).scrollMarginTop) || 0) - 1;
  const headingTop = (/** @type {Element} */ item) =>
    (item.querySelector(":scope > header") ?? item).getBoundingClientRect().top;

  return offset > 0
    ? items.find((i) => headingTop(i) >= line)
    : items.findLast((i) => headingTop(i) < line);
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
  addEventListener("pageshow", () => {
    document.querySelector(".search-results > li:target")?.scrollIntoView({ behavior: "instant" });
  });

  document.addEventListener("keydown", (e) => {
    if (e.altKey || e.ctrlKey || e.metaKey || isEditable(e.target)) return;

    const list = document.querySelector(".search-results");
    if (!list) return;

    if (e.key === "j") {
      e.preventDefault();
      select(getRelative(list, 1));
    } else if (e.key === "k") {
      e.preventDefault();
      select(getRelative(list, -1));
    } else if (e.key === "J") {
      e.preventDefault();
      select(getItems(list)[0]);
    } else if (e.key === "K") {
      e.preventDefault();
      select(getItems(list).at(-1));
    } else if (e.key === "Enter" && !e.shiftKey) {
      if (e.target instanceof HTMLElement && e.target.matches("a, button")) return;

      const link = getCurrent(list)?.querySelector(":scope > header a");
      if (!(link instanceof HTMLAnchorElement)) return;

      e.preventDefault();
      link.click();
    }
  });
}
