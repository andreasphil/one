import { useListNavigation } from "../lib/listNavigation.js";

export class TagsView extends HTMLElement {
  static tag = "x-tagsview";

  static define(tag = this.tag) {
    this.tag = tag;
    customElements.define(tag, this);
  }

  static get observedAttributes() {
    return [];
  }

  // Lifecycle ----------------------------------------------

  constructor() {
    super();
  }

  connectedCallback() {
    useListNavigation(".search-results > li", { activate: ":scope > header a" });
  }
}
