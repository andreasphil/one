function createSavedSearchState(key = "searches", capacity = 10) {
  /** @type {string[]} */
  const searches = [];

  /** @param {string} search */
  function save(search) {
    const existing = searches.indexOf(search);
    if (existing > -1) searches.splice(existing, 1);

    searches.splice(0, 0, search);
    if (searches.length > capacity) searches.splice(capacity, searches.length - capacity);

    saveToStorage();
  }

  /** @param {string} search */
  function remove(search) {
    const existing = searches.indexOf(search);
    if (existing > -1) searches.splice(existing, 1);

    saveToStorage();
  }

  // Persisting ---------------------------------------------

  function loadFromStorage() {
    const raw = localStorage.getItem(key);
    if (!raw) return;

    try {
      searches.push(...JSON.parse(raw));
    } catch (e) {
      console.error("failed to restore searches from localStorage", e);
    }
  }

  function saveToStorage() {
    localStorage.setItem(key, JSON.stringify(searches));
  }

  loadFromStorage();

  return {
    removeSearch: remove,
    saveSearch: save,
    searches: () => [...searches],
  };
}

/** @type {ReturnType<typeof createSavedSearchState>} */
let state;

export function getSavedSearch() {
  return (state ??= createSavedSearchState());
}
