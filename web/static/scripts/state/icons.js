function createIconsState(selector = "template[id^='clientState.Icons.']") {
  /** @type {Map<string, Element>} */
  const icons = new Map();

  /**
   * @param {string} key
   * @returns {Element | undefined}
   */
  function getElement(key) {
    const node = icons.get(key)?.cloneNode(true);
    return node instanceof Element ? node : undefined;
  }

  const templates = document.querySelectorAll(selector);

  templates.entries().reduce((all, [, i]) => {
    if (i instanceof HTMLTemplateElement) {
      const name = i.id.substring(i.id.lastIndexOf(".") + 1);
      if (i.content.firstElementChild) all.set(name, i.content.firstElementChild);
    }
    return all;
  }, icons);

  return {
    getElement,
  };
}

/** @type {ReturnType<typeof createIconsState>} */
let state;

export function getIcons() {
  return (state ??= createIconsState());
}
