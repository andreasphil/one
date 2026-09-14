import { newNavigationAction, useCommands } from "../lib/commands.js";
import { formatDate } from "../lib/format.js";
import { getIcons } from "../state/icons.js";
import { getNotes } from "../state/notes.js";
import { getTags } from "../state/tags.js";

export function init() {
  const { getElement: icon } = getIcons();
  const { notes } = getNotes();
  const { tags } = getTags();
  const { register } = useCommands();

  /** @param {import("../lib/types.js").NoteMeta} note */
  function formatNoteTitle(note) {
    let title = note.Title;
    if (note.IsChildNote && note.Date) title = `${title} (${formatDate(new Date(note.Date))})`;

    return title;
  }

  /** @type {import("../lib/commands.js").Command[]} */
  const notesCommands = notes.map((i) => ({
    id: `note:open:${i.Slug}`,
    name: formatNoteTitle(i),
    groupName: "Notes",
    icon: icon("StickyNote"),
    action: newNavigationAction(`/notes/${i.Slug}/`),
  }));

  /** @type {import("../lib/commands.js").Command[]} */
  const tagsCommands = tags.map((i) => ({
    id: `tag:open:${i}`,
    name: i,
    alias: [`#${i}`],
    groupName: "Tags",
    icon: icon("Tag"),
    action: newNavigationAction(`/tags/${i}/`),
  }));

  /** @type {import("../lib/commands.js").Command[]} */
  const commands = [
    // Navigation
    {
      id: "open:search",
      name: "Search",
      chord: "gs",
      groupName: "Open",
      icon: icon("Search"),
      action: newNavigationAction("/search/"),
    },
    {
      id: "open:notes",
      name: "Notes",
      chord: "gn",
      groupName: "Open",
      icon: icon("StickyNote"),
      action: newNavigationAction("/notes/"),
    },
    {
      id: "open:tags",
      name: "Tags",
      chord: "gx",
      groupName: "Open",
      icon: icon("Tag"),
      action: newNavigationAction("/tags/"),
    },
    {
      id: "open:today",
      name: "Today",
      chord: "gt",
      groupName: "Open",
      icon: icon("Calendar"),
      action: newNavigationAction(`/notes/${new Date().toISOString().substring(0, 10)}/`),
    },

    // Utils
    {
      id: "reveal-in-sidebar",
      name: "Reveal in sidebar",
      groupName: "App",
      icon: icon("Pipette"),
      action: () => {
        const current = document.querySelector('[aria-current="page"]');
        current?.scrollIntoView({ behavior: "smooth", block: "center" });
      },
    },
  ];

  register(...notesCommands, ...tagsCommands, ...commands);
}
