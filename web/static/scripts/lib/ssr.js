/**
 * @template T
 * @param {string} src
 * @returns {T}
 */
export function hydrate(src) {
  const el = document.getElementById(src);
  if (!el) throw new Error(`could not read serialized data from ${src}`);

  try {
    return JSON.parse(el.textContent);
  } catch (e) {
    throw new Error("failed to parse data", { cause: e });
  }
}
