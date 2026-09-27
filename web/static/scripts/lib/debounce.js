/**
 * @template {unknown[]} T
 * @param {(...args: T) => void} fn
 * @param {number} wait
 * @param {number} [maxWait] Calls `fn` at the latest after this many ms, even if calls keep coming in.
 */
export function debounce(fn, wait, maxWait = Infinity) {
  /** @type {ReturnType<typeof setTimeout> | undefined} */
  let timeout;

  /** @type {number | undefined} */
  let start;

  /** @param {T} args */
  function debounced(...args) {
    start ??= Date.now();
    clearTimeout(timeout);

    const remaining = Math.min(wait, maxWait - (Date.now() - start));
    timeout = setTimeout(() => {
      start = undefined;
      fn(...args);
    }, Math.max(remaining, 0));
  }

  debounced.cancel = () => {
    clearTimeout(timeout);
    start = undefined;
  };

  return debounced;
}
