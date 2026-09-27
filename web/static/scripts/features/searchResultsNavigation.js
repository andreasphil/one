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

/**
 * @param {Element} list
 * @param {1 | -1} offset
 */
function getRelative(list, offset) {
  const items = getItems(list);
  const current = getSelected(list);
  const index = current ? items.indexOf(current) + offset : offset > 0 ? 0 : -1;
  return items.at(index % items.length);
}

/** @param {Element | undefined} item */
function select(item) {
  if (item) location.replace(`#${item.id}`);
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

      const link = getSelected(list)?.querySelector(":scope > header a");
      if (!(link instanceof HTMLAnchorElement)) return;

      e.preventDefault();
      link.click();
    }
  });
}
