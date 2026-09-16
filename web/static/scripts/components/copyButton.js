import { getIcons } from "../state/icons.js";

/** @param {string | null} value */
function plainTextItem(value) {
  return new ClipboardItem({ "text/plain": new Blob([value ?? ""], { type: "text/plain" }) });
}

export class CopyButton extends HTMLElement {
  static tag = "x-copybutton";

  static define(tag = this.tag) {
    this.tag = tag;
    customElements.define(tag, this);
  }

  // State, refs --------------------------------------------

  #element;

  #didCopy = false;

  // Public API ---------------------------------------------

  get label() {
    return this.getAttribute("label");
  }

  /** Plain string value * */
  get value() {
    return this.getAttribute("value");
  }

  /**
   * Should be set to the ID of an HTML element. When set, the value is taken
   * from that HTML element. If it is a script tag with type `text/plain`,
   * the content is copied as plain text. If it is a regular HTML element, the
   * content is copied as rich text.
   */
  get valueFrom() {
    return this.getAttribute("valuefrom");
  }

  get copiedMessage() {
    return this.getAttribute("copiedmessage") ?? "Copied!";
  }

  get iconOnly() {
    return this.hasAttribute("icononly");
  }

  /** One of the variants supported by the button styles of the design system. */
  get variant() {
    return this.getAttribute("variant");
  }

  // Lifecycle ----------------------------------------------

  #disconnect = new AbortController();

  constructor() {
    super();

    const { getElement: icon } = getIcons();

    this.#element = document.createElement("button");
    this.#element.setAttribute("type", "button");

    const slot = document.createElement("span");
    slot.classList.add("slot");
    this.#element.append(slot);

    const label = document.createElement("span");
    label.classList.add("label");
    label.append(icon("Copy"));
    label.append(document.createElement("span"));
    slot.append(label);

    const copied = document.createElement("span");
    copied.classList.add("copied-message");
    copied.append(icon("Check"));
    copied.append(document.createElement("span"));
    slot.append(copied);
  }

  connectedCallback() {
    this.#disconnect = new AbortController();

    this.append(this.#element);

    this.#element.addEventListener("click", this, { signal: this.#disconnect.signal });

    if (this.valueFrom && this.value !== null) {
      console.warn("value and valuefrom are both set, valuefrom takes precedence", this);
    }

    this.#render();
  }

  disconnectedCallback() {
    this.#disconnect.abort();
  }

  async handleEvent() {
    try {
      await navigator.clipboard.write([this.#clipboardItem()]);
      this.#showDidCopy();
    } catch (e) {
      alert("Could not copy the value to the clipboard.");
      console.error("error copying the value to the clipboard", e);
    }
  }

  #render() {
    const label = this.#element.querySelector(".label > span");
    label.textContent = this.label;
    if (this.iconOnly) label.classList.add("visually-hidden");

    const copied = this.#element.querySelector(".copied-message > span");
    copied.textContent = this.copiedMessage;
    if (this.iconOnly) copied.classList.add("visually-hidden");

    this.#element.removeAttribute("variant");
    if (this.variant) this.#element.setAttribute("variant", this.variant);

    this.classList.remove("copied");
    if (this.#didCopy) this.classList.add("copied");
  }

  // Internal -----------------------------------------------

  #showDidCopy() {
    this.#didCopy = true;
    this.#render();

    setTimeout(() => {
      this.#didCopy = false;
      this.#render();
    }, 2000);
  }

  #clipboardItem() {
    if (!this.valueFrom) return plainTextItem(this.value);

    const source = document.getElementById(this.valueFrom);
    if (!source) throw new Error(`element with id ${this.valueFrom} does not exist`);
    else if (source instanceof HTMLScriptElement) {
      if (source.type !== "text/plain") {
        throw new Error(`script ${this.valueFrom} is not of type text/plain`);
      }

      return plainTextItem(source.textContent);
    }

    return new ClipboardItem({
      "text/html": new Blob([source.outerHTML], { type: "text/html" }),
      "text/plain": new Blob([source.innerText], { type: "text/plain" }),
    });
  }
}
