import { hydrate } from "../lib/ssr.js";

function createTagsState(src = "clientState.Tags") {
  /** @type {string[]} */
  const tags = hydrate(src);

  return { tags };
}

/** @type {ReturnType<typeof createTagsState>} */
let state;

export function getTags() {
  return (state ??= createTagsState());
}
