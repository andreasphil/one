import { getIcons } from "../state/icons.js";

export class CopyButton extends HTMLElement {
  static tag = "x-copybutton";

  static define(tag = this.tag) {
    this.tag = tag;
    customElements.define(tag, this);
  }

  #element;

  #didCopy = false;

  get label() {
    return this.getAttribute("label");
  }

  get value() {
    return this.getAttribute("value");
  }

  get copiedMessage() {
    return this.getAttribute("copiedmessage") ?? "Copied!";
  }

  #disconnect = new AbortController();

  constructor() {
    super();

    const { getElement: icon } = getIcons();

    this.#element = document.createElement("button");
    this.#element.setAttribute("variant", "muted");
    this.#element.setAttribute("type", "button");
    this.append(this.#element);

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

    this.#render();
  }

  connectedCallback() {
    this.#disconnect = new AbortController();

    this.#element.addEventListener("click", this, { signal: this.#disconnect.signal });
  }

  disconnectedCallback() {
    this.#disconnect.abort();
  }

  async handleEvent() {
    try {
      const content = new ClipboardItem({ "text/plain": this.value });
      await navigator.clipboard.write([content]);
      this.#showDidCopy();
    } catch (e) {
      alert("Could not copy the value to the clipboard.");
      console.error("error copying the value to the clipboard", e);
    }
  }

  #render() {
    const labelEl = this.#element.querySelector(".label > span");
    labelEl.textContent = this.label;

    const messageEl = this.#element.querySelector(".copied-message > span");
    messageEl.textContent = this.copiedMessage;

    this.classList.remove("copied");
    if (this.#didCopy) this.classList.add("copied");
  }

  #showDidCopy() {
    this.#didCopy = true;
    this.#render();

    setTimeout(() => {
      this.#didCopy = false;
      this.#render();
    }, 2000);
  }
}
