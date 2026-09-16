import { CopyButton } from "../components/copyButton.js";

export function init() {
  for (const code of document.querySelectorAll("pre > code")) {
    const pre = code.parentElement;
    if (!(pre instanceof HTMLPreElement)) continue;

    const wrapper = document.createElement("div");
    wrapper.classList.add("code-block");
    pre.replaceWith(wrapper);
    wrapper.append(pre);

    const button = document.createElement(CopyButton.tag);
    button.setAttribute("label", "Copy");
    button.setAttribute("value", code.textContent ?? "");
    button.toggleAttribute("icononly", true);
    button.setAttribute("variant", "haptic");
    button.dataset.tooltip = "Copy this code";
    wrapper.append(button);
  }
}
