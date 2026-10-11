import { getIcons } from "../state/icons.js";

/** @param {string} query */
function toSearchLink(query, base = location.href) {
  const url = new URL(base);
  url.searchParams.set("query", query);
  return url.toString();
}

export class RecentSearches extends HTMLElement {
  static tag = "x-recentsearches";

  static define(tag = this.tag) {
    this.tag = tag;
    customElements.define(tag, this);
  }

  #element;

  #fallback;

  /** @type {string[]} */
  #value = [];

  get value() {
    return this.#value;
  }

  /** @param {string[]} value */
  set value(value) {
    this.#value = value;
    this.#render();
  }

  // Lifecycle ----------------------------------------------

  #disconnect = new AbortController();

  constructor() {
    super();

    const { getElement: icon } = getIcons();

    this.#element = document.createElement("ul");

    this.#fallback = document.createElement("div");
    this.#fallback.setAttribute("fallback-for", "empty");

    const fallbackIcon = document.createElement("p");
    fallbackIcon.append(icon("Search"));

    const fallbackText = document.createElement("p");
    fallbackText.textContent = "No recent searches.";

    this.#fallback.append(fallbackIcon, fallbackText);
  }

  connectedCallback() {
    this.#disconnect = new AbortController();

    this.#element.addEventListener("click", this, { signal: this.#disconnect.signal });

    this.#render();
  }

  disconnectedCallback() {
    this.#disconnect.abort();
  }

  /** @param {Event} event */
  handleEvent(event) {
    if (!(event.target instanceof Element)) return;

    const button = event.target.closest("button[data-query]");
    if (!(button instanceof HTMLButtonElement)) return;

    this.dispatchEvent(new CustomEvent("remove", { detail: button.dataset.query }));
  }

  #render() {
    const { getElement: icon } = getIcons();

    const items = this.#value.map((i) => {
      const item = document.createElement("li");

      const link = document.createElement("a");
      link.setAttribute("href", toSearchLink(i));
      link.textContent = i;

      const button = document.createElement("button");
      button.setAttribute("type", "button");
      button.setAttribute("variant", "muted");
      button.setAttribute("size", "small");
      button.dataset.query = i;
      button.append(icon("X"));

      const label = document.createElement("span");
      label.classList.add("visually-hidden");
      label.textContent = `Remove "${i}"`;
      button.append(label);

      item.append(link, button);
      return item;
    });

    this.#element.replaceChildren(...items);
    this.replaceChildren(items.length ? this.#element : this.#fallback);
  }
}
