import { debounce } from "../lib/debounce.js";

export class StickyTitle extends HTMLElement {
  static tag = "x-stickytitle";

  static define(tag = this.tag) {
    this.tag = tag;
    customElements.define(tag, this);
  }

  // State, refs --------------------------------------------

  #pill;

  /** @type {Element | null} */
  #current = null;

  #scheduleUpdate = debounce(() => this.#update(), 50, 200);

  // Public API ---------------------------------------------

  /** ID of the list whose items' titles are shown. */
  get for() {
    return this.getAttribute("for");
  }

  // Lifecycle ----------------------------------------------

  #disconnect = new AbortController();

  constructor() {
    super();

    this.#pill = document.createElement("div");
    this.#pill.classList.add("pill");
    this.#pill.setAttribute("aria-hidden", "true");
  }

  connectedCallback() {
    this.#disconnect = new AbortController();
    const { signal } = this.#disconnect;

    this.append(this.#pill);

    document.addEventListener("scroll", this, { capture: true, passive: true, signal });
    window.addEventListener("resize", this, { passive: true, signal });

    this.#update();
  }

  disconnectedCallback() {
    this.#disconnect.abort();
    this.#scheduleUpdate.cancel();
  }

  handleEvent() {
    this.#scheduleUpdate();
  }

  // Internal -----------------------------------------------

  #update() {
    const id = this.for;
    if (!id) return;

    const list = document.getElementById(id);
    if (!list) return;

    const items = [...list.querySelectorAll(":scope > li")];
    if (!items.length) return;

    const { top } = this.getBoundingClientRect();
    const scrollMargin = parseFloat(getComputedStyle(items[0]).scrollMarginTop) || 0;
    const line = Math.max(top + this.#pill.offsetHeight, scrollMargin);
    const minHeight = window.innerHeight * 0.3;

    const next =
      items.find((i) => {
        const rect = i.getBoundingClientRect();
        if (rect.height <= minHeight || rect.top >= line || rect.bottom <= line) return false;

        const header = i.querySelector(":scope > header");
        return header && header.getBoundingClientRect().bottom < top;
      }) ?? null;

    if (next === this.#current) return;

    this.#current = next;
    this.classList.toggle("active", !!next);
    if (next) this.#pill.textContent = next.querySelector(":scope > header a")?.textContent?.trim() ?? "";
  }
}
