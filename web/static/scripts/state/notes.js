import { hydrate } from "../lib/ssr.js";

function createNotesState(src = "clientState.NotesMeta") {
  /** @type {import("../lib/types.js").NoteMeta[]} */
  const notes = hydrate(src);

  return { notes };
}

/** @type {ReturnType<typeof createNotesState>} */
let state;

export function getNotes() {
  return (state ??= createNotesState());
}
