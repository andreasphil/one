import { RecentSearches } from "../components/recentSearches.js";
import { useListNavigation } from "../lib/listNavigation.js";
import { getSavedSearch } from "../state/savedSearch.js";

export class SearchView extends HTMLElement {
  static tag = "x-searchview";

  static define(tag = this.tag) {
    this.tag = tag;
    customElements.define(tag, this);
  }

  static get observedAttributes() {
    return [];
  }

  #store = getSavedSearch();

  // Lifecycle ----------------------------------------------

  #disconnect = new AbortController();

  constructor() {
    super();
  }

  connectedCallback() {
    this.#disconnect = new AbortController();
    const { signal } = this.#disconnect;

    useListNavigation(".search-results > li", { activate: ":scope > header a" });

    this.querySelector("form")?.addEventListener("submit", this, { signal });
    this.querySelector(RecentSearches.tag)?.addEventListener("remove", this, { signal });

    this.#render();
  }

  disconnectedCallback() {
    this.#disconnect.abort();
  }

  /** @param {Event} event */
  handleEvent(event) {
    if (event instanceof SubmitEvent) this.#handleSubmit(event);
    else if (event instanceof CustomEvent && event.type === "remove") this.#handleRemove(event);
  }

  #render() {
    const recent = this.querySelector(RecentSearches.tag);
    if (recent instanceof RecentSearches) recent.value = this.#store.searches();
  }

  // Internal -----------------------------------------------

  /** @param {SubmitEvent} event */
  #handleSubmit(event) {
    if (!(event.target instanceof HTMLFormElement)) return;

    const query = new FormData(event.target).get("query");
    if (typeof query === "string" && query) this.#store.saveSearch(query);
  }

  /** @param {CustomEvent<string>} event */
  #handleRemove(event) {
    this.#store.removeSearch(event.detail);
    this.#render();
  }
}
