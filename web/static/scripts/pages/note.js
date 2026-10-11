import { useCommands } from "../lib/commands.js";
import { isEditable } from "../lib/keyboard.js";
import { useListNavigation } from "../lib/listNavigation.js";
import { getIcons } from "../state/icons.js";

const maximizedKey = "maximize_note";

export class NoteView extends HTMLElement {
  static tag = "x-noteview";

  static define(tag = this.tag) {
    this.tag = tag;
    customElements.define(tag, this);
  }

  static get observedAttributes() {
    return [];
  }

  get maximized() {
    return this.hasAttribute("maximized");
  }

  set maximized(value) {
    this.toggleAttribute("maximized", value);

    if (value) localStorage.setItem(maximizedKey, "true");
    else localStorage.removeItem(maximizedKey);
  }

  toggleMaximized(value = !this.maximized) {
    this.maximized = value;
  }

  // Lifecycle ----------------------------------------------

  #disconnect = new AbortController();

  #removeCommands = () => {};

  constructor() {
    super();
  }

  connectedCallback() {
    this.#disconnect = new AbortController();

    this.maximized = localStorage.getItem(maximizedKey) === "true";

    useListNavigation("#note h2[id]");

    document.addEventListener("keydown", this, { signal: this.#disconnect.signal });

    const { getElement: icon } = getIcons();
    const { register } = useCommands();

    this.#removeCommands = register({
      id: "note:fullscreen",
      name: "Toggle fullscreen",
      groupName: "Note",
      chord: "f",
      icon: icon("Maximize2"),
      action: () => this.toggleMaximized(),
    });
  }

  disconnectedCallback() {
    this.#disconnect.abort();
    this.#removeCommands();
  }

  /** @param {Event} event */
  handleEvent(event) {
    if (event instanceof KeyboardEvent) this.#handleKeydown(event);
  }

  // Internal -----------------------------------------------

  /** @param {KeyboardEvent} event */
  #handleKeydown(event) {
    if (event.key !== "f" || isEditable(event.target)) return;
    this.toggleMaximized();
  }
}
